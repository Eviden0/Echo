<template>
  <div class="monitor-container">
    <h2 class="monitor-title">当前监听 0.0.0.0:9999 给我弹 !!!</h2>
    <div class="buttons">
      <el-button @click="start" type="primary" class="btn-start">点我启动</el-button>
      <el-button @click="toggleTerminal" type="info" class="btn-toggle">隐藏/显示 终端</el-button>
    </div>
    <h1 class="session-id">当前id: {{ session }}</h1>
  </div>
  <div ref="terminal" id="terminal" :style="{ display: terminalVisible ? 'block' : 'none' }"></div>

</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import {GetId, Start} from "../../wailsjs/go/main/App.js";


const terminal = ref(null)
const fitAddon = new FitAddon()
const socket = ref(null)
const term = ref(null)
const inputBuffer = ref('')  // 当前输入行
const session = ref("")//sessionID
const terminalVisible = ref(true) // 控制终端可见性
let intervalId = null
let RemoteAdd=ref('')
// 切换终端可见性
function toggleTerminal() {
  terminalVisible.value = !terminalVisible.value
}

// 得到sessionID
Start()
function start(){
  GetId().then((id)=>{
    session.value=id
  })
  initWebSocket()
  // initTerm()
}
// 初始化 WebSocket
const initWebSocket = () => {
  //  拿到后端传来的session
  socket.value = new WebSocket(`ws://127.0.0.1:8080/ws?session=`+session.value)
  socket.value.onopen = () => console.log("WebSocket open")

  socket.value.onmessage = (event) => {
    term.value.write(event.data)
  }

  socket.value.onclose = () => console.log("WebSocket close")

  socket.value.onerror = (error) => console.log("WebSocket error:", error)
}

// 初始化 Terminal
const initTerm = () => {
  term.value = new Terminal({
    rows: 30,
    cols: 80,
    convertEol: true,
    cursorBlink: true,
    cursorStyle: 'block',
    theme: {
      foreground: '#ECECEC',
      background: '#000000',
      cursor: 'help',
      lineHeight: 20,
    },
  })

  term.value.open(document.getElementById('terminal'))
  term.value.loadAddon(fitAddon)
  setTimeout(() => fitAddon.fit(), 5)
  term.value.focus()

  // 挂载事件
  termData()
}

// 处理终端事件
const termData = () => {
  term.value.onData((data) => {
    // 退格处理
    if (data === '\x7F') {
      if (inputBuffer.value.length > 0) {
        inputBuffer.value = inputBuffer.value.slice(0, -1)
        term.value.write('\b \b')
      }
    }
    // 回车发送整行
    else if (data === '\r') {
      term.value.write('\r\n')
      socket.value.send(inputBuffer.value)  // 发送完整命令
      inputBuffer.value = ''
    }
    // 常规输入
    else {
      inputBuffer.value += data
      term.value.write(data)
    }
  })
}

// 销毁 WebSocket 和 Terminal
const closeSocketAndTerm = () => {
  socket.value && socket.value.close()
  term.value && term.value.dispose()
}

onMounted(() => {
  // initWebSocket()
  initTerm()
  intervalId = setInterval(() => {
    GetId().then((id) => {
      if (id ==="") {
      session.value="等待连接中..."
      }else {
        session.value=id
      }
    })
  }, 1000)
})

onBeforeUnmount(() => {
  closeSocketAndTerm()
  if (intervalId) {
    clearInterval(intervalId)
  }
})
</script>

<style  scoped>
#terminal {
  width: 100%;
  height: 100%;
}
.monitor-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.monitor-title {
  font-size: 1.5rem;
  color: #409eff;
  margin-bottom: 20px;
  text-align: center;
}

.buttons {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.session-id {
  font-size: 1.2rem;
  color: #606266;
}
</style>

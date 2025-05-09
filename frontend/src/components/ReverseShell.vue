<template>
  <div class="monitor-container">
    <h2 class="monitor-title">当前监听 0.0.0.0:9999 给我弹 !!!</h2>
    <div class="input-row">
<!--      <el-input v-model="session" placeholder="请输入 Session ID" class="session-input" />-->
      <el-button @click="start" type="primary" class="btn-start">点我启动</el-button>
      <el-button @click="toggleTerminal" type="info" class="btn-toggle">隐藏/显示终端</el-button>
      <el-button @click="closeShell" type="danger" class="btn-start">断开</el-button>
    </div>
    <h1 class="session-id">当前 ID: {{ session }}</h1>
    <div ref="terminal" id="terminal" :style="{ display: terminalVisible ? 'block' : 'none' }"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import {Callgologger, CloseCon, GetId, Start} from "../../wailsjs/go/main/App.js";
import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {LogInfo} from "../store/interface";


const terminal = ref(null)
const fitAddon = new FitAddon()
const socket = ref(null)
const term = ref(null)
const inputBuffer = ref('')  // 当前输入行
const session = ref("")//sessionID
const terminalVisible = ref(true) // 控制终端可见性
let intervalId = null

  // 断开连接并清空终端
function closeShell() {
    if (!session.value) {
      alert("Session ID 为空，无法关闭连接！")
      return
    }

    {
      // 发送 GET 请求到后端，通知关闭当前 session
      const response =  fetch(`http://127.0.0.1:8080/close?session=${session.value}`, {
        method: 'GET',  // 使用 GET 请求关闭 session
      })
      .then((response) => {
        if (response.ok) {
          console.log("关闭连接成功")
        } else {
          console.error("关闭连接失败")
        }
      })

      Callgologger("info","关闭连接"+session.value)

      // 清空 xterm 面板
      term.value.write('\x1b[2J\x1b[0f');  // 清屏 + 光标复位（不依赖缓冲区）
      term.value.clear()
      term.value.reset()
      // 关闭 WebSocket 连接
      socket.value?.close()

      // 清空输入缓存
      inputBuffer.value = ''
      session.value = ''  // 清空 session ID
      socket.value = null  // 确保断开后 socket 重置为 null

    }
  }
  //TODO 前端写个日志系统,进行展示,现在先在控制台打印
  EventsOn("gologger", (log: LogInfo) => {
    console.log(log.Level, log.Msg);
  });




// 切换终端可见性
function toggleTerminal() {
  terminalVisible.value = !terminalVisible.value
}

// 启动监听
Start()
function start(){
  // 如果已存在连接且连接未关闭，则不再重新初始化 WebSocket
  if (socket.value && socket.value.readyState === WebSocket.OPEN) {
    console.log("WebSocket 已经连接，无需重新连接")
    return
  }
  Callgologger("info","启动websocket")
  GetId().then((id)=>{
    session.value=id
    initWebSocket()  // 启动新的 WebSocket 连接
  })




}
// 初始化 WebSocket
const initWebSocket = () => {
  // 确保只有当 socket 为 null 时才初始化新的 WebSocket
  if (socket.value) {
    console.log("已经存在 WebSocket 连接")
    return
  }
  GetId().then((id)=>{
    session.value=id
  })
  //  拿到后端传来的session
  socket.value = new WebSocket(`ws://127.0.0.1:8080/ws?session=`+session.value)

  socket.value.onopen = () => Callgologger("info","收到sessionID: "+session.value)
  socket.value.onmessage = (event) => term.value?.write(event.data)
  socket.value.onclose = () => {
    Callgologger("info",session.value+" 已关闭")
    socket.value = null  // 确保 WebSocket 关闭后重置为 null
  }
  socket.value.onerror = (error) => Callgologger("info","WebSocket 错误:"+error)
}

// 初始化 Terminal
const initTerm = () => {
  term.value = new Terminal({
    rows: 30,
    cols: 80,
    convertEol: true,
    cursorBlink: true,
    rightClickSelectsWord: true,
    cursorStyle: 'block',
    fontFamily: '"Cascadia Code", Menlo, monospace',
    screenKeys: true,
    useStyle: true,
    theme: {
      foreground: '#ECECEC',
      background: '#000000',
      cursor: 'help',
      lineHeight: 20,
    },
  })

  term.value.open(document.getElementById('terminal')!)
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
  initTerm()
})

onBeforeUnmount(() => {
  socket.value?.close()
  term.value?.dispose()
})
</script>

<style scoped>
#terminal {
  width: 100%;
  height: 100%;
  min-height: 400px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #000;
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

.input-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
}

.session-input {
  width: 300px;
}

.session-id {
  font-size: 1.2rem;
  color: #606266;
}
</style>

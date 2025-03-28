<template>

  <el-dialog
    v-model="navStore.dialogVisiable"
    title="payload"
    width="500"
    align-center
  >
    <template #footer>
      <div class="dialog-footer">
        <el-input v-model="fileName" placeholder="请输入文件名称" class="keyInput"></el-input>
                <el-button @click="navStore.dialogVisiable = false">取消</el-button>
                <el-button type="primary" @click="inputK">
                    确认
                </el-button>
      </div>
    </template>
  </el-dialog>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import { useNavStore } from '../../store/nav';
const navStore=useNavStore()
const fileName=('payload.sh')
const inputK = () => {
    navStore.dialogVisiable=false
    downloadContent(fileName,navStore.payloadContent)
}
 function downloadContent(fileName, content) {
    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileName
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
}
</script>

<style>
.keyInput {
    margin-bottom: 30px;
}
</style> 
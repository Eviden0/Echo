<template>
    <div class="reverse">
        <div class="upper">
            <el-card class="mainContent" :style="{ backgroundColor: isDark ? '#303030' : '#ecefff' }" shadow="always"
                @click="focusEditable">
                <div class="content">
                    <span @click="updateCommand">🚀</span>
                    <span ref="editable" contenteditable class="editable"></span>
                </div>
            </el-card>
        </div>

        <div class="options">
            <el-select v-model="shell" filterable placeholder="Select" style="width: 120px">
                <el-option v-for="item in shellsOptions" :key="item.value" :label="item.label" :value="item.value"
                    @click="updateCommand(currentShellCommand)" />
            </el-select>
            <el-select v-model="encoding" placeholder="Select" style="width: 120px">
                <el-option v-for="item in encodingsOptions" :key="item.value" :label="item.label" :value="item.value"
                    @click="updateCommand(currentShellCommand)" />
            </el-select>
            <el-button type="primary" style="width: 100px;margin-top: 10px;" @click="copyToClipboard">COPY</el-button>
        </div>
    </div>
</template>

<script setup lang='js'>
import { ref, defineExpose, onMounted } from 'vue'
import { useDark } from '@vueuse/core'
import { shells, encodings, reverseShellCommands } from '../../utils/myData'
import { useDataStore } from '../../store/data';
import { ElMessage } from 'element-plus'
const dataStore = useDataStore()
const isDark = useDark()
const currentShellCommand = ref('Bash -i')
const editable = ref(null)

function focusEditable() {
    if (editable.value) {
        editable.value.focus()
    }
}
const shell = ref('powershell')
const encoding = ref('None')
const shellsOptions = shells.map(shell => ({
    value: shell,
    label: shell,
}))
const encodingsOptions = encodings.map(shell => ({
    value: shell,
    label: shell,
}))
function escapeHTML(str) {
    return str.replace(/[&<>"']/g, function (match) {
        const escape = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#39;'
        };
        return escape[match];
    });
}
function updateCommand(commandName) {
    if (!commandName)
        return
    currentShellCommand.value = commandName
    let command = reverseShellCommands.find(item => item.name === commandName).command
    if (!command)
        return
    if (encoding.value === 'None') {
        command = escapeHTML(command)
        command = dataStore.rsg.highlightParameters(command)
        command = command.replace('{port}', dataStore.rsg.getPort())
        command = command.replace('{ip}', dataStore.rsg.getIP())
        command = command.replace('{payload}', dataStore.rsg.getPayload())
        command = command.replace('{type}', dataStore.rsg.getType())
        command = command.replace('{shell}', shell.value)

    
    }
    if (encoding.value === 'Base64') {
        command = command.replace('{port}', dataStore.rsg.getPort())
        command = command.replace('{ip}', dataStore.rsg.getIP())
        command = command.replace('{payload}', dataStore.rsg.getPayload())
        command = command.replace('{type}', dataStore.rsg.getType())
        command = command.replace('{shell}', shell.value)
        command = btoa(command)
    } else if (encoding.value === 'URL Encode') {

        command = dataStore.rsg.fixedEncodeURIComponent(command)
        command = dataStore.rsg.highlightParameters(command)
        command = command.replace('%7Bport%7D', dataStore.rsg.getPort)
        command = command.replace('%7Bip%7D', dataStore.rsg.getIP())
        command = command.replace('%7Bpayload%7D', dataStore.rsg.getPayload())
        command = command.replace('%7Btype%7D', dataStore.rsg.getType())
        command = command.replace('%7Bshell%7D', shell.value)
    }
    else if(encoding.value === 'Double URL Encoude')
    {
        command = dataStore.rsg.fixedEncodeURIComponent(command)
        command = dataStore.rsg.fixedEncodeURIComponent(command)
        command = dataStore.rsg.highlightParameters(command,dataStore.rsg.fixedEncodeURIComponent)
        command = command.replace('%257Bport%257D', dataStore.rsg.getPort)
        command = command.replace('%257Bip%257D', dataStore.rsg.getIP())
        command = command.replace('%257Bpayload%257D', dataStore.rsg.getPayload())
        command = command.replace('%257Btype%257D', dataStore.rsg.getType())
        command = command.replace('%257Bshell%257D', shell.value)
    }
    editable.value.innerHTML = command
    // console.log(command)
}

defineExpose({
    updateCommand,
    getEditableContent() {
        return editable.value ? editable.value.textContent : ''
    }
})
onMounted(() => {
    updateCommand('Bash -i')
})
function copyToClipboard() {
    if (editable.value) {
        navigator.clipboard.writeText(editable.value.innerText).then(() => {
            ElMessage({
                type: 'success',
                message: '复制成功!',
                duration: 1000
            })
        }).catch(err => {
            ElMessage.error({
                message: '复制失败!',
                duration: 1000
            })
        })
    }
}
</script>

<style lang="scss" scoped>
.reverse {
    display: flex;
    flex-direction: column;
    padding: 10px 0 0 20px;

    .upper {
        height: 300px;
    }

    .mainContent {
        max-height: 280px;
        /* 其他样式 */
        cursor: text;
        /* 鼠标指针变为文本选择 */
    }

    .content {
        max-height: 280px;
        overflow-y: auto;
        display: flex;
        align-items: flex-start;
        /* 确保文本在顶部对齐 */
    }

    .editable {
        outline: none;
        /* 移除编辑时的边框 */
        flex: 1;
        /* 使文本区域占据剩余空间 */
        word-break: break-word;
        /* 确保长单词换行 */
        white-space: pre-wrap;
        /* 保留空白符并换行 */
    }

    .options {
        flex: 1;
        display: flex;
        justify-content: space-between;
        padding: 10px 0 0 0;
    }
}
</style>
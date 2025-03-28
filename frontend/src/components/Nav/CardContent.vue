<template>
    <div class="CardContent">
        <div class="upper">
            <el-card class="mainContent" :style="{ backgroundColor: isDark ? '#303030' : '#ecefff' }" shadow="always"
                @click="focusEditable">
                <div class="content">
                    <span @click="updateCommand">🚀</span>
                    <span ref="editable" contenteditable class="editable"></span>
                </div>
            </el-card>
        </div>

        <div class="options" v-show="showOptions">
                 <el-select v-model="shell" filterable placeholder="Select" style="width: 120px">
                    <el-option v-for="item in shellsOptions" :key="item.value" :label="item.label" :value="item.value"
                        @click="updateCommand(currentShellCommand)" />
                </el-select>
                <el-select v-model="encoding" placeholder="Select" style="width: 120px">
                    <el-option v-for="item in encodingsOptions" :key="item.value" :label="item.label"
                        :value="item.value" @click="updateCommand(currentShellCommand)" />
                </el-select>
        </div>
        <div class="acctions">
            <el-button type="primary" style="width: 150px;margin-top: 10px;" @click="openWebPage" v-show="showDownload">Download Listener</el-button>
            <el-button type="primary" style="width: 80px;margin-top: 10px;" @click="copyToClipboard">COPY</el-button>
        </div>

    </div>
</template>

<script setup lang='js'>
import { ref, defineExpose, onMounted } from 'vue'
import { useDark } from '@vueuse/core'
import { shells, encodings } from '../../utils/myData'
import { useDataStore } from '../../store/data';
import { ElMessage } from 'element-plus'
const dataStore = useDataStore()
const isDark = useDark()
const currentShellCommand = ref('Bash -i')
const editable = ref(null)
const props = defineProps(['shellCommands', 'showOptions', 'showDownload'])
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
    console.log(commandName, props.shellCommands)
    let command = props.shellCommands.find(item => item.name === commandName).command
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
    else if (encoding.value === 'Double URL Encoude') {
        command = dataStore.rsg.fixedEncodeURIComponent(command)
        command = dataStore.rsg.fixedEncodeURIComponent(command)
        command = dataStore.rsg.highlightParameters(command, dataStore.rsg.fixedEncodeURIComponent)
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
    // updateCommand('Bash -i')
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
function openWebPage() {
    window.open('https://github.com/t3l3machus/hoaxshell/tree/main/revshells', '_blank')
}
</script>

<style lang="scss" scoped>
.CardContent {
    display: flex;
    flex-direction: column;
    padding: 10px 0 0 20px;

    .upper {
        height: 300px;
    }

    .mainContent {
        max-height: 296px;
        /* 其他样式 */
        cursor: text;
        /* 鼠标指针变为文本选择 */
    }

    .content {
        max-height: 270px;
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
        justify-content: space-evenly;
        padding: 10px 0 0 0;
    }
    .acctions{
        display: flex;
        justify-content: end;
    }
}
</style>
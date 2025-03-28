<template>
    <div class="setting">
        <div class="ipAndPort">
            <el-card class="custom-card" :style="{ backgroundColor: isDark ? '#303030' : 'white' }">
                <div class="header">
                    <span>IP & Port</span>
                </div>
                <div class="content">
                    <el-input v-model="dataStore.rsg.ip" style="width: 190px" placeholder="ip" @input="onInputChange">
                        <template #prepend>IP</template>
                    </el-input>
                    <div class="port">
                        <el-input v-model="dataStore.rsg.port" style="width: 190px" placeholder="port"
                            @input="onInputChange">
                            <template #prepend>PORT</template>
                        </el-input>
                        <el-button style="width: 10px;" @click="imcrement">+1</el-button>
                    </div>
                </div>
            </el-card>
        </div>
        <div class="listener">
            <el-card class="custom-card" :style="{ backgroundColor: isDark ? '#303030' : 'white' }">
                <div class="header">
                    <span>Listener</span>
                    <div>
                        <!-- <el-button class="ebtn" @click="reset" type="danger" :icon="RefreshLeft" circle /> -->
                        <el-switch v-model="isShowAD" />
                    </div>
                </div>
                <div class="content">
                    <el-card class="mainContent" :style="{ backgroundColor: isDark ? '#303030' : '#ecefff' }" shadow="always" @click="focusEditable">
                        <div class="editable-container">
                            <span>🚀</span>
                            <span ref="editableDiv" contenteditable class="editable"></span>
                        </div>
                    </el-card>
                    <el-select v-model="advance" filterable placeholder="Select" size="large"
                        style="width: 200px;margin-top: 10px;" v-if="isShowAD" @change="updateListenerCommand">
                        <el-option v-for="item in listenerCommands" :key="item.type" :label="item.type"
                            :value="item.type" />
                    </el-select>
                    <div class="copy">
                        <el-button type="primary" style="width: 100px;margin-top: 10px;"
                            @click="copyToClipboard">COPY</el-button>
                    </div>
                </div>
            </el-card>
        </div>
    </div>
</template>

<script setup lang='js'>
import { ref,onMounted } from 'vue'
import { useDataStore } from '../store/data';
import listenerCommandsData from '../assets/json/listenerCommands.json'
import { ElMessage } from 'element-plus'
import { useDark, useToggle } from '@vueuse/core'
const isDark = useDark()
const dataStore = useDataStore()
const isShowAD = ref(true)
const advance = ref('nc')
const listenerCommands = ref(listenerCommandsData) // 正确引用导入的 JSON 数据
const editableDiv = ref(null)

function imcrement() {
    dataStore.portImcrement()
    onInputChange()
}

function copyToClipboard() {
    if (editableDiv.value) {
        navigator.clipboard.writeText(editableDiv.value.innerText).then(() => {
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

function onInputChange() {
    if (dataStore.rsg.port > 65535 || dataStore.rsg.port < 0) {
        dataStore.rsg.port = 9001
    }
    updateListenerCommand()
}

function updateListenerCommand() {
    if (!advance.value) return
    let command = listenerCommands.value.find(item => item.type === advance.value).result
    command = dataStore.rsg.highlightParameters(command)
    command = command.replace('{port}', dataStore.rsg.getPort())
    command = command.replace('{ip}', dataStore.rsg.getIP())
    command = command.replace('{payload}', dataStore.rsg.getPayload())
    command = command.replace('{type}', dataStore.rsg.getType())
    if (dataStore.rsg.getPort() < 1024) {
        command = `<span class="highlighted-warning">sudo</span> ${command}`
    }
    editableDiv.value.innerHTML = command
}

function focusEditable() {
    if (editableDiv.value) {
        editableDiv.value.focus()
    }
}
onMounted(() => {
    updateListenerCommand()
})
</script>

<style lang="scss">
.setting {
    display: flex;
    justify-content: space-evenly;
    /* 确保两个卡片之间有间距 */
    margin-top: 16px;
    flex: 1;
    flex-wrap: wrap;
    /* 使设置区域填满剩余空间 */

    .ipAndPort,
    .listener {
        display: flex;
        flex-direction: column;
        margin-bottom: 10px;

        .header {
            margin-bottom: 20px;
        }

        .content {
            display: flex;
            justify-content: space-evenly;
            flex: 1;
            /* 使内容区域填满卡片 */
        }
    }

    .ipAndPort {
        .port {
            display: flex;
            align-items: center;
        }
    }

    .listener {
        .header {
            display: flex;
            justify-content: space-between;
            line-height: 20px;
        }

        .content {
            display: flex;
            flex-direction: column;

            .mainContent {
                max-height: 160px;
                overflow-y: auto;
                /* 其他样式 */
                cursor: text; /* 鼠标指针变为文本选择 */

                .editable-container {
                    display: flex;
                    align-items: flex-start; /* 确保文本在顶部对齐 */
                    width: 100%;
                }

                .editable {
                    outline: none; /* 移除编辑时的边框 */
                    flex: 1; /* 使文本区域占据剩余空间 */
                    word-break: break-word; /* 确保长单词换行 */
                    white-space: pre-wrap; /* 保留空白符并换行 */
                }
            }
        }
    }

    .ebtn {
        margin-right: 20px;
    }

    .copy {
        display: flex;
        justify-content: end;
        margin-top: 10px;
    }
}

.custom-card {
    width: 500px;
    display: flex;
    flex-direction: column;
    flex: 1;
    /* 使卡片等高 */
}

.highlighted-parameter {
    color: #7223B5;
    font-weight: bold;
}

.highlighted-warning {
    color: red;
    font-weight: bold;
}
</style>
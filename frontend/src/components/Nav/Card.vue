<template>
    <el-card shadow="always" :style="{ backgroundColor: isDark ? '#303030' : 'white' }">
        <template #header>
            <div class="card-header">
                <span class="OS">OS</span>
                <el-select class="select" v-model="os" :empty-values="[null, undefined]" :value-on-clear="null"
                    placeholder="Select" style="width: 120px" @change="onOSchange">
                    <el-option v-for="item in osList" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
                <span>Name</span>
                <el-input class="input" v-model="name" placeholder="Name" size="large" @input="onNameChange"></el-input>
                <el-switch class="switch" v-model="isDark" :active-icon="Moon" :inactive-icon="Sunny" inline-prompt
                    @change="toggleDark" />
                <el-button type="primary" @click="getEditableContent">Generate</el-button>
            </div>
        </template>
        <div class="content">
            <div class="aside">
                <ul class="scrollable-list">
                    <li v-for="item in listData" :key="item.name" :class="getItemClass(item.name)" @click="onListCilck(item.name)">
                        {{ item.name }}
                    </li>
                </ul>
            </div>
            <div class="mainContent">
                 <CardContent v-show="navStore.nav=='Reverse'" ref="reverse" :shellCommands="reverseShellCommands" :showOptions="true"/>
                 <CardContent v-show="navStore.nav=='Bind'" ref="bind" :shellCommands="bindShellCommands" />
                 <CardContent v-show="navStore.nav=='MSFVenom'" ref="msfvenom" :shellCommands="msfvenomCommands" />
                 <CardContent v-show="navStore.nav=='HoaxShell'" ref="hoaxShell" :shellCommands="hoaxShellCommands" :showDownload="true"/>

            </div>
        </div>
    </el-card>
</template>

<script setup lang='js'>
import { ref, defineExpose } from 'vue'
import { useDark } from '@vueuse/core'
import { osList, reverseShellCommands, bindShellCommands, msfvenomCommands, hoaxShellCommands } from '../../utils/myData'
import { useNavStore } from '../../store/nav'
import Reverse from './Reverse.vue'
import CardContent from './CardContent.vue'

const navStore = useNavStore()
const isDark = useDark()
const name = ref('')
const os = ref('All')
const reverse = ref(null)
const bind = ref(null)
const msfvenom = ref(null)
const hoaxShell = ref(null)
const selectedItem = ref(null)
function onOSchange() {
    updateList()
}

function onNameChange() {
    updateList()
}

const listData = ref(reverseShellCommands) // 默认为这个

function updateList() {
    let currentShellCommands = reverseShellCommands
    if (navStore.nav === 'Reverse') {
        currentShellCommands = reverseShellCommands
    } else if (navStore.nav === 'Bind') {
        currentShellCommands = bindShellCommands
    } else if (navStore.nav === 'MSFVenom') {
        currentShellCommands = msfvenomCommands
    } else if (navStore.nav === 'HoaxShell') {
        currentShellCommands = hoaxShellCommands
    }
    if (os.value === 'All') {
        listData.value = currentShellCommands
    } else {
        listData.value = currentShellCommands.filter(item => item.meta.includes(os.value.toLowerCase()))
    }
    if (name.value) {
        listData.value = listData.value.filter(item => item.name.toLowerCase().includes(name.value.toLowerCase()))
    }

    //这里要修改
    if(!selectedItem)
    selectedItem.value=listData.value[0].name
}

function onListCilck(commandName) {
    selectedItem.value = commandName
    if (navStore.nav === 'Reverse') {
        reverse.value.updateCommand(commandName)
    }
    if(navStore.nav === 'Bind'){
        bind.value.updateCommand(commandName)
    }
    if(navStore.nav === 'MSFVenom'){
        msfvenom.value.updateCommand(commandName)
    }
    if(navStore.nav === 'HoaxShell'){
        hoaxShell.value.updateCommand(commandName)
    }
}

function getItemClass(itemName) {
    return {
        selected: selectedItem.value === itemName,
        'selected-dark': selectedItem.value === itemName && isDark.value
    }
}
function getEditableContent() {
    if(navStore.nav === 'Reverse'){
        const content =reverse.value.getEditableContent()
        navStore.setPayloadContent(content)
    }
    if(navStore.nav === 'Bind'){
        const content =bind.value.getEditableContent()
        navStore.setPayloadContent(content)
    }
    if(navStore.nav === 'MSFVenom'){
        const content =msfvenom.value.getEditableContent()
        navStore.setPayloadContent(content)
    }
    if(navStore.nav === 'HoaxShell'){
        const content =hoaxShell.value.getEditableContent()
        navStore.setPayloadContent(content)
    }
}
defineExpose({
    updateList, os
})
</script>

<style lang="scss" scoped>
span,
.select,
.input,
.switch {
    margin-right: 20px;
}

.input {
    min-width: 150px;
    max-width: 450px;
}

.content {
    display: flex;

    .aside {
        flex: 2;
        // display: flex;
        // flex-direction: column;
    }

    .scrollable-list {
        // flex: 2;
        display: flex;
        flex-direction: column;
        padding: 0;
        margin: 0;
        list-style: none;
        max-height: 380px;
        /* 设置固定高度 */
        overflow-y: auto;
        /* 启用垂直滚动条 */
    }

    li {
        flex: none;
        // color: #7f7f7f;
        font-size: 14px;
        line-height: 30px;
        cursor: pointer;
        padding: 10px;
        border: 1px solid gray;
        border-top: none;
        /* 确保两个相邻的上下边框不重叠 */
        box-sizing: border-box;
        border-radius: 2px;
    }

    li:first-child {
        border-top: 1px solid black;
        /* 为第一个 li 元素添加顶部边框 */
    }

    li:hover {
        background-color: skyblue
    }

    li.selected {
        background-color: lightblue;
    }

    li.selected-dark {
        background-color: #375A7F;
    }

    .mainContent {
        flex: 7;
    }
}
</style>
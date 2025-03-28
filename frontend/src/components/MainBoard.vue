<template>
    <el-card class="mainBoard">
        <div class="theme">
            <el-text class="mx-1" type="info">theme</el-text>
            <el-switch v-model="isDark" :active-icon="Moon" :inactive-icon="Sunny" inline-prompt @change="toggleDark" />
        </div>
        <div class="title">
            <el-text class="mx-1">Reverse Shell Generator</el-text>
        </div>
        <setting />
        <Nav />
    </el-card>
</template>

<script setup lang='js'>

import { ref } from 'vue'
import { Sunny, Moon } from '@element-plus/icons-vue';
import { useDark, useToggle } from '@vueuse/core'
import { useDataStore } from '../store/data';
import Setting from './Setting.vue';
import Nav from './Nav/index.vue';
const dataStore = useDataStore()
const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<style lang="scss" scoped>
.mainBoard {
    display: flex;
    width: 1200px;
    min-width: 600px;
    min-height: 900px;
    border-radius: 15px;
    flex-direction: column;

    .theme {
        display: flex;
        align-items: center;
        margin-bottom: 10px;

        .mx-1 {
            margin-right: 20px;
        }
    }

    .title {
        display: flex;
        align-items: center;
        justify-content: center;

        .mx-1 {
            font-size: 35px;
        }
    }

}

.hightLight {
    color: #7223B5;
    font-weight: bold;
}
</style>
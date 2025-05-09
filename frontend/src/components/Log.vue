<script setup lang="ts">
import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import { LogInfo } from "../store/interface";
// 监听 'gologger' 事件
EventsOn("gologger", (log: LogInfo) => {
  console.log(log.Level, log.Msg);
});
</script>

<template>

</template>

<style scoped>

</style>
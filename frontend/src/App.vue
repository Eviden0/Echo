<script lang="ts" setup>
import Main from "./components/Main.vue";
import MainBoard from "./components/MainBoard.vue";
</script>

<template>
  <div class="common-layout">
    <Main/>
    <MainBoard></MainBoard>
  </div>
</template>

<style>
.common-layout {
  margin: 0; /* 消除父元素的外边距 */
  padding: 0; /* 消除父元素的内边距 */
  background-image: url('./assets/images/gopher.png');
  background-size: contain; /* 或者 'contain' 根据需要选择 */
  background-position: center; /* 图片居中 */
  height: 100vh; /* 设置元素的高度 */
  position: relative; /* 使子元素绝对定位时基于此元素 */
}

.common-layout::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.9); /* 设定白色背景和透明度 */
  z-index: 1; /* 确保背景在最底层 */
}

.common-layout > * {
  position: relative; /* 使子元素在最上层 */
  z-index: 2; /* 确保内容在背景和透明层之上 */
}
*{
  margin: 0;
  padding: 0;
}
</style>

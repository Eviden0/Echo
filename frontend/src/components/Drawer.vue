<template>
  <el-button @click="visible = true">
    管理
  </el-button>
  <el-drawer v-model="visible" :show-close="false" size="60%" class="custom-drawer">
    <template #header="{ close, titleId, titleClass }">
      <h3 :id="titleId" :class="titleClass">管理页</h3>
      <el-button type="danger" @click="close" class="close-button">
        <el-icon class="el-icon--left"><CircleCloseFilled /></el-icon>
        Close
      </el-button>
    </template>

    <div class="button-group">
      <el-button type="primary" @click="redirectToSignup">跳转到注册页面</el-button>
      <el-button type="primary" @click="redirectToSignin">跳转到登录页面</el-button>
      <el-button type="primary" @click="toggleIframe">{{ showIframe ? '隐藏' : '显示' }}</el-button>
      <el-button type="primary" @click="handleGenerate">生成用户信息</el-button>
    </div>

    <div v-if="showIframe" class="iframe-container">
      <div style="width: 100%; overflow: hidden">
        <iframe
            src="http://localhost:9200/"
            width="107%"
            height="490"
            style="border: none;">
        </iframe>
      </div>
    </div>

    <el-table :data="userData" style="width: 100%" class="user-table">
      <el-table-column prop="uname" label="用户名">
        <template #default="{ row }">
          <span>{{ row.uname }}</span>
          <el-button size="small" @click="copyToClipboard(row.uname)">复制</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="uemail" label="电子邮件">
        <template #default="{ row }">
          <span>{{ row.uemail }}</span>
          <el-button size="small" @click="copyToClipboard(row.uemail)">复制</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="uphone" label="电话号码">
        <template #default="{ row }">
          <span>{{ row.uphone }}</span>
          <el-button size="small" @click="copyToClipboard(row.uphone)">复制</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="upassword" label="密码">
        <template #default="{ row }">
          <span>{{ row.upassword }}</span>
          <el-button size="small" @click="copyToClipboard(row.upassword)">复制</el-button>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="{ $index }">
          <el-button size="small" type="danger" @click="deleteUser($index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-drawer>
</template>

<script lang="ts" setup>
import {h, ref} from 'vue'
import {ElButton, ElDrawer, ElNotification,ElMessage} from 'element-plus'
import { CircleCloseFilled } from '@element-plus/icons-vue'
import {Generate} from "../../wailsjs/go/main/App";

const visible = ref(false)
const showIframe = ref(false);
const toggleIframe = () => {
  showIframe.value = !showIframe.value;
};
const redirectToSignup = () => {
  ElNotification({
    title: '通知',
    message: '确保登录信息唯一，建议用生成用户信息登录！',
    type: 'warning',
  })
  window.open('https://dashboard.cpolar.com/signup', '_blank'); // 在新标签页中打开注册页面
};
const redirectToSignin = () => {
  ElNotification({
    title: '通知',
    message: '确保登录信息唯一，建议用生成用户信息登录！',
    type: 'warning',
  })
  window.open('http://localhost:9200/', '_blank'); // 在新标签页中打开注册页面
};
const userData = ref([]);

const handleGenerate = async () => {
  try {
    const result = await Generate(); // 获取返回的 UserInfo 对象
    console.log('Generated Data:', result); // 打印返回的对象

    userData.value.push(result); // 将生成的用户信息添加到 userData 数组中
  } catch (error) {
    console.error('Error generating data:', error);
  }
};
const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text)
      .then(() => {
        ElMessage('复制成功.')
      })
      .catch(err => {
        ElMessage('复制失败:');
      });
};
const deleteUser = (index) => {
  userData.value.splice(index, 1); // 删除指定索引的用户信息
  ElNotification({
    title: '通知',
    message: h('i', { style: 'color: teal' }, '信息已经删除！'),
  })
};
</script>
<style>
.custom-drawer {
  background-color: #f5f7fa; /* 背景色 */
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1); /* 阴影效果 */
}

.button-group {
  margin: 20px 0; /* 上下间距 */
}

.iframe-container {
  margin: 20px 0; /* 上下间距 */
}

.user-table {
  margin-top: 20px; /* 表格与其他内容的间距 */
}

.close-button {
  margin-left: 10px; /* Close按钮与标题的间距 */
}

</style>
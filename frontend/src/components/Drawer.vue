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
      <el-button type="danger" @click="DeleteConfigFile">重置配置信息</el-button>
      <el-button type="primary" @click="handleGenerate">生成用户信息</el-button>
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
    <div>
      <el-button @click="fetchUserTunnels" type="primary">获取隧道信息</el-button>
      <el-table :data="userTunnels" style="width: 100%" class="tunnel-table">
        <el-table-column prop="id" label="隧道ID"></el-table-column>
        <el-table-column prop="name" label="名称"></el-table-column>
        <el-table-column prop="public_url" label="URL"></el-table-column>
        <el-table-column prop="proto" label="协议"></el-table-column>
        <el-table-column prop="addr" label="地址"></el-table-column>
        <el-table-column prop="create_datetime" label="创建时间"></el-table-column>
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDeleteTunnel(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
<!--      添加tunnel-->
      <div class="add-tunnel-form">
        <el-select v-model="selectedProtocol" placeholder="选择协议">
          <el-option label="HTTP" value="http"></el-option>
          <el-option label="TCP" value="tcp"></el-option>
          <el-option label="FTP" value="ftp"></el-option>
          <el-option label="TLS" value="tls"></el-option>
        </el-select>
        <el-input v-model="tunnelAddress" placeholder="输入地址"></el-input>
        <el-button type="primary" @click="handleAddTunnel">添加隧道</el-button>
      </div>
    </div>
  </el-drawer>
</template>

<script lang="ts" setup>
import {h, ref} from 'vue'
import {ElButton, ElDrawer, ElNotification,ElMessage} from 'element-plus'
import { CircleCloseFilled } from '@element-plus/icons-vue'
import {AddTunnel, DeleteConfigFile, DeleteTunnel, Generate, GenerateTunnel} from "../../wailsjs/go/main/App";

const userTunnels = ref([])


const loading = ref(false);
const adding = ref(false);
const deleting = ref(false);
//  拿到接口的tunnel信息

const fetchUserTunnels = async () => {
  loading.value = true;
  try {
    const user = await GenerateTunnel();
    userTunnels.value = user.Tunnels || [];
  } catch (error) {
    console.error('Error fetching tunnels:', error);
    ElMessage.error('获取隧道信息失败');
  } finally {
    loading.value = false;
  }
};
// 添加tunnel
const selectedProtocol = ref('tcp')
const tunnelAddress = ref('0.0.0.0:9999')

const handleAddTunnel = async () => {
  if (!tunnelAddress.value) {
    ElMessage.warning('请输入隧道地址');
    return;
  }

  adding.value = true;
  try {
    await AddTunnel(tunnelAddress.value, selectedProtocol.value);
    ElMessage.success('隧道添加成功');
    await fetchUserTunnels();  // 添加后刷新列表
  } catch (error) {
    console.error('Error adding tunnel:', error);
    ElMessage.error('添加隧道失败');
  } finally {
    adding.value = false;
  }
};

//  删除对应的隧道
const handleDeleteTunnel = async (id: string) => {
  deleting.value = true;
  try {
    await DeleteTunnel(id);
    ElMessage.success('隧道删除成功');
    await fetchUserTunnels();  // 删除后刷新列表
  } catch (error) {
    console.error('Error deleting tunnel:', error);
    ElMessage.error('删除隧道失败');
  } finally {
    deleting.value = false;
  }
};



const visible = ref(false)
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
.tunnel-table {
  margin-top: 20px;
}

.add-tunnel-form {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}
</style>
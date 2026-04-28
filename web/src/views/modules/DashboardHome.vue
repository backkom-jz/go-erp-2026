<script setup>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import client from "../../api/client";

const profile = ref(null);
const loading = ref(false);

const loadProfile = async () => {
  loading.value = true;
  try {
    const resp = await client.get("/users/me");
    profile.value = resp.data;
  } catch (err) {
    ElMessage.error(err.message || "获取用户信息失败");
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="page-card">
    <h3>系统概览</h3>
    <p>这是 Go-ERP 后台管理 MVP 页面，可用于联调用户/商品/库存/订单接口。</p>
    <el-button type="primary" :loading="loading" @click="loadProfile">获取我的信息</el-button>
    <el-descriptions v-if="profile" style="margin-top: 16px" :column="1" border>
      <el-descriptions-item label="用户编号">{{ profile.user_no }}</el-descriptions-item>
      <el-descriptions-item label="用户名">{{ profile.name }}</el-descriptions-item>
      <el-descriptions-item label="租户">{{ profile.tenant_id }}</el-descriptions-item>
      <el-descriptions-item label="角色">{{ profile.role }}</el-descriptions-item>
    </el-descriptions>
  </div>
</template>

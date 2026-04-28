<script setup>
import { reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";
import client from "../api/client";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const authStore = useAuthStore();
const loading = ref(false);
const form = reactive({
  user_no: "admin",
  tenant_id: "t-default",
  role: "admin"
});

const submit = async () => {
  loading.value = true;
  try {
    const result = await client.post("/auth/login", form);
    authStore.setSession({
      token: result.data.access_token,
      userNo: form.user_no,
      tenantId: form.tenant_id
    });
    ElMessage.success("登录成功");
    router.push("/");
  } catch (err) {
    ElMessage.error(err.message || "登录失败");
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="login-wrap">
    <div class="login-hero">
      <h1>Go-ERP Admin</h1>
      <p>简约现代化后台 · Vue 3 + Element Plus</p>
    </div>
    <div class="login-card">
      <h2>登录系统</h2>
      <el-form :model="form" label-position="top">
        <el-form-item label="用户编号">
          <el-input v-model="form.user_no" size="large" />
        </el-form-item>
        <el-form-item label="租户 ID">
          <el-input v-model="form.tenant_id" size="large" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role" size="large" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="查看者" value="viewer" />
          </el-select>
        </el-form-item>
        <el-button type="primary" :loading="loading" size="large" style="width: 100%" @click="submit">
          进入后台
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login-wrap {
  min-height: 100vh;
  display: flex;
  gap: 56px;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at top left, #e0e7ff, #f8fafc 45%, #eff6ff);
}

.login-hero {
  max-width: 360px;
}

.login-hero h1 {
  margin: 0;
  font-size: 40px;
  line-height: 1.1;
  color: #1e293b;
}

.login-hero p {
  margin-top: 14px;
  color: #64748b;
}

.login-card {
  width: 420px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
}

h2 {
  margin-top: 0;
  margin-bottom: 18px;
  color: #0f172a;
}
</style>

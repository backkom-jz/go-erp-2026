<script setup>
import { reactive } from "vue";
import { ElMessage } from "element-plus";
import client from "../../api/client";

const form = reactive({
  user_no: "",
  name: "",
  tenant_id: "t-default",
  role: "viewer"
});

const submit = async () => {
  try {
    await client.post("/users", form);
    ElMessage.success("用户创建成功");
    form.user_no = "";
    form.name = "";
  } catch (err) {
    ElMessage.error(err.message || "创建失败");
  }
};
</script>

<template>
  <div class="page-card">
    <h3>用户管理</h3>
    <el-form :model="form" label-width="100px" style="max-width: 560px">
      <el-form-item label="用户编号">
        <el-input v-model="form.user_no" />
      </el-form-item>
      <el-form-item label="用户名">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="租户 ID">
        <el-input v-model="form.tenant_id" />
      </el-form-item>
      <el-form-item label="角色">
        <el-select v-model="form.role" style="width: 100%">
          <el-option label="管理员" value="admin" />
          <el-option label="查看者" value="viewer" />
        </el-select>
      </el-form-item>
      <el-button type="primary" @click="submit">创建用户</el-button>
    </el-form>
  </div>
</template>

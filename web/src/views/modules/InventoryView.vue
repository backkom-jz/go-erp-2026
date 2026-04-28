<script setup>
import { reactive } from "vue";
import { ElMessage } from "element-plus";
import client from "../../api/client";

const form = reactive({
  sku_id: 1,
  qty: 1,
  business_no: `manual-${Date.now()}`
});

const deduct = async () => {
  try {
    await client.post("/inventory/deduct", form);
    ElMessage.success("扣减成功");
    form.business_no = `manual-${Date.now()}`;
  } catch (err) {
    ElMessage.error(err.message || "扣减失败");
  }
};
</script>

<template>
  <div class="page-card">
    <h3>库存管理</h3>
    <el-alert type="warning" show-icon :closable="false" title="该页面调用扣减接口，建议先创建 SKU 并初始化库存数据。" />
    <el-form :model="form" label-width="120px" style="max-width: 560px; margin-top: 12px">
      <el-form-item label="SKU ID"><el-input-number v-model="form.sku_id" :min="1" /></el-form-item>
      <el-form-item label="扣减数量"><el-input-number v-model="form.qty" :min="1" /></el-form-item>
      <el-form-item label="业务单号"><el-input v-model="form.business_no" /></el-form-item>
      <el-button type="primary" @click="deduct">扣减库存</el-button>
    </el-form>
  </div>
</template>

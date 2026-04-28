<script setup>
import { reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import client from "../../api/client";

const createForm = reactive({
  user_id: 1,
  tenant_id: "t-default",
  items: [
    {
      sku_id: 1,
      qty: 1,
      price_cents: 1000
    }
  ]
});

const queryOrderId = ref("");
const latestOrder = ref(null);
const queryResult = ref(null);

const createOrder = async () => {
  try {
    const result = await client.post("/order/create", createForm);
    latestOrder.value = result.data;
    queryOrderId.value = String(result.data.id);
    ElMessage.success("订单创建成功");
  } catch (err) {
    ElMessage.error(err.message || "创建订单失败");
  }
};

const queryOrder = async () => {
  if (!queryOrderId.value) {
    ElMessage.warning("请输入订单 ID");
    return;
  }
  try {
    const result = await client.get(`/order/${queryOrderId.value}`);
    queryResult.value = result.data;
  } catch (err) {
    ElMessage.error(err.message || "查询失败");
  }
};
</script>

<template>
  <div class="page-card">
    <h3>订单管理</h3>
    <el-row :gutter="20">
      <el-col :span="12">
        <h4>创建订单</h4>
        <el-form :model="createForm" label-width="100px">
          <el-form-item label="用户 ID">
            <el-input-number v-model="createForm.user_id" :min="1" />
          </el-form-item>
          <el-form-item label="租户 ID">
            <el-input v-model="createForm.tenant_id" />
          </el-form-item>
          <el-form-item label="SKU ID">
            <el-input-number v-model="createForm.items[0].sku_id" :min="1" />
          </el-form-item>
          <el-form-item label="数量">
            <el-input-number v-model="createForm.items[0].qty" :min="1" />
          </el-form-item>
          <el-form-item label="单价(分)">
            <el-input-number v-model="createForm.items[0].price_cents" :min="1" />
          </el-form-item>
          <el-button type="primary" @click="createOrder">创建订单</el-button>
        </el-form>
        <el-alert v-if="latestOrder" style="margin-top: 10px" type="success" :closable="false" :title="`订单号: ${latestOrder.order_no}`" />
      </el-col>
      <el-col :span="12">
        <h4>查询订单</h4>
        <el-input v-model="queryOrderId" placeholder="输入订单ID" />
        <el-button style="margin-top: 10px" @click="queryOrder">查询</el-button>

        <el-descriptions v-if="queryResult?.order" style="margin-top: 12px" :column="1" border>
          <el-descriptions-item label="订单号">{{ queryResult.order.order_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ queryResult.order.status }}</el-descriptions-item>
          <el-descriptions-item label="总金额(分)">{{ queryResult.order.total_cents }}</el-descriptions-item>
        </el-descriptions>
      </el-col>
    </el-row>
  </div>
</template>

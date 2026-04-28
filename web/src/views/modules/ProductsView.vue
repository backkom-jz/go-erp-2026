<script setup>
import { reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import client from "../../api/client";

const spuForm = reactive({
  name: "",
  category_id: 1,
  brand: ""
});

const skuForm = reactive({
  spu_id: 0,
  code: "",
  name: "",
  price_cents: 1000
});

const spuList = ref([]);

const createSpu = async () => {
  try {
    await client.post("/products/spu", spuForm);
    ElMessage.success("SPU 创建成功");
    await querySpu();
  } catch (err) {
    ElMessage.error(err.message || "SPU 创建失败");
  }
};

const createSku = async () => {
  try {
    await client.post("/products/sku", skuForm);
    ElMessage.success("SKU 创建成功");
  } catch (err) {
    ElMessage.error(err.message || "SKU 创建失败");
  }
};

const querySpu = async () => {
  try {
    const result = await client.get("/products/spu?limit=20");
    spuList.value = result.data || [];
  } catch (err) {
    ElMessage.error(err.message || "查询 SPU 失败");
  }
};

querySpu();
</script>

<template>
  <div class="page-card">
    <h3>商品管理</h3>
    <el-row :gutter="20">
      <el-col :span="12">
        <h4>创建 SPU</h4>
        <el-form :model="spuForm" label-width="100px">
          <el-form-item label="名称"><el-input v-model="spuForm.name" /></el-form-item>
          <el-form-item label="分类 ID"><el-input-number v-model="spuForm.category_id" :min="1" /></el-form-item>
          <el-form-item label="品牌"><el-input v-model="spuForm.brand" /></el-form-item>
          <el-button type="primary" @click="createSpu">提交</el-button>
        </el-form>
      </el-col>
      <el-col :span="12">
        <h4>创建 SKU</h4>
        <el-form :model="skuForm" label-width="100px">
          <el-form-item label="SPU ID"><el-input-number v-model="skuForm.spu_id" :min="1" /></el-form-item>
          <el-form-item label="编码"><el-input v-model="skuForm.code" /></el-form-item>
          <el-form-item label="名称"><el-input v-model="skuForm.name" /></el-form-item>
          <el-form-item label="价格(分)"><el-input-number v-model="skuForm.price_cents" :min="1" /></el-form-item>
          <el-button type="primary" @click="createSku">提交</el-button>
        </el-form>
      </el-col>
    </el-row>

    <h4 style="margin-top: 20px">SPU 列表</h4>
    <el-button plain @click="querySpu">刷新</el-button>
    <el-table :data="spuList" style="margin-top: 12px" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="category_id" label="分类 ID" />
      <el-table-column prop="brand" label="品牌" />
    </el-table>
  </div>
</template>

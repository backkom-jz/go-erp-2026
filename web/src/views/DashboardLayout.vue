<script setup>
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const menus = [
  { path: "/dashboard", label: "首页", icon: "House" },
  { path: "/users", label: "用户管理", icon: "User" },
  { path: "/products", label: "商品管理", icon: "Goods" },
  { path: "/inventory", label: "库存管理", icon: "Box" },
  { path: "/orders", label: "订单管理", icon: "Document" },
  { path: "/ai-chat", label: "AI 对话", icon: "ChatDotRound" }
];

const activeMenu = computed(() => route.path);

const logout = () => {
  authStore.clearSession();
  router.push("/login");
};
</script>

<template>
  <el-container class="layout-wrap">
    <el-aside width="220px" class="aside">
      <div class="logo">Go-ERP Admin</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#0f172a"
        text-color="#cbd5e1"
        active-text-color="#ffffff"
      >
        <el-menu-item v-for="menu in menus" :key="menu.path" :index="menu.path">
          <el-icon><component :is="menu.icon" /></el-icon>
          <span>{{ menu.label }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span class="tenant-tag">当前租户：{{ authStore.tenantId || "-" }}</span>
        <el-button text @click="logout">退出登录</el-button>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-wrap {
  min-height: 100vh;
}

.aside {
  background: #0f172a;
  color: #fff;
}

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  letter-spacing: 0.4px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.25);
}

.header {
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.tenant-tag {
  color: #334155;
  font-weight: 500;
}

.main {
  background: #f8fafc;
}
</style>

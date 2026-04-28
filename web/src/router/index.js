import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
  {
    path: "/login",
    name: "login",
    component: () => import("../views/LoginView.vue")
  },
  {
    path: "/",
    component: () => import("../views/DashboardLayout.vue"),
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "dashboard",
        component: () => import("../views/modules/DashboardHome.vue")
      },
      {
        path: "users",
        name: "users",
        component: () => import("../views/modules/UsersView.vue")
      },
      {
        path: "products",
        name: "products",
        component: () => import("../views/modules/ProductsView.vue")
      },
      {
        path: "inventory",
        name: "inventory",
        component: () => import("../views/modules/InventoryView.vue")
      },
      {
        path: "orders",
        name: "orders",
        component: () => import("../views/modules/OrdersView.vue")
      },
      {
        path: "ai-chat",
        name: "ai-chat",
        component: () => import("../views/modules/AIChatView.vue")
      }
    ]
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem("erp_admin_token");
  if (to.path !== "/login" && !token) {
    next("/login");
    return;
  }
  if (to.path === "/login" && token) {
    next("/");
    return;
  }
  next();
});

export default router;

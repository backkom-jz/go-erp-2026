import { defineStore } from "pinia";

const TOKEN_KEY = "erp_admin_token";
const USER_NO_KEY = "erp_admin_user_no";
const TENANT_KEY = "erp_admin_tenant";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || "",
    userNo: localStorage.getItem(USER_NO_KEY) || "",
    tenantId: localStorage.getItem(TENANT_KEY) || ""
  }),
  getters: {
    isLoggedIn: (state) => !!state.token
  },
  actions: {
    setSession(payload) {
      this.token = payload.token;
      this.userNo = payload.userNo;
      this.tenantId = payload.tenantId;
      localStorage.setItem(TOKEN_KEY, this.token);
      localStorage.setItem(USER_NO_KEY, this.userNo);
      localStorage.setItem(TENANT_KEY, this.tenantId);
    },
    clearSession() {
      this.token = "";
      this.userNo = "";
      this.tenantId = "";
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_NO_KEY);
      localStorage.removeItem(TENANT_KEY);
    }
  }
});

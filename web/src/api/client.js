import axios from "axios";

const client = axios.create({
  baseURL: "/api/v1",
  timeout: 12000,
  withCredentials: true
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem("erp_admin_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (resp) => {
    const body = resp.data || {};
    if (typeof body.code === "number" && body.code !== 0) {
      const err = new Error(body.message || "请求失败");
      err.code = body.code;
      throw err;
    }
    return body;
  },
  (err) => {
    if (err?.response?.data?.message) {
      return Promise.reject(new Error(err.response.data.message));
    }
    return Promise.reject(err);
  }
);

export default client;

<script setup>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import client from "../../api/client";

const loading = ref(false);
const prompt = ref("");
const messages = ref([
  {
    role: "assistant",
    content: "你好，我是 DeepSeekV4-pro 对话助手。你可以问我业务设计、代码问题或运营分析。"
  }
]);

const sendMessage = async () => {
  const text = prompt.value.trim();
  if (!text || loading.value) {
    return;
  }

  messages.value.push({ role: "user", content: text });
  prompt.value = "";
  loading.value = true;

  try {
    const result = await client.post("/ai/chat", {
      message: text,
      history: messages.value.slice(0, -1).map((item) => ({
        role: item.role,
        content: item.content
      }))
    });
    messages.value.push({
      role: "assistant",
      content: result.data.reply
    });
  } catch (err) {
    ElMessage.error(err.message || "AI 对话失败");
    messages.value.push({
      role: "assistant",
      content: "当前请求失败，请检查后端 AI 配置（ai.enabled、ai.api_key）后重试。"
    });
  } finally {
    loading.value = false;
  }
};

const clearChat = () => {
  messages.value = [
    {
      role: "assistant",
      content: "会话已清空。请继续提问。"
    }
  ];
};
</script>

<template>
  <div class="page-card ai-page">
    <div class="ai-head">
      <div>
        <h3>AI 对话（DeepSeekV4-pro）</h3>
        <p>通过后端代理调用模型，前端不暴露密钥。</p>
      </div>
      <el-button plain @click="clearChat">清空会话</el-button>
    </div>

    <div class="chat-box">
      <div v-for="(item, idx) in messages" :key="idx" class="chat-row" :class="item.role">
        <div class="bubble">
          <div class="role">{{ item.role === "user" ? "你" : "AI" }}</div>
          <div class="content">{{ item.content }}</div>
        </div>
      </div>
      <div v-if="loading" class="chat-row assistant">
        <div class="bubble">
          <div class="role">AI</div>
          <div class="content">思考中...</div>
        </div>
      </div>
    </div>

    <div class="input-bar">
      <el-input
        v-model="prompt"
        type="textarea"
        :rows="3"
        placeholder="输入你的问题，Enter 发送，Shift+Enter 换行"
        @keydown.enter.prevent="sendMessage"
      />
      <el-button type="primary" :loading="loading" @click="sendMessage">发送</el-button>
    </div>
  </div>
</template>

<style scoped>
.ai-page {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ai-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ai-head p {
  margin: 4px 0 0;
  color: #64748b;
}

.chat-box {
  min-height: 440px;
  max-height: 60vh;
  overflow: auto;
  padding: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #f8fafc;
}

.chat-row {
  display: flex;
  margin-bottom: 10px;
}

.chat-row.user {
  justify-content: flex-end;
}

.bubble {
  max-width: 80%;
  padding: 10px 12px;
  border-radius: 10px;
  background: #fff;
  border: 1px solid #e2e8f0;
}

.chat-row.user .bubble {
  background: #1d4ed8;
  color: #fff;
  border-color: #1d4ed8;
}

.role {
  font-size: 12px;
  opacity: 0.8;
  margin-bottom: 4px;
}

.content {
  white-space: pre-wrap;
  line-height: 1.5;
}

.input-bar {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}
</style>

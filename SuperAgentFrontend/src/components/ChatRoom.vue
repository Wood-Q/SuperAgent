<template>
  <div class="chat-container">
    <!-- 头部 -->
    <div class="chat-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="title">🌙 MoonAgent小助手</h1>
        </div>
        <div class="header-right">
          <div class="session-id">会话ID: {{ sessionId }}</div>
        </div>
      </div>
    </div>

    <!-- 消息区域 -->
    <div class="messages-container" ref="messagesContainer">
      <div class="messages-list">
        <!-- 欢迎消息 -->
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="message-wrapper">
            <div class="ai-message">
              <div class="message-avatar ai-avatar">
                <img src="/image.png" alt="image" class="ai-avatar-image" />
              </div>
              <div class="message-content">
                <div class="message-bubble ai-bubble">
                  您好！我是MoonAgent小助手，有什么可以帮助您的吗？
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 聊天消息 -->
        <div v-for="message in messages" :key="message.id" class="message-item">
          <div
            class="message-wrapper"
            :class="{ 'user-message-wrapper': message.role === 'user' }"
          >
            <!-- 用户消息 -->
            <div v-if="message.role === 'user'" class="user-message">
              <div class="message-content">
                <div class="message-bubble user-bubble">
                  {{ message.content }}
                </div>
              </div>
              <div class="message-avatar user-avatar">👤</div>
            </div>

            <!-- AI消息 -->
            <div v-else class="ai-message">
              <div class="message-avatar ai-avatar">
                <img src="/image.png" alt="image" class="ai-avatar-image" />
              </div>
              <div class="message-content">
                <div class="message-bubble ai-bubble">
                  <div v-if="message.content" class="message-text">
                    {{ message.content
                    }}<span v-if="message.isTyping" class="cursor">|</span>
                  </div>
                  <div
                    v-if="message.isTyping && !message.content"
                    class="typing-indicator"
                  >
                    <span></span>
                    <span></span>
                    <span></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-container">
      <div class="input-wrapper">
        <div class="input-box">
          <!-- 流式输出开关 -->
          <div class="stream-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="isStreamMode" />
              <span class="slider"></span>
            </label>
            <span class="toggle-label">{{
              isStreamMode ? "流式" : "整体"
            }}</span>
          </div>

          <textarea
            v-model="inputMessage"
            @keydown="handleKeyDown"
            @input="adjustTextareaHeight"
            ref="messageInput"
            class="message-input"
            placeholder="输入您的问题..."
            rows="1"
            :disabled="isLoading"
          ></textarea>
          <button
            @click="sendMessage"
            class="send-button"
            :disabled="!inputMessage.trim() || isLoading"
          >
            <svg
              v-if="!isLoading"
              width="16"
              height="16"
              viewBox="0 0 16 16"
              fill="none"
            >
              <path
                d="M.5 1.163A1 1 0 0 1 1.97.28l12.868 6.837a1 1 0 0 1 0 1.766L1.969 15.72A1 1 0 0 1 .5 14.836V10.33a1 1 0 0 1 .816-.983L8.5 8 1.316 6.653A1 1 0 0 1 .5 5.67V1.163Z"
                fill="currentColor"
              />
            </svg>
            <div v-else class="loading-spinner"></div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, nextTick } from "vue";
import { v4 as uuidv4 } from "uuid";
import { streamChatWithModel, chatWithModel } from "../utils/api";

export default {
  name: "ChatRoom",
  setup() {
    const sessionId = ref("");
    const messages = reactive([]);
    const inputMessage = ref("");
    const isLoading = ref(false);
    const messagesContainer = ref(null);
    const messageInput = ref(null);
    const isStreamMode = ref(true);

    // 生成会话ID
    const generateSessionId = () => {
      sessionId.value = uuidv4().substring(0, 8);
    };

    // 滚动到底部
    const scrollToBottom = () => {
      nextTick(() => {
        if (messagesContainer.value) {
          const container = messagesContainer.value;
          container.scrollTop = container.scrollHeight;
        }
      });
    };

    // 平滑滚动到底部
    const smoothScrollToBottom = () => {
      nextTick(() => {
        if (messagesContainer.value) {
          const container = messagesContainer.value;
          container.scrollTo({
            top: container.scrollHeight,
            behavior: "smooth",
          });
        }
      });
    };

    // 调整输入框高度
    const adjustTextareaHeight = () => {
      const textarea = messageInput.value;
      if (textarea) {
        textarea.style.height = "auto";
        textarea.style.height = Math.min(textarea.scrollHeight, 120) + "px";
      }
    };

    // 处理键盘事件
    const handleKeyDown = (event) => {
      if (event.key === "Enter" && !event.shiftKey) {
        event.preventDefault();
        sendMessage();
      }
    };

    // 发送消息
    const sendMessage = async () => {
      const message = inputMessage.value.trim();
      if (!message || isLoading.value) return;

      // 添加用户消息
      const userMessage = {
        id: uuidv4(),
        role: "user",
        content: message,
        timestamp: new Date(),
      };
      messages.push(userMessage);

      // 清空输入框
      inputMessage.value = "";
      adjustTextareaHeight();

      // 添加AI消息占位符
      const aiMessage = {
        id: uuidv4(),
        role: "assistant",
        content: "",
        isTyping: isStreamMode.value,
        timestamp: new Date(),
      };
      messages.push(aiMessage);

      isLoading.value = true;
      scrollToBottom();

      try {
        if (isStreamMode.value) {
          // 流式输出模式
          await streamChatWithModel(
            message,
            (chunk) => {
              // 更新AI消息内容 - 逐字符追加
              aiMessage.content += chunk;
              aiMessage.isTyping = true; // 保持光标显示

              // 更频繁的滚动更新，让体验更丝滑
              if (messagesContainer.value) {
                const container = messagesContainer.value;
                const isAtBottom =
                  container.scrollHeight - container.clientHeight <=
                  container.scrollTop + 100;
                if (isAtBottom) {
                  nextTick(() => {
                    container.scrollTop = container.scrollHeight;
                  });
                }
              }
            },
            (error) => {
              console.error("Stream error:", error);
              aiMessage.content = "抱歉，发生了错误，请稍后重试。";
              aiMessage.isTyping = false;
            }
          );
        } else {
          // 非流式输出模式
          const response = await chatWithModel(message);
          aiMessage.content = response.message || "抱歉，没有收到回复。";
          aiMessage.isTyping = false;
        }
      } catch (error) {
        console.error("API error:", error);
        aiMessage.content = "抱歉，发生了错误，请稍后重试。";
        aiMessage.isTyping = false;
      } finally {
        // 完成后隐藏光标和加载状态
        aiMessage.isTyping = false;
        isLoading.value = false;
        smoothScrollToBottom();
      }
    };

    // 组件挂载
    onMounted(() => {
      generateSessionId();
    });

    return {
      sessionId,
      messages,
      inputMessage,
      isLoading,
      messagesContainer,
      messageInput,
      sendMessage,
      handleKeyDown,
      adjustTextareaHeight,
      isStreamMode,
    };
  },
};
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background: #343541;
  color: #ffffff;
}

.chat-header {
  padding: 0;
  background: #40414f;
  border-bottom: 1px solid #4d4d4f;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #ffffff;
  margin: 0;
}

.session-id {
  font-size: 12px;
  color: #8e8ea0;
  font-family: monospace;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  background: #343541;
  scroll-behavior: smooth;
}

.messages-list {
  padding: 24px 0;
}

.welcome-message {
  margin-bottom: 24px;
}

.message-item {
  margin-bottom: 0;
  animation: fadeInUp 0.3s ease-out;
}

.message-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.message-wrapper:hover {
  background: rgba(255, 255, 255, 0.025);
}

.user-message-wrapper {
  background: #444654;
}

.user-message-wrapper:hover {
  background: #4a4b5a;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.user-message {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 16px;
}

.ai-message {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 16px;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}

.user-avatar {
  background: #19c37d;
  color: white;
}

.ai-avatar {
  background: #da70d6;
  color: white;
}

.ai-avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 4px;
}

.message-content {
  flex: 1;
  max-width: calc(100% - 56px);
}

.message-bubble {
  padding: 0;
  border-radius: 0;
  word-wrap: break-word;
  line-height: 1.6;
  position: relative;
  background: transparent;
  border: none;
}

.user-bubble {
  color: #ffffff;
}

.ai-bubble {
  color: #ffffff;
}

.message-text {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 16px;
  line-height: 1.6;
}

/* 打字机光标效果 */
.cursor {
  display: inline-block;
  background-color: #ffffff;
  margin-left: 2px;
  width: 2px;
  animation: blink 1s infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #8e8ea0;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {
  0%,
  80%,
  100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.input-container {
  padding: 32px 0;
  background: #343541;
  position: sticky;
  bottom: 0;
}

.input-wrapper {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0 24px;
}

.input-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #40414f;
  border: 1px solid #565869;
  border-radius: 12px;
  padding: 12px 16px;
  transition: border-color 0.2s;
}

.input-box:focus-within {
  border-color: #10a37f;
  box-shadow: 0 0 0 2px rgba(16, 163, 127, 0.1);
}

/* 流式输出开关样式 */
.stream-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 20px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #565869;
  transition: 0.3s;
  border-radius: 20px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #10a37f;
}

input:checked + .slider:before {
  transform: translateX(20px);
}

.toggle-label {
  font-size: 12px;
  color: #8e8ea0;
  white-space: nowrap;
  user-select: none;
}

.message-input {
  flex: 1;
  min-height: 24px;
  max-height: 120px;
  padding: 0;
  border: none;
  background: transparent;
  outline: none;
  resize: none;
  font-family: inherit;
  font-size: 16px;
  line-height: 1.5;
  color: #ffffff;
}

.message-input::placeholder {
  color: #8e8ea0;
}

.message-input:disabled {
  opacity: 0.6;
}

.send-button {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: #10a37f;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-button:hover:not(:disabled) {
  background: #0d8f69;
}

.send-button:disabled {
  background: #2d2d37;
  cursor: not-allowed;
  opacity: 0.4;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    padding: 16px 20px;
  }

  .title {
    font-size: 20px;
  }

  .message-wrapper {
    padding: 20px;
  }

  .input-wrapper {
    padding: 0 20px;
  }

  .input-container {
    padding: 20px 0;
  }
}

/* 自定义滚动条 */
.messages-container::-webkit-scrollbar {
  width: 8px;
}

.messages-container::-webkit-scrollbar-track {
  background: transparent;
}

.messages-container::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(255, 255, 255, 0.3);
}
</style>

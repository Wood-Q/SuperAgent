<template>
  <div class="chat-container">
    <!-- 头部 -->
    <div class="chat-header">
      <h1 class="title">
        🌙 MoonAgent小助手
      </h1>
      <div class="session-id">
        会话ID: {{ sessionId }}
      </div>
    </div>

    <!-- 消息区域 -->
    <div class="messages-container" ref="messagesContainer">
      <div class="messages-list">
        <!-- 欢迎消息 -->
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="ai-message">
            <div class="message-avatar ai-avatar">🤖</div>
            <div class="message-content">
              <div class="message-bubble ai-bubble">
                您好！我是MoonAgent小助手，有什么可以帮助您的吗？
              </div>
            </div>
          </div>
        </div>

        <!-- 聊天消息 -->
        <div v-for="message in messages" :key="message.id" class="message-item">
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
            <div class="message-avatar ai-avatar">🤖</div>
            <div class="message-content">
              <div class="message-bubble ai-bubble">
                <div v-if="message.content" class="message-text">
                  {{ message.content }}<span v-if="message.isTyping" class="cursor">|</span>
                </div>
                <div v-if="message.isTyping && !message.content" class="typing-indicator">
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

    <!-- 输入区域 -->
    <div class="input-container">
      <div class="input-wrapper">
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
          <svg v-if="!isLoading" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M.5 1.163A1 1 0 0 1 1.97.28l12.868 6.837a1 1 0 0 1 0 1.766L1.969 15.72A1 1 0 0 1 .5 14.836V10.33a1 1 0 0 1 .816-.983L8.5 8 1.316 6.653A1 1 0 0 1 .5 5.67V1.163Z" fill="currentColor"/>
          </svg>
          <div v-else class="loading-spinner"></div>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { streamChatWithModel } from '../utils/api'

export default {
  name: 'ChatRoom',
  setup() {
    const sessionId = ref('')
    const messages = reactive([])
    const inputMessage = ref('')
    const isLoading = ref(false)
    const messagesContainer = ref(null)
    const messageInput = ref(null)

    // 生成会话ID
    const generateSessionId = () => {
      sessionId.value = uuidv4().substring(0, 8)
    }

    // 滚动到底部
    const scrollToBottom = () => {
      nextTick(() => {
        if (messagesContainer.value) {
          const container = messagesContainer.value
          container.scrollTop = container.scrollHeight
        }
      })
    }

    // 平滑滚动到底部
    const smoothScrollToBottom = () => {
      nextTick(() => {
        if (messagesContainer.value) {
          const container = messagesContainer.value
          container.scrollTo({
            top: container.scrollHeight,
            behavior: 'smooth'
          })
        }
      })
    }

    // 调整输入框高度
    const adjustTextareaHeight = () => {
      const textarea = messageInput.value
      if (textarea) {
        textarea.style.height = 'auto'
        textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px'
      }
    }

    // 处理键盘事件
    const handleKeyDown = (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        sendMessage()
      }
    }

    // 发送消息
    const sendMessage = async () => {
      const message = inputMessage.value.trim()
      if (!message || isLoading.value) return

      // 添加用户消息
      const userMessage = {
        id: uuidv4(),
        role: 'user',
        content: message,
        timestamp: new Date()
      }
      messages.push(userMessage)

      // 清空输入框
      inputMessage.value = ''
      adjustTextareaHeight()
      
      // 添加AI消息占位符
      const aiMessage = {
        id: uuidv4(),
        role: 'assistant',
        content: '',
        isTyping: true,
        timestamp: new Date()
      }
      messages.push(aiMessage)

      isLoading.value = true
      scrollToBottom()

      try {
        // 调用流式API
        await streamChatWithModel(message, (chunk) => {
          // 更新AI消息内容 - 逐字符追加
          aiMessage.content += chunk
          aiMessage.isTyping = true // 保持光标显示
          
          // 更频繁的滚动更新，让体验更丝滑
          if (messagesContainer.value) {
            const container = messagesContainer.value
            const isAtBottom = container.scrollHeight - container.clientHeight <= container.scrollTop + 100
            if (isAtBottom) {
              nextTick(() => {
                container.scrollTop = container.scrollHeight
              })
            }
          }
        }, (error) => {
          console.error('Stream error:', error)
          aiMessage.content = '抱歉，发生了错误，请稍后重试。'
          aiMessage.isTyping = false
        })
      } catch (error) {
        console.error('API error:', error)
        aiMessage.content = '抱歉，发生了错误，请稍后重试。'
        aiMessage.isTyping = false
      } finally {
        // 流式传输结束，隐藏光标
        aiMessage.isTyping = false
        isLoading.value = false
        smoothScrollToBottom()
      }
    }

    // 组件挂载
    onMounted(() => {
      generateSessionId()
    })

    return {
      sessionId,
      messages,
      inputMessage,
      isLoading,
      messagesContainer,
      messageInput,
      sendMessage,
      handleKeyDown,
      adjustTextareaHeight
    }
  }
}
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 800px;
  margin: 0 auto;
  background: white;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
}

.chat-header {
  padding: 20px 24px;
  border-bottom: 1px solid #e5e5e5;
  background: white;
  z-index: 10;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.session-id {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 0 24px;
  background: #f9fafb;
  scroll-behavior: smooth;
}

.messages-list {
  padding: 24px 0;
}

.welcome-message {
  margin-bottom: 24px;
}

.message-item {
  margin-bottom: 24px;
  animation: fadeInUp 0.3s ease-out;
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
  align-items: flex-end;
  gap: 12px;
}

.ai-message {
  display: flex;
  justify-content: flex-start;
  align-items: flex-end;
  gap: 12px;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.user-avatar {
  background: #3b82f6;
  color: white;
}

.ai-avatar {
  background: #f3f4f6;
  color: #374151;
}

.message-content {
  max-width: 70%;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 18px;
  word-wrap: break-word;
  line-height: 1.5;
  position: relative;
}

.user-bubble {
  background: #3b82f6;
  color: white;
}

.ai-bubble {
  background: white;
  color: #374151;
  border: 1px solid #e5e7eb;
}

.message-text {
  white-space: pre-wrap;
  word-break: break-word;
}

/* 打字机光标效果 */
.cursor {
  display: inline-block;
  background-color: #374151;
  margin-left: 2px;
  width: 2px;
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 50% {
    opacity: 1;
  }
  51%, 100% {
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
  background: #9ca3af;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.input-container {
  padding: 20px 24px;
  border-top: 1px solid #e5e5e5;
  background: white;
}

.input-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  max-width: 100%;
}

.message-input {
  flex: 1;
  min-height: 44px;
  max-height: 120px;
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 22px;
  outline: none;
  resize: none;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.5;
  transition: border-color 0.2s;
}

.message-input:focus {
  border-color: #3b82f6;
}

.message-input:disabled {
  background: #f9fafb;
  color: #6b7280;
}

.send-button {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  border: none;
  background: #3b82f6;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-button:hover:not(:disabled) {
  background: #2563eb;
  transform: scale(1.05);
}

.send-button:disabled {
  background: #d1d5db;
  cursor: not-allowed;
  transform: none;
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
  .chat-container {
    height: 100vh;
    max-width: 100%;
  }
  
  .chat-header {
    padding: 16px 20px;
  }
  
  .title {
    font-size: 20px;
  }
  
  .messages-container {
    padding: 0 20px;
  }
  
  .input-container {
    padding: 16px 20px;
  }
  
  .message-content {
    max-width: 85%;
  }
}
</style> 
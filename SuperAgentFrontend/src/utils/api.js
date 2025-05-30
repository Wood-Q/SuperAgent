import axios from "axios";

// API基础配置
const API_BASE_URL = "http://localhost:8888/api/chat";

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

/**
 * 流式聊天接口（使用简化的fetch实现）
 * @param {string} userInput - 用户输入
 * @param {function} onChunk - 处理流数据的回调函数
 * @param {function} onError - 错误处理回调函数
 * @returns {Promise}
 */
export const streamChatWithModel = async (userInput, onChunk, onError) => {
  try {
    const response = await fetch(`${API_BASE_URL}/stream`, {
      method: "POST",
      mode: "cors",
      credentials: "omit",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        userInput: userInput,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";
    let currentEvent = null;

    try {
      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          break;
        }

        // 解码数据并添加到缓冲区
        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");

        // 保留最后一行（可能不完整）
        buffer = lines.pop() || "";

        for (const line of lines) {
          const trimmedLine = line.trim();
          if (trimmedLine === "") {
            // 空行表示一个事件结束，处理当前事件
            if (
              currentEvent &&
              currentEvent.event === "message" &&
              currentEvent.data !== undefined
            ) {
              // 处理消息数据
              if (onChunk) {
                onChunk(currentEvent.data);
              }
            } else if (currentEvent && currentEvent.event === "done") {
              // 流结束
              console.log("Stream completed");
              return;
            } else if (currentEvent && currentEvent.event === "error") {
              // 错误处理
              if (onError) {
                onError(new Error(currentEvent.data || "Stream error"));
              }
              return;
            }
            // 重置当前事件
            currentEvent = null;
            continue;
          }

          // 解析SSE格式
          if (trimmedLine.startsWith("event:")) {
            const eventType = trimmedLine.slice(6).trim();
            currentEvent = { event: eventType };
          } else if (trimmedLine.startsWith("data:")) {
            const data = trimmedLine.slice(5); // 保留data:后面的内容，包括空格
            if (currentEvent) {
              currentEvent.data = data;
            }
          }
        }
      }

      // 处理缓冲区中剩余的事件
      if (
        currentEvent &&
        currentEvent.event === "message" &&
        currentEvent.data !== undefined
      ) {
        if (onChunk) {
          onChunk(currentEvent.data);
        }
      }
    } finally {
      reader.releaseLock();
    }
  } catch (error) {
    console.error("Stream error:", error);
    if (onError) {
      onError(error);
    }
    throw error;
  }
};

/**
 * 备用的POST方式调用（非流式）
 * @param {string} userInput - 用户输入
 * @returns {Promise}
 */
export const chatWithModel = async (userInput) => {
  try {
    const response = await apiClient.post("", {
      userInput: userInput,
    });
    return response.data;
  } catch (error) {
    console.error("API Error:", error);
    throw error;
  }
};

// 导出默认配置
export default {
  streamChatWithModel,
  chatWithModel,
  API_BASE_URL,
};

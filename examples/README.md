# 流式聊天 API 使用说明

## 概述

本项目提供了两种聊天 API：

1. 普通聊天 API：`/api/chat` - 同步返回完整响应
2. 流式聊天 API：`/api/chat/stream` - 使用 Server-Sent Events (SSE)实时返回响应

## API 接口

### 1. 普通聊天 API

```bash
POST /api/chat
Content-Type: application/json

{
  "userInput": "你的问题"
}
```

响应：

```json
{
  "message": "完整的回答内容"
}
```

### 2. 流式聊天 API

```bash
POST /api/chat/stream
Content-Type: application/json

{
  "userInput": "你的问题"
}
```

响应格式（SSE 流）：

```
event: message
data: 部分响应内容

event: message
data: 更多响应内容

event: done
data: Stream completed
```

## 使用示例

### 1. 使用 curl 测试

```bash
# 普通聊天
curl -X POST http://localhost:8888/api/chat \
  -H "Content-Type: application/json" \
  -d '{"userInput": "你好"}'

# 流式聊天
curl -N -X POST http://localhost:8888/api/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"userInput": "你好"}'
```

### 2. 使用 HTML 客户端

打开 `examples/stream_client.html` 文件，即可在浏览器中体验流式聊天功能。

### 3. 使用 JavaScript

```javascript
// 创建EventSource连接
async function streamChat(message) {
  const response = await fetch("/api/chat/stream", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      userInput: message,
    }),
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const text = decoder.decode(value);
    console.log("收到:", text);
  }
}
```

## 错误处理

API 会返回以下错误状态码：

- `400 Bad Request` - 请求参数错误或缺少必填参数
- `500 Internal Server Error` - 服务器内部错误

错误响应格式：

```json
{
  "error": "错误描述"
}
```

## 注意事项

1. 确保输入的 `userInput` 不为空
2. 流式 API 使用 SSE 协议，客户端需要支持流式读取
3. 建议设置合理的超时时间，避免长时间等待

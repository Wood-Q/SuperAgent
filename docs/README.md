# SuperAgent 项目文档

SuperAgent 是一个基于 Go 语言开发的 AI 聊天服务框架，使用 CloudWeGo Hertz 作为 Web 框架，支持多种 LLM 模型集成。

## 项目特点

- 基于 CloudWeGo Hertz 框架，高性能、低延迟
- 支持多种 LLM 模型集成（如 Qwen、Ollama 等）
- 模块化设计，易于扩展
- 完整的项目结构和开发规范
- Docker 支持，快速部署

## 项目结构

```
├── api/            # API 接口层，处理 HTTP 请求
├── config/         # 配置文件和配置结构定义
├── constants/      # 常量定义
├── docs/           # 项目文档
├── global/         # 全局变量和配置
├── initialize/     # 初始化相关代码
├── middleware/     # 中间件（认证、日志等）
├── message_model/  # 消息模型定义
├── router/         # 路由定义
├── tools/          # 工具函数
└── main.go         # 程序入口文件
```

## 快速开始

### 环境要求

- Go 1.18+
- Docker & Docker Compose（可选）

### 安装步骤

1. 克隆项目

```bash
git clone https://github.com/Wood-Q/SuperAgent.git
cd SuperAgent
```

2. 配置项目

```bash
cp config-example.yaml config.yaml
# 编辑 config.yaml 文件，配置服务器和 LLM 模型参数
```

3. 运行项目

```bash
go mod tidy
go run main.go
```

### Docker 部署

```bash
docker-compose up -d
```

## 配置说明

配置文件 `config.yaml` 包含以下主要配置项：

```yaml
host: "127.0.0.1" # 服务器主机地址
port: "9090" # 服务器端口
llm:
  base_url: "http://localhost:11434" # LLM 服务地址
  model: "qwen3:8b" # 使用的模型名称
```

## API 文档

### 聊天接口

#### 发送聊天消息

- 请求路径：`/api/chat/send`
- 请求方法：`POST`
- 请求参数：
  ```json
  {
    "message": "string" // 用户输入的消息
  }
  ```
- 响应格式：
  ```json
  {
    "code": 200, // 状态码
    "message": "string", // 响应信息
    "data": {
      "message": "string", // 用户输入的消息
      "response": "string", // AI 的响应
      "createTime": 0 // 消息创建时间
    }
  }
  ```

## 开发指南

### 添加新的 API

1. 在 `api` 目录下创建新的处理函数
2. 在 `router` 目录下注册新的路由
3. 在 `service` 层实现业务逻辑

### 集成新的 LLM 模型

1. 在 `initialize` 目录下添加新的模型初始化代码
2. 在配置文件中添加相应的模型配置
3. 实现模型接口
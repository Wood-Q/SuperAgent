# MoonAgent - 智能 AI 助手系统

<div align="center">

![MoonAgent Logo](https://img.shields.io/badge/MoonAgent-智能助手-blue?style=for-the-badge)

一个基于 Go 和 Vue3 构建的现代化智能 AI 助手系统，集成了 RAG（检索增强生成）、向量数据库、流式对话等先进技术。

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.3+-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

</div>

## 📖 项目介绍

MoonAgent 是一个功能强大的智能 AI 助手系统，旨在提供高质量的对话体验和智能问答服务。系统采用现代化的微服务架构，支持实时流式对话、文档检索、网页浏览等多种功能。

### 🏗️ 系统架构

```
MoonAgent/
├── cmd/                    # 应用程序入口
│   ├── server/            # 服务器启动程序
│   └── di/                # 依赖注入配置
├── internal/              # 内部业务逻辑
│   ├── api/               # API层
│   │   ├── handler/       # 请求处理器
│   │   ├── router/        # 路由配置
│   │   └── middleware/    # 中间件
│   ├── pipeline/          # AI流水线
│   ├── agents/            # AI代理
│   ├── constants/         # 常量定义
│   └── global/            # 全局配置
├── pkg/                   # 公共包
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── embedder/          # 向量嵌入
│   ├── indexer/           # 索引器
│   ├── retriever/         # 检索器
│   ├── splitter/          # 文档分割器
│   ├── vectorDB/          # 向量数据库
│   ├── milvus/            # Milvus集成
│   ├── tools/             # 工具集成
│   └── logger/            # 日志系统
├── SuperAgentFrontend/    # Vue3前端应用
│   ├── src/               # 源代码
│   │   ├── components/    # Vue组件
│   │   ├── utils/         # 工具函数
│   │   └── style.css      # 样式文件
│   └── public/            # 静态资源
├── configs/               # 配置文件
├── assets/                # 资源文件
│   └── documents/         # 文档库
└── tests/                 # 测试文件
```

### 🎯 核心特性

- **🤖 智能对话**: 基于大语言模型的自然语言对话
- **📚 RAG 检索**: 检索增强生成，提供准确的知识问答
- **🌊 流式响应**: 实时流式输出，提升用户体验
- **🔍 文档检索**: 支持文档上传和智能检索
- **🌐 网页浏览**: 集成网页浏览和搜索功能
- **💾 向量存储**: 基于 Milvus 的高性能向量数据库
- **🎨 现代 UI**: 基于 Vue3 的响应式前端界面
- **⚡ 高性能**: 基于 Hertz 框架的高性能后端服务

## 🛠️ 技术栈

### 后端技术

| 技术       | 版本    | 用途             |
| ---------- | ------- | ---------------- |
| **Go**     | 1.24+   | 主要编程语言     |
| **Hertz**  | v0.10.0 | 高性能 HTTP 框架 |
| **Eino**   | v0.3.33 | AI 编排框架      |
| **Milvus** | v2.4.2  | 向量数据库       |
| **Viper**  | v1.20.1 | 配置管理         |
| **Zap**    | v1.27.0 | 结构化日志       |
| **Wire**   | v0.6.0  | 依赖注入         |

### 前端技术

| 技术      | 版本   | 用途        |
| --------- | ------ | ----------- |
| **Vue 3** | 3.3.4+ | 前端框架    |
| **Vite**  | 4.4.5+ | 构建工具    |
| **Axios** | 1.5.0+ | HTTP 客户端 |
| **UUID**  | 9.0.0+ | 会话管理    |

### AI 与集成

- **字节跳动豆包大模型**: 主要的语言模型服务
- **向量嵌入模型**: 文档向量化处理
- **Google 搜索 API**: 网页搜索功能
- **浏览器自动化**: 网页内容抓取
- **SSE 流式传输**: 实时响应流

## 🚀 快速开始

### 环境要求

- **Go**: 1.24 或更高版本
- **Node.js**: 16.0 或更高版本
- **Milvus**: 2.4 或更高版本（可选，用于向量存储）

### 1. 克隆项目

```bash
git clone https://github.com/your-username/MoonAgent.git
cd MoonAgent
```

### 2. 配置环境

启动milvus服务

```bash
cd pkg/vectorDB
docker compose up -d
```

可以通过访问http://127.0.0.1:8000/
进入milvus可视化面板attu

复制配置文件并填写必要信息：

```bash
cp configs/config-example.yaml configs/config.yaml
```

编辑 `configs/config.yaml`：

```yaml
# 服务配置
host: "127.0.0.1"
port: "8888"

# 大模型配置
llm:
  base_url: "your_llm_base_url"
  api_key: "your_api_key"
  model: "your_model_name"

# 向量数据库配置
document:
  addr: "127.0.0.1:19530"
  api_key: "your_embedding_api_key"
  model: "your_embedding_model"

# 浏览器配置
browser:
  api_key: "your_google_api_key"
  search_engine_id: "your_search_engine_id"
```

### 3. 启动后端服务

```bash
# 安装Go依赖
go mod tidy

# 启动服务
go run cmd/server/main.go
```

服务将在 `http://localhost:8888` 启动

### 4. 启动前端应用

```bash
# 进入前端目录
cd SuperAgentFrontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将在 `http://localhost:3000` 启动

### 5. 访问应用

打开浏览器访问 `http://localhost:3000`，开始与 AI 助手对话！

## 📋 API 接口

### 聊天接口

#### 普通聊天

```http
POST /api/chat
Content-Type: application/json

{
  "userInput": "你好，请介绍一下自己"
}
```

#### 流式聊天

```http
POST /api/chat/stream
Content-Type: application/json

{
  "userInput": "请详细解释什么是人工智能"
}
```

响应格式：Server-Sent Events (SSE)

```
event: message
data: 人工智能是...

event: message
data: 一种模拟人类智能...

event: done
data: Stream completed
```

### 文档管理

将需要检索的文档放入 `assets/documents/` 目录。

- [ ] 提供向量化存储的对外 API

## 🧪 开发指南

### 项目结构说明

- **cmd/**: 应用程序入口点，包含服务器启动逻辑
- **internal/**: 内部业务逻辑，不对外暴露
- **pkg/**: 可复用的公共包
- **configs/**: 配置文件目录
- **SuperAgentFrontend/**: Vue3 前端应用

### 添加新功能

1. **添加新的工具**: 在 `pkg/tools/` 目录下创建新的工具文件
2. **扩展 API**: 在 `internal/api/handler/` 添加新的处理器
3. **修改流水线**: 在 `internal/pipeline/` 调整 AI 处理流程

## 🔮 未来规划

### 目标

- [ ] **对话历史**: 实现对话记录和历史查询
- [ ] **API 文档**: 完善 Swagger API 文档
- [ ] **检索策略优化**: 尝试采用多种检索策略进行优化
- [ ] **安全机制引入**: 尝试引入安全作用机制
- [ ] **思维与工具步骤记录**: 引入思维链输出与工具使用的步骤输出记录

## 🤝 贡献指南

我们欢迎所有形式的贡献！

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个星标！**

</div>

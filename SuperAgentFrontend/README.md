# MoonAgent 小助手 - 前端项目

一个基于 Vue3 的智能聊天助手前端应用，提供流式对话体验。

## 功能特性

- 🌙 **现代化聊天界面** - 参考 ChatGPT 设计，简洁美观
- 💬 **实时流式对话** - 使用 SSE 技术实现流式响应
- 📱 **响应式设计** - 完美支持桌面端和移动端
- 🎯 **自动会话管理** - 每次访问自动生成唯一会话 ID
- ⚡ **快速响应** - 优化的用户体验和性能

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 现代前端构建工具
- **Axios** - HTTP 客户端库
- **UUID** - 会话 ID 生成
- **Modern CSS** - 响应式布局和动画

## 项目结构

```
SuperAgentFrontend/
├── src/
│   ├── components/
│   │   └── ChatRoom.vue      # 聊天室主组件
│   ├── utils/
│   │   └── api.js           # API工具类
│   ├── App.vue              # 主应用组件
│   ├── main.js              # 应用入口
│   └── style.css            # 全局样式
├── index.html               # HTML模板
├── package.json             # 项目配置
├── vite.config.js          # Vite配置
└── README.md               # 项目文档
```

## 安装与运行

### 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖

```bash
cd SuperAgentFrontend
npm install
```

### 启动开发服务器

```bash
npm run dev
```

项目将在 http://localhost:3000 启动

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 后端接口

项目需要配合后端 API 使用：

- **接口地址**: `http://localhost:8888/api/chat/stream`
- **请求方式**: POST
- **请求格式**:
  ```json
  {
    "userInput": "用户输入的消息"
  }
  ```
- **响应格式**: Server-Sent Events (SSE)

### 启动后端服务

确保后端 MoonAgent 服务运行在 8888 端口：

```bash
cd ../  # 返回项目根目录
go run main.go  # 启动后端服务
```

## 使用说明

1. **启动应用** - 访问 http://localhost:3000
2. **自动分配会话 ID** - 每次访问自动生成唯一会话标识
3. **发送消息** - 在底部输入框输入问题，按 Enter 或点击发送按钮
4. **实时对话** - AI 助手将以流式方式实时回复
5. **支持多行输入** - 使用 Shift+Enter 换行

## 界面说明

- **顶部标题栏** - 显示应用名称和当前会话 ID
- **消息区域** - 显示对话历史，用户消息在右侧（蓝色），AI 消息在左侧（白色）
- **输入区域** - 底部消息输入框和发送按钮
- **加载状态** - 发送消息时显示加载动画
- **错误处理** - 网络错误时显示友好的错误提示

## 开发说明

### 主要组件

- **ChatRoom.vue** - 聊天室主组件，包含消息显示、输入处理、SSE 通信
- **api.js** - 封装了流式 API 调用逻辑

### 核心功能

- **流式响应处理** - 使用 fetch API 读取 SSE 流数据
- **消息状态管理** - Vue3 Composition API 管理聊天状态
- **自动滚动** - 新消息时自动滚动到底部
- **输入框自适应** - 根据内容自动调整高度

## 浏览器兼容性

- Chrome >= 88
- Firefox >= 84
- Safari >= 14
- Edge >= 88

## 许可证

MIT License

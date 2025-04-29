# SuperAgent 项目文档

## 项目结构

```
├── api/            # API接口层，处理HTTP请求
├── config/         # 配置文件和配置结构定义
├── constants/      # 常量定义
├── docs/           # 项目文档
├── global/         # 全局变量
├── initialize/     # 初始化相关代码
├── middleware/     # 中间件
├── model/          # 数据模型定义
├── router/         # 路由定义
├── service/        # 业务逻辑层
├── utils/          # 工具函数
└── main.go         # 程序入口文件
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
      // 响应数据
      "message": "string", // 用户输入的消息
      "response": "string", // AI的响应
      "createTime": 0 // 消息创建时间
    }
  }
  ```

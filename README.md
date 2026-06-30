# EinoChat

一个基于 [Eino](https://github.com/cloudwego/eino) 框架的 AI 对话服务，主要用练习 **API Key 调用大模型** 和 **Eino 图编排初体验**。

## 项目简介

通过 Eino 的 Graph 编排能力，搭建了一个支持多分支对话、会话持久化的聊天服务。底层调用豆包（Doubao）大模型 API，HTTP 层用 Gin 提供 RESTful 接口，聊天记录通过 GORM + MySQL 存储。

> ⚠️ 这是一个练手项目，主要用于熟悉 Eino 框架的 ChatModel、ChatTemplate、Lambda、Branch 等节点类型，以及 Graph 的组合编排方式。

## 技术栈

- **语言**: Go 1.24
- **AI 框架**: [Eino](https://github.com/cloudwego/eino) — CloudWeGo 出品的 Go AI 应用框架
- **模型接入**: Ark（豆包/火山引擎大模型平台），通过 API Key 调用 `doubao-seed-1-6-250615`
- **HTTP 框架**: Gin
- **数据库**: MySQL + GORM
- **配置**: `.env` 文件（godotenv）

## 快速开始

### 1. 环境要求

- Go 1.24+
- MySQL（本地或远程均可）

### 2. 配置 API Key

在 `Eino_chat/` 目录下创建 `.env` 文件：

```env
ARK_API_KEY=你的火山引擎Ark_API_Key
MYSQL_DSN=用户名:密码@tcp(127.0.0.1:3306)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
```

### 3. 运行

```bash
cd Eino_chat
go mod tidy
go run main.go
```

服务启动后默认监听 `http://localhost:8080`。

### 4. API 调用示例

```bash
# 发送消息
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，介绍一下你自己", "session_id": "test-session"}'
```

## 项目结构

```
einochat/
└── Eino_chat/
    ├── main.go                    # 入口：初始化数据库 → 启动 HTTP 服务
    ├── go.mod / go.sum            # 依赖管理
    ├── .env                       # 环境变量（API Key、数据库连接等）
    └── internal/
        ├── model.go               # ChatModel 节点初始化（Ark API Key 配置）
        ├── prompt.go              # ChatTemplate（提示词模板）
        ├── orchestration.go       # Eino Graph 编排（核心：构建对话流水线）
        ├── lambda_func.go         # Lambda 节点（自定义处理逻辑）
        ├── branch.go              # Branch 节点（对话分支路由）
        ├── chatService.go         # 对话服务层
        ├── api.go                 # Gin HTTP API 路由和处理器
        ├── database.go            # 数据库初始化和连接管理
        └── types.go               # 数据结构定义
```

## 学习要点

### Eino Graph 编排

- **ChatModel 节点**: 封装大模型 API 调用（Ark → 豆包），只需配置 API Key 和 Model 名称即可接入
- **ChatTemplate 节点**: 使用模板格式化 Prompt，实现 System/User 角色消息的组装
- **Lambda 节点**: 编写自定义 Go 函数插入流水线，做消息预处理或后处理
- **Branch 节点**: 根据条件路由到不同分支，实现对话分流

### API Key 调用体验

通过 Ark SDK 调用豆包大模型，核心只需三步：

1. 从火山引擎控制台获取 API Key
2. 配置 `ark.ChatModelConfig{APIKey, Model, Timeout}`
3. `ark.NewChatModel(ctx, config)` 创建模型实例，塞进 Graph 即可使用

Eino 帮你屏蔽了底层的 HTTP 请求、流式响应处理、重试等细节，只需关心 Graph 的节点编排逻辑。

## 作者

练手项目，用于体验 Eino 框架和云端大模型 API Key 调用流程。

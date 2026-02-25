
这也是一份精心设计的 **《Go 语言微服务全栈学习指南 - Frontend to Backend》**。

这份文档不仅包含了项目的基本信息，还特别针对**前端开发者**的视角定制了学习路径和 AI 交互规则。你可以直接保存为项目根目录下的 `LEARNING_GUIDE.md` 或 `.cursorrules`（如果 Cursor 支持读取特定指令文件），并在后续对话中指引 AI。

---

# **Go 微服务学习指南：基于 go-zero-looklook**

## **1. 项目基本信息**

- **项目名称**：`go-zero-learning` (基于 `go-zero-looklook` 复刻精简版)
- **参考源项目**：Mikaelemmmm/go-zero-looklook
- **学员背景**：高级前端开发（熟悉 JS/TS, HTTP, JSON, 组件化思想；对 Go 语言、微服务架构、RPC 通信尚处于入门阶段）。
- **核心技术栈**：
    - **语言**：Go (Golang)
    - **框架**：go-zero (微服务框架)
    - **协议**：HTTP (API), gRPC (RPC)
    - **存储**：MySQL (GORM), Redis
    - **工具**：goctl (代码生成神器), Docker

## **2. 核心教学原则 (AI 指令)**

作为 AI 助教，在指导过程中**必须**严格遵守以下原则：

1. **前端视角类比**：
    - 解释 Go 概念时，尽量使用 TS/JS 做类比（例如：`Struct` ≈ `Interface/Type`，`Goroutine` ≈ `Async/Promise` 但并发模型不同，`Channel` ≈ `Event/Observable`）。
    - 强调类型安全和编译时检查的优势。
2. **循序渐进构建 (Step-by-Step)**：
    - **不要**一次性生成所有代码。
    - 每次只实现一个小的功能闭环（例如：先写 `.api` 定义，再生成代码，最后填空逻辑）。
    - **拒绝黑盒**：在执行 `goctl` 生成代码前，先解释生成了什么，以及为什么要这样生成。
3. **代码风格与规范**：
    - 严格遵守 Go 官方格式 (`gofmt`)。
    - 强制错误处理（Go 的 `if err != nil` 范式），并解释为什么要这么做（对比 try-catch）。
4. **环境优先**：
    - 确保 Docker 容器（MySQL/Redis）先行启动，再运行 Go 代码。
    - 优先使用 `Makefile` 或简单的 Shell 脚本管理命令。

## **3. 项目结构说明**

我们将从空目录开始，最终构建出以下结构（参考 `looklook` 但简化）：

- 
- 
- 
- 

## **4. 学习流程示例**

每个功能模块的开发都遵循 **"定义 -> 生成 -> 实现"** 的标准 `go-zero` 流程：

1. **Define (定义)**: 编写 `.api` (HTTP) 或 `.proto` (RPC) 文件。
    - *AI 任务*: 解释 DSL 语法，对比 RESTful API 设计。
2. **Generate (生成)**: 使用 `goctl` 自动生成 CRUD 代码骨架。
    - *AI 任务*: 执行命令，并概览生成的文件结构。
3. **Implement (实现)**: 在 `internal/logic` 目录中填充业务逻辑。
    - *AI 任务*: 指导编写数据库查询、Redis 缓存逻辑。
4. **Verify (验证)**: 使用 Postman 或 curl 进行测试。

## **5. 分阶段学习计划**

### Phase 1: 环境搭建与 Go 语言热身

- **目标**: 跑通 "Hello World"，理解 Go 的基础语法。
- **任务**:
    1. 安装 Go 环境, Docker, goctl, protoc。
    2. 初始化 Git 仓库与 `go.mod`。
    3. 编写通过 `http/net` 启动的最简单 Web Server。
    4. **知识点**: `go run`, `go mod`, 变量定义, 结构体, 基础 HTTP 处理。

### Phase 2: 单体服务与 API 定义 (User API)

- **目标**: 完成用户注册/登录接口（不涉及 RPC）。
- **任务**:
    1. 编写 `user.api` 文件（定义 `/user/register`, `/user/login`）。
    2. 使用 `goctl api go` 生成代码框架。
    3. 配置 MySQL Docker 容器。
    4. 引入 `GORM` 或 `sqlx` 连接数据库。
    5. 实现注册与登录逻辑（JWT Token 签发）。
    6. **知识点**: `.api` 语法, Handler/Logic 分层, JWT 认证, 数据库操作。

### Phase 3: 微服务拆分与 RPC 通信

- **目标**: 将业务拆分为 "身份服务" 和 "用户中心"，体验微服务调用。
- **任务**:
    1. 创建 `identity` (RPC) 和 `usercenter` (API+RPC) 服务。
    2. 编写 `.proto` 文件定义 RPC 接口。
    3. API 服务通过 `zRPC` 客户端调用 RPC 服务获取数据。
    4. **知识点**: Protobuf, gRPC 原理, 服务发现 (Etcd), `go-zero` 的 zRPC 组件。

### Phase 4: 中间件与进阶组件

- **目标**: 引入 Redis 缓存与全局异常处理。
- **任务**:
    1. 集成 Redis 缓存用户信息。
    2. 实现全局中间件（Middleware）处理权限验证。
    3. 统一错误码设计（参考 `common/xerr`）。
    4. **知识点**: Redis 常用命令, 中间件模式, 错误处理最佳实践。

### Phase 5: 部署与全链路 (Ops)

- **目标**: 像 `looklook` 一样，用 Docker Compose 编排整个系统。
- **任务**:
    1. 编写 `Dockerfile`。
    2. 编写 `docker-compose.yaml` 编排 API, RPC, MySQL, Redis, Etcd。
    3. 前端（可选）调用本地 Docker 启动的接口。

---

### 🚀 如何开始？

请复制以下指令发送给 AI (Cursor)：

> "你好，我已保存了《Go 微服务学习指南》。现在请开始 **Phase 1: 环境搭建与 Go 语言热身**。请检查我的当前目录，帮我初始化 `go.mod`，并创建一个 `main.go` 演示最基本的 Web 服务，顺便给我讲解一下 Go 的 `package` 和 `func` 概念，用 TypeScript 里的模块概念做个类比。"
>
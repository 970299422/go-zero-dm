# Go 微服务学习指南：基于 go-zero-looklook

## 1. 项目基本信息

- **项目名称**：`go-zero-learning` (基于 `go-zero-looklook` 复刻精简版)
- **参考源项目**：[Mikaelemmmm/go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook)
- **学员背景**：高级前端开发（熟悉 JS/TS, HTTP, JSON, 组件化思想；对 Go 语言、微服务架构、RPC 通信尚处于入门阶段）。
- **学习目标**：
    - 精通 Go 语言基础与常用 Web 开发范式。
    - 掌握 `go-zero` 微服务框架的核心概念（API, RPC, Middleware）。
    - 能够从零搭建支持高并发的微服务后端（包含 MySQL, Redis, Etcd）。
    - 理解微服务与前端交互的全链路流程（JWT, 跨域, 错误码）。
- **核心技术栈**：
    - **语言**：Go (Golang)
    - **框架**：go-zero (微服务框架)
    - **协议**：HTTP (API), gRPC (RPC)
    - **存储**：MySQL (GORM), Redis
    - **工具**：goctl (代码生成神器), Docker

## 2. 核心教学原则 (AI 指令)

作为 AI 助教，在指导过程中**必须**严格遵守以下原则：

1. **角色定位**：
    - **AI** 是资深技术专家/导师，负责引导、解释和 Code Review。
    - **User** 是学徒/执行者，负责手动输入代码、运行命令并观察结果。
2. **前端视角类比**：
    - 解释 Go 概念时，**必须**使用 TS/JS 做类比（例如：`Struct` ≈ `Interface`，`Goroutine` ≈ `Async/Promise` 但并发模型不同，`Channel` ≈ `Event/Observable`）。
3. **代码输出策略 (循序渐进与试错)**：
    - **分小步输出**：绝不一次性输出所有代码。每次只实现一个文件或函数。
    - **先给"Low"的写法**：优先提供最直观但可能存在问题（如硬编码、缺乏错误处理、性能差）的代码，让代码先跑起来。
    - **暴露问题**：引导用户运行并发现警告、报错或逻辑缺陷。
    - **解释问题**：以资深开发者的身份，解释为什么刚才的写法在生产环境不可行（如：内存泄漏风险、并发不安全）。
    - **给出优化方案**：在用户理解问题后，再提供 `go-zero` 官方推荐的最佳实践或优化代码。
4. **环境优先**：
    - 确保 Docker 容器（MySQL/Redis）先行启动，再运行 Go 代码。
    - 优先使用 `Makefile` 或简单的 Shell 脚本管理命令。

## 3. 学习流程示例

每个功能点的学习都遵循以下 **"试错 -> 理解 -> 优化"** 循环：

```text
步骤1：AI 给出简单粗暴的实现（可能含 Hardcode 或忽略 Error）
  ↓
步骤2：User 运行代码，观察现象（成功运行，但发现潜在问题或 AI 提示隐患）
  ↓
步骤3：AI 解释问题（例如："这种写法会导致数据库连接池耗尽..."）
  ↓
步骤4：AI 给出优化方案（引入连接池或依赖注入），完成最佳实践
```

## 4. 项目结构说明

我们将从空目录开始，最终构建出以下结构（参考 `looklook` 但简化）：

```text
go-zero-learning/
├── app/
│   ├── identity/      # 认证服务 (业务模块1)
│   │   ├── api/       # HTTP 接口层 (直接面向前端)
│   │   │   ├── internal/
│   │   │   │   ├── config/  # 配置解析
│   │   │   │   ├── handler/ # 路由处理 (Controller)
│   │   │   │   ├── logic/   # 业务逻辑 (Service)
│   │   │   │   ├── svc/     # 依赖注入上下文 (Context)
│   │   │   │   └── types/   # 请求响应结构体 (DTO)
│   │   │   └── identity.api # API 定义文件 (DSL)
│   │   └── rpc/       # RPC 服务层 (微服务内部调用)
│   ├── usercenter/    # 用户中心 (业务模块2)
│   │   ├── api/
│   │   └── rpc/
├── common/            # 通用库 (工具函数, 全局错误码, JWT解析)
├── deploy/            # 部署文件 (Docker Compose, Nginx, Prometheus)
├── go.mod             # 依赖管理
└── go.work            # (可选) 多模块工作区
```

**模块关系**：
- **API 层**：直接对接前端，处理 HTTP 请求，进行参数校验，聚合多个 RPC 服务的结果。
- **RPC 层**：原子化的业务服务，操作数据库/缓存，仅对内提供 gRPC 接口。
- **Model 层**：ORM 映射，直接操作 SQL。

## 5. 分阶段学习计划

### Phase 1: 环境搭建与 Go 语言热身

- **目标**: 跑通 "Hello World"，理解 Go 的基础语法与工程结构。
- **核心知识点**: `go run`, `go mod`, `package` vs `module`, 变量定义, 结构体, 基础 HTTP 处理。
- **任务**:
    1. 安装 Go 环境, Docker, goctl, protoc。
    2. 初始化 Git 仓库与 `go.mod`。
    3. 编写通过 `http/net` 启动的最简单 Web Server（体验手动写路由的繁琐）。

### Phase 2: 单体服务与 API 定义 (User API)

- **目标**: 完成用户注册/登录接口（不涉及 RPC），掌握 `go-zero` 的开发流。
- **核心知识点**: `.api` 语法 (DSL), `goctl` 代码生成原理, Handler/Logic 分层, 依赖注入 (ServiceContext), JWT 认证。
- **任务**:
    1. 编写 `identity.api` 文件。
    2. 使用 `goctl api go` 生成代码框架。
    3. **试错环节**：直接在 Controller 写业务逻辑 -> 发现代码臃肿 -> 重构到 Logic 层。
    4. 实现注册与登录逻辑（JWT Token 签发）。

### Phase 3: 数据库交互与 ORM 集成

- **目标**: 让服务真正持久化数据到 MySQL。
- **核心知识点**: Docker Compose 编排, MySQL 基础, GORM/Sqlx 集成, 配置管理 (Configuration)。
- **任务**:
    1. 使用 Docker Compose 启动 MySQL。
    2. **试错环节**：手动连接 DB -> 遇到硬编码配置问题 -> 引入 `config.yaml`。
    3. **试错环节**：原生 SQL 写法 -> 繁琐且易错 -> 引入 GORM。
    4. 完成用户数据的 CRUD。

### Phase 4: 微服务拆分与 RPC 通信

- **目标**: 将业务拆分为 "身份服务" 和 "用户中心"，体验微服务调用。
- **核心知识点**: Protobuf 语法, gRPC 原理, Etcd 服务发现, `zRPC` 客户端调用, 分布式场景下的调试。
- **任务**:
    1. 将 `identity` 服务拆分为 API + RPC。
    2. 编写 `.proto` 文件定义 RPC 接口。
    3. API 服务通过 `zRPC` 客户端调用 RPC 服务。

### Phase 5: 中间件、缓存与进阶

- **目标**: 这里是生产环境的分水岭，引入 Redis 和全局处理。
- **核心知识点**: Redis 缓存策略, 中间件 (Middleware) 原理, 全局错误码 (xerr), 链路追踪 (Jaeger)。
- **任务**:
    1. 集成 Redis 缓存用户信息。
    2. 实现 Global Custom Error Handler（返回标准 JSON 错误结构）。
    3. 编写自定义中间件处理 RBAC 权限。

### Phase 6: 部署与全链路 (Ops)

- **目标**: 像 `looklook` 一样，用 Docker Compose 编排整个系统。
- **任务**:
    1. 编写 `Dockerfile`。
    2. 编写 `docker-compose.yaml` 编排 API, RPC, MySQL, Redis, Etcd。
    3. 完整联调。

---

### 🚀 如何开始？

请复制以下指令发送给 AI (Cursor)：

> "你好，我已保存了《Go 微服务学习指南》。现在请开始 **Phase 1: 环境搭建与 Go 语言热身**。请检查我的当前目录，帮我初始化 `go.mod`，并创建一个 `main.go` 演示最基本的 Web 服务，顺便给我讲解一下 Go 的 `package` 和 `func` 概念，用 TypeScript 里的模块概念做个类比。"

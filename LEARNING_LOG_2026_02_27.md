# 2026-02-27 学习日志：从单体到微服务的进阶之路

今天的学习密度非常大，不仅完成了数据库的真实接入，还跨越了 Go 语言开发中最陡峭的一段坡——**RPC 微服务拆分**。

## 1. 核心实战：数据库与认证

### 1.1 真实数据库接入 (GORM + SQLite)
不再使用 Mock 数据，而是通过 `GORM` 连接 SQLite 数据库，实现了持久化存储。

- **模型定义**：学会了如何定义 GORM 的 Struct，并使用 Tag (`gorm:"column:..."`) 来映射数据库字段。
- **配置切换**：学会了根据环境使用不同的数据库驱动（Windows 下的纯 Go 驱动 `glebarez/sqlite`）。

### 1.2 密码安全 (bcrypt)
在生产环境中，绝对不能明文存储密码。我们引入了 `golang.org/x/crypto/bcrypt` 包：

- **加密**：`bcrypt.GenerateFromPassword(password, cost)` —— 用于注册时生成哈希。
- **比对**：`bcrypt.CompareHashAndPassword(hashedPassword, password)` —— 用于登录时验证。

### 1.3 JWT (JSON Web Token)
- **颁发**：登录成功后，使用 `golang-jwt/jwt` 生成 Token，包含 `userId` 和过期时间 (`exp`)。
- **解析**：在 API 服务的 Middleware 中解析 Token，将 `userId` 注入到 Context 中，供后续逻辑使用。

---

## 2. Go 语言进阶与填坑

### 2.1 CGO 问题与编译器
- **现象**：在 Windows 下直接使用 `mattn/go-sqlite3` 报错，提示还需要 GCC 编译器，因为它底层是用 C 写的。
- **解决**：切换到 **纯 Go 实现** 的驱动 `github.com/glebarez/sqlite`，完美解决了跨平台编译问题，不需要安装 MinGW 等 C 编译器。

### 2.2 指针与引用
深入探讨了 Go 语言中 `&` (取地址) 和 `*` (取值) 的区别，以及在函数传参时的应用。
- **值传递**：默认行为，复制一份数据。
- **引用传递 (指针)**：传递内存地址，函数内部修改会影响外部变量。

### 2.3 目录结构重构 (Internal 限制)
- **Go 限制**：`internal` 目录下的代码只能被父级和同级目录访问，外部包无法导入。
- **实战**：为了让 `API` 服务和 `RPC` 服务都能复用 `User` 模型，我们将 `backend/app/identity/internal/model` 移动到了公共目录 `backend/common/model`。

---

## 3. 微服务架构：RPC 拆分 (重头戏)

### 3.1 什么是 RPC？
- **概念**：**R**emote **P**rocedure **C**all (远程过程调用)。简单来说，就是像调用本地函数一样调用远程服务器上的函数。
- **位置**：微服务内部通信的基石。相比 HTTP RESTful API，RPC 更高效、更严格（强类型）。

### 3.2 Protobuf (.proto) 的作用
这是 API 和 RPC 服务之间的 **"契约"**。

- **定义服务**：`service Identity { ... }` 定义了有哪些方法可以调。
- **定义消息**：`message GetUserReq { ... }` 定义了入参和出参的数据结构。
- **关键配置**：`option go_package` 决定了生成的 Go 代码的包名。如果设置不当（如包名与文件夹名不一致），会导致 `undefined` 错误，今天我们通过修改 `.proto` 或手动别名解决了这个问题。

### 3.3 架构变更
现在获取用户信息的流程变成了：
1. **浏览器** 发起 HTTP 请求 -> `API 服务 (端口 8888)`
2. API 服务解析 JWT，拿到 `userId`。
3. API 服务发起 **RPC 调用** -> `RPC 服务 (端口 8080)`
4. RPC 服务连接数据库查询用户信息 -> 返回结果给 API
5. API 服务组装 JSON -> 返回给浏览器

---

## 4. 调试与运维技巧

- **goctl 代码生成**：熟练使用 `goctl rpc protoc ...` 一键生成 RPC 代码。
- **端口冲突解决**：遇到端口被占用（如 8888），使用 `netstat -ano | findstr :8888` 查进程 ID，然后 `taskkill /F /PID <pid>` 杀掉进程。
- **多服务启动**：现在的架构需要 **同时启动两个服务**（API 和 RPC）才能正常工作。
- **go mod tidy**：每次引入新依赖或重构目录后，务必执行此命令清理 `go.mod`。

---

## 5. 明日建议

现在的系统架构已经初具雏形，明天我们可以尝试：
1. **全链路测试**：跑通 `API (8888) -> RPC (8080) -> DB` 的完整流程。
2. **服务发现**：目前 RPC 是写死 IP 直连的，后续可以学习使用 Etcd 进行服务发现。
3. **业务完善**：将注册 (`Register`) 逻辑也拆分进 RPC，让 API 层变得更薄。

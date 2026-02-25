# Go 微服务实战：后端实现路线图 (Backend Implementation Roadmap)

这份文档是根据 `frontend` 目录中的 API 定义生成的后端开发大纲。我们将按照此路线图，配合 `frontend` 的需求，一步步实现一个完整的微服务电商后台。

## 1. 项目目标 (Project Goal)
构建一个支持 [Go-Zero-Looklook 前端] 的高性能微服务后端。
最终实现前后端联调，跑通所有业务流程。

## 2. 模块与服务拆分 (Service Splitting)
根据前端 `src/api` 目录下的文件结构，我们将后端拆分为以下核心服务：

| 前端模块 (API) | 对应后端服务 (Service) | 功能描述 |
| :--- | :--- | :--- |
| `auth.ts`, `permission.ts`, `role.ts`, `menu.ts` | **identity** (身份认证服务) | 负责登录认证、JWT 签发、RBAC 权限控制、菜单管理 |
| `user.ts` | **usercenter** (用户中心服务) | 用户信息的增删改查、用户列表 |
| `product.ts`, `category.ts` | **product** (商品服务) | 商品管理、分类管理、库存管理 |
| `order.ts`, `cart.ts` | **order** (订单服务) | 购物车、下单、订单状态流转 |
| `upload.ts` | **common** (公共服务) | 文件上传 (对接 OSS 或本地存储) |
| `dashboard.ts` | **mq / job** (或者聚合服务) | 数据统计 (通常涉及异步计算) |

## 3. 技术栈 (Tech Stack)
- **Framework**: go-zero (v1.8+)
- **Database**: MySQL 8.0
- **Cache**: Redis 7.0
- **Message Queue**: Kafka / Asynq (用于订单超时取消等)
- **Registry**: Etcd (服务发现)
- **Protocol**: HTTP (RESTful API) + gRPC (Internal)

---

## 4. 分阶段实现大纲 (Implementation Outline)

我们将学习过程分为 5 个阶段，每个阶段都有明确的产出。

### 🚀 Phase 1: 基础架构与身份认证 (Identity Service) (当前阶段)
**目标**: 跑通最基础的 "注册/登录" 流程，让前端能登录进系统。
- [ ] **Infrastructure**: 搭建 MySQL, Redis, Etcd 环境 (Docker Compose)。
- [ ] **Identity API**: 定义 `identity.api`。
    - 实现 `/api/auth/login` (登录)
    - 实现 `/api/auth/login` (Admin 登录)
    - 实现 `/api/auth/register` (注册 - 为了方便测试)
    - 实现 `/api/auth/info` (获取当前用户信息)
- [ ] **JWT**: 集成 JWT Middleware，实现 Token 解析与验证。
- [ ] **Database**: 设计 `user` 表，使用 GORM/Sqlx 连接数据库。

### 🛒 Phase 2: 商品与类目管理 (Product Service)
**目标**: 实现商品上架、分类管理，让前端能展示商品列表。
- [ ] **Product RPC**: 定义 `product.proto` (CRUD)。
- [ ] **Product API**: 定义 `product.api`，调用 RPC。
    - `/api/product/list`, `/api/product/detail`
    - `/api/category/list` ...
- [ ] **Upload**: 实现图片上传接口 `/api/upload` (映射 `upload.ts`)。

### 🛍️ Phase 3: 购物车与订单 (Order Service)
**目标**: 实现电商核心交易链路。
- [ ] **Cart**: 实现购物车逻辑 (Redis Hash 结构存储)。
- [ ] **Order RPC**: 定义订单创建、支付(模拟)、发货接口。
- [ ] **Distributed Tx**: 学习 DTM 或 简单的 TCC/SAGA 模式处理库存扣减一致性。

### 🛡️ Phase 4: 权限与菜单 (RBAC)
**目标**: 完善后台管理功能，实现基于角色的动态菜单。
- [ ] **RBAC**: 扩展 `identity` 服务，增加 `Role` 和 `Menu` 表。
- [ ] **Interceptor**: 在 API 网关层实现权限拦截 (Casbin 或 自定义 Middleware)。

### 📊 Phase 5: 仪表盘与部署 (Dashboard & Ops)
**目标**: 生产级特性与数据统计。
- [ ] **Dashboard**: 聚合多个服务的数据，提供 `/api/dashboard` 接口。
- [ ] **Deploy**: 编写 K8s 部署文件或完整 Docker Compose。

---

## 5. 开发规范 (AI Instructions)
*所有后续生成的代码必须遵守：*
1. **API 定义优先**: 必须先写 `.api` 文件，再用 `goctl` 生成 Go 代码。
2. **规范的错误处理**: 使用 `xerr` 包封装错误码，确保前端能收到标准 JSON 错误。
    ```json
    { "code": 1001, "msg": "用户名已存在" }
    ```
3. **分层架构**:
   - `api/internal/logic`: 业务逻辑 (类似于 Controller + Service)
   - `api/internal/svc`: 依赖注入 (ServiceContext)
   - `api/internal/types`: 请求/响应结构体 (DTO)
   - `rpc/model`: 数据库模型 (DAO)

---
**Next Step**: 请告诉 AI "开始 Phase 1，帮我分析 auth.ts 并编写 identity.api"。

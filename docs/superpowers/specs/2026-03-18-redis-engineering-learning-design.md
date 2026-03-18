# 设计：工程实战型 Redis 学习方案（5×10 分钟）

日期：2026-03-18
负责人：Dongm + Codex
状态：草案（待审阅）

## 目标
在今天完成 5 个 10 分钟的工程实战训练，通过一个最小 Demo API 把 Redis 从“能用”学到“可维护”。

## 范围
包含：
- 新建最小 Demo API：`backend/app/demo/api`
- 使用 SQLite 作为数据源
- Cache-Aside 缓存策略 + TTL
- 更新后缓存失效
- Redis 不可用时降级
- 输出简要复盘说明

不包含：
- 复杂中间件体系
- 监控/指标上报
- 高并发压力测试

## 总体方案（5×10 分钟）
1. **最小 Demo API + SQLite 查询**
   - 建立 `GET /demo/item/:id`
   - SQLite 查询并返回
2. **接入 Redis（Cache-Aside）**
   - 先查 Redis
   - Miss 时查 SQLite，再写入 Redis（TTL=60s）
3. **更新与缓存失效**
   - 增加一个最小更新接口：`PUT /demo/item/:id`
   - 请求体：`{ name }`
   - 返回：`{ id, name, updatedAt }`
   - SQLite 更新后删除缓存
4. **降级保护**
   - Redis 超时或不可用时，直接走 SQLite
5. **收尾复盘**
   - 写出 3 条要点：何时命中、何时失效、为何要降级

## Demo API 设计
- 路径：`backend/app/demo/api`
- 端点：`GET /demo/item/:id`
- 返回结构：`{ id, name, updatedAt }`
- 更新端点：`PUT /demo/item/:id`
- 请求体：`{ name }`
- 返回结构：`{ id, name, updatedAt }`

## 缓存策略
- 读：Cache-Aside
  - 先 Redis
  - Miss -> SQLite -> 写 Redis（带 TTL=60s）
- 写：更新后删缓存
- 缓存 Key：`demo:item:{id}`
- 缓存值：JSON 序列化 `{ id, name, updatedAt }`
- 降级：Redis 异常时不影响主链路
  - 超时阈值：100ms
  - 重试：0 次（fail-fast）

## 成功标准
- 可以用一句话解释：Redis 在该接口中起到的作用
- 可以解释“命中、失效、降级”的因果关系
- Demo API 可重复演示
- 可验证演示：三次请求分别触发“命中 / 失效 / 降级”

## 风险与应对
- Redis 不可用：确保降级路径始终可用
- SQLite 初始化：提供最小表结构与示例数据
- 时间限制：每次 10 分钟只做一个明确目标

## SQLite 最小数据定义
- 表：`demo_items`
- 字段：
  - `id` INTEGER PRIMARY KEY
  - `name` TEXT NOT NULL
  - `updated_at` TEXT NOT NULL (ISO8601)
- 示例数据：
  - `id=1, name='apple', updated_at=<now>`
  - `id=2, name='banana', updated_at=<now>`

## 交付物
- Demo API 最小实现
- 说明文档（3 条要点）
- 可复现步骤

# Redis Engineering Learning Demo API Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a minimal Demo API (`/demo/item/:id`) backed by SQLite and Redis cache-aside, with update + cache invalidation and graceful Redis fallback, completed in 5x10-minute steps.

**Architecture:** A new `backend/app/demo/api` service serves a single GET and a PUT update. Reads follow cache-aside with Redis, writes update SQLite then invalidate Redis. Redis failures are fail-fast (100ms timeout, no retries) and fall back to SQLite.

**Tech Stack:** Go, go-zero, SQLite, Redis, Docker (optional for Redis), PowerShell

---

## File/Artifact Map
- Create: `backend/app/demo/api/demo.api` (API DSL)
- Create: `backend/app/demo/api/main.go` (service entry)
- Create: `backend/app/demo/api/internal/config/config.go`
- Create: `backend/app/demo/api/etc/demo.yaml`
- Create: `backend/app/demo/api/internal/handler/getitemhandler.go`
- Create: `backend/app/demo/api/internal/handler/updateitemhandler.go`
- Create: `backend/app/demo/api/internal/logic/getitemlogic.go`
- Create: `backend/app/demo/api/internal/logic/updateitemlogic.go`
- Create: `backend/app/demo/api/internal/svc/servicecontext.go`
- Create: `backend/app/demo/api/internal/types/types.go`
- Create: `backend/app/demo/api/internal/store/sqlite.go`
- Create: `backend/app/demo/api/internal/cache/redis.go`
- Create: `backend/app/demo/api/internal/cache/keys.go`
- Create: `backend/app/demo/api/internal/seed/seed.go`
- Create: `backend/app/demo/api/internal/util/json.go` (json helper)
- Create: `docs/experiments/redis-demo-api/README.md` (3-point summary + repro steps)

---

### Task 1: Define API DSL + Generate Scaffold

**Files:**
- Create: `backend/app/demo/api/demo.api`
- Generate: goctl outputs under `backend/app/demo/api/internal/...`

- [ ] **Step 1: Write demo.api (API DSL)**

```text
info(
    title: "demo"
    desc: "demo api"
    author: "codex"
    version: "v1"
)

type (
    ItemResp {
        id: int64
        name: string
        updatedAt: string
    }

    UpdateItemReq {
        name: string
    }
)

@server(
    prefix: /demo
)
service demo-api {
    @handler GetItem
    get /item/:id returns (ItemResp)

    @handler UpdateItem
    put /item/:id (UpdateItemReq) returns (ItemResp)
}
```

- [ ] **Step 2: Generate scaffold**

Run: `goctl api go -api backend/app/demo/api/demo.api -dir backend/app/demo/api`
Expected: handler/logic/svc/types/config files generated.

- [ ] **Step 3: Commit**

```bash
git add backend/app/demo/api
git commit -m "feat: scaffold demo api"
```

---

### Task 2: SQLite Store + Seed Data

**Files:**
- Create: `backend/app/demo/api/internal/store/sqlite.go`
- Create: `backend/app/demo/api/internal/seed/seed.go`
- Modify: `backend/app/demo/api/internal/svc/servicecontext.go`
- Modify: `backend/app/demo/api/internal/config/config.go`
- Modify: `backend/app/demo/api/etc/demo.yaml`

- [ ] **Step 1: Add DB config**

Update config structs to include:
```go
type Config struct {
    rest.RestConf
    SQLite struct {
        Path string `json:"Path"`
    } `json:"SQLite"`
}
```

Update `demo.yaml`:
```yaml
SQLite:
  Path: ./demo.db
```

- [ ] **Step 2: Implement SQLite store**

Create `sqlite.go`:
```go
package store

import (
    "database/sql"
    _ "modernc.org/sqlite"
    "time"
)

type Item struct {
    ID        int64
    Name      string
    UpdatedAt string
}

type SQLiteStore struct {
    db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, err
    }
    return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Init() error {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS demo_items (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            updated_at TEXT NOT NULL
        );
    `)
    return err
}

func (s *SQLiteStore) SeedIfEmpty() error {
    var count int
    if err := s.db.QueryRow(`SELECT COUNT(1) FROM demo_items`).Scan(&count); err != nil {
        return err
    }
    if count > 0 {
        return nil
    }
    now := time.Now().UTC().Format(time.RFC3339)
    _, err := s.db.Exec(`INSERT INTO demo_items (id, name, updated_at) VALUES (1, 'apple', ?), (2, 'banana', ?)`, now, now)
    return err
}

func (s *SQLiteStore) GetByID(id int64) (*Item, error) {
    var it Item
    err := s.db.QueryRow(`SELECT id, name, updated_at FROM demo_items WHERE id = ?`, id).Scan(&it.ID, &it.Name, &it.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return &it, nil
}

func (s *SQLiteStore) UpdateName(id int64, name string) (*Item, error) {
    now := time.Now().UTC().Format(time.RFC3339)
    _, err := s.db.Exec(`UPDATE demo_items SET name = ?, updated_at = ? WHERE id = ?`, name, now, id)
    if err != nil {
        return nil, err
    }
    return s.GetByID(id)
}
```

- [ ] **Step 3: Seed on service startup**

Create `seed.go`:
```go
package seed

import "go-zero-learning/backend/app/demo/api/internal/store"

type Seeder struct {
    Store *store.SQLiteStore
}

func (s *Seeder) Run() error {
    if err := s.Store.Init(); err != nil {
        return err
    }
    return s.Store.SeedIfEmpty()
}
```

- [ ] **Step 4: Wire in ServiceContext**

Update `servicecontext.go`:
```go
import (
    "go-zero-learning/backend/app/demo/api/internal/store"
    "go-zero-learning/backend/app/demo/api/internal/seed"
)

type ServiceContext struct {
    Config Config
    Store  *store.SQLiteStore
}

func NewServiceContext(c Config) *ServiceContext {
    st, err := store.NewSQLiteStore(c.SQLite.Path)
    if err != nil {
        panic(err)
    }
    if err := (&seed.Seeder{Store: st}).Run(); err != nil {
        panic(err)
    }
    return &ServiceContext{Config: c, Store: st}
}
```

- [ ] **Step 5: Sync deps**

Run: `go mod tidy`
Expected: sqlite/redis deps added as needed.

- [ ] **Step 6: Commit**

```bash
git add backend/app/demo/api/internal backend/app/demo/api/internal/config
git commit -m "feat: add sqlite store and seed"
```

---

### Task 3: Redis Cache Layer (Cache-Aside)

**Files:**
- Create: `backend/app/demo/api/internal/cache/redis.go`
- Create: `backend/app/demo/api/internal/cache/keys.go`
- Modify: `backend/app/demo/api/internal/svc/servicecontext.go`
- Modify: `backend/app/demo/api/internal/config/config.go`
- Modify: `backend/app/demo/api/etc/demo.yaml`
- Modify: `backend/app/demo/api/internal/logic/getitemlogic.go`

- [ ] **Step 1: Add Redis config**

```go
type Config struct {
    rest.RestConf
    SQLite struct {
        Path string `json:"Path"`
    } `json:"SQLite"`
    Redis struct {
        Addr string `json:"Addr"`
        Password string `json:"Password"`
        DB int `json:"DB"`
        TimeoutMs int `json:"TimeoutMs"`
    } `json:"Redis"`
}
```

`demo.yaml`:
```yaml
Redis:
  Addr: 127.0.0.1:6379
  Password: ""
  DB: 0
  TimeoutMs: 100
```

- [ ] **Step 2: Implement Redis client wrapper**

`redis.go`:
```go
package cache

import (
    "context"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    timeout time.Duration
}

func NewRedisCache(addr, password string, db int, timeoutMs int) *RedisCache {
    timeout := time.Duration(timeoutMs) * time.Millisecond
    rdb := redis.NewClient(&redis.Options{
        Addr: addr,
        Password: password,
        DB: db,
        MaxRetries: 0,
        DialTimeout: timeout,
        ReadTimeout: timeout,
        WriteTimeout: timeout,
    })
    return &RedisCache{client: rdb, timeout: time.Duration(timeoutMs) * time.Millisecond}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
    c, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()
    return r.client.Get(c, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
    c, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()
    return r.client.Set(c, key, val, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
    c, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()
    return r.client.Del(c, key).Err()
}
```

`keys.go`:
```go
package cache

import "fmt"

func ItemKey(id int64) string {
    return fmt.Sprintf("demo:item:%d", id)
}
```

- [ ] **Step 3: Wire in ServiceContext**

```go
import "go-zero-learning/backend/app/demo/api/internal/cache"

type ServiceContext struct {
    Config Config
    Store  *store.SQLiteStore
    Cache  *cache.RedisCache
}

func NewServiceContext(c Config) *ServiceContext {
    ...
    rc := cache.NewRedisCache(c.Redis.Addr, c.Redis.Password, c.Redis.DB, c.Redis.TimeoutMs)
    return &ServiceContext{Config: c, Store: st, Cache: rc}
}
```

- [ ] **Step 4: Update GetItem logic (Cache-Aside)**

In `getitemlogic.go`:
```go
key := cache.ItemKey(req.Id)
// 1) try cache
if val, err := l.svcCtx.Cache.Get(l.ctx, key); err == nil {
    var resp types.ItemResp
    if json.Unmarshal([]byte(val), &resp) == nil {
        return &resp, nil
    }
}
// 2) fallback to db
item, err := l.svcCtx.Store.GetByID(req.Id)
if err != nil {
    return nil, err
}
resp := types.ItemResp{Id: item.ID, Name: item.Name, UpdatedAt: item.UpdatedAt}
// 3) set cache with TTL=60s (ignore error)
_ = l.svcCtx.Cache.Set(l.ctx, key, util.MustJSON(resp), 60*time.Second)
return &resp, nil
```

Add helper in `internal/util/json.go`:
```go
package util

import "encoding/json"

func MustJSON(v any) string {
    b, _ := json.Marshal(v)
    return string(b)
}
```

- [ ] **Step 5: Commit**

```bash
git add backend/app/demo/api/internal
git commit -m "feat: add redis cache-aside"
```

---

### Task 4: Update Endpoint + Cache Invalidation

**Files:**
- Modify: `backend/app/demo/api/internal/logic/updateitemlogic.go`

- [ ] **Step 1: Implement UpdateItem logic**

```go
item, err := l.svcCtx.Store.UpdateName(req.Id, req.Name)
if err != nil {
    return nil, err
}
// invalidate cache
_ = l.svcCtx.Cache.Del(l.ctx, cache.ItemKey(req.Id))
resp := types.ItemResp{Id: item.ID, Name: item.Name, UpdatedAt: item.UpdatedAt}
return &resp, nil
```

- [ ] **Step 2: Commit**

```bash
git add backend/app/demo/api/internal/logic/updateitemlogic.go
git commit -m "feat: add update and cache invalidation"
```

---

### Task 5: Graceful Redis Fallback + Doc

**Files:**
- Modify: `backend/app/demo/api/internal/logic/getitemlogic.go`
- Modify: `docs/experiments/redis-demo-api/README.md`

- [ ] **Step 1: Fail-fast behavior**

Ensure cache errors never fail the request:
- Any Redis error should fall through to SQLite
- Cache Set/Del errors are ignored

- [ ] **Step 2: Add doc with 3-point summary + repro**

`docs/experiments/redis-demo-api/README.md`:
```text
1) 命中：第一次 GET -> miss -> DB -> set cache；第二次 GET -> hit
2) 失效：PUT 更新后删除缓存，再 GET -> miss
3) 降级：停止 Redis 后 GET -> 仍返回 DB 数据

Repro:
- start redis (docker)
- run demo api
- curl GET/PUT:
  - curl http://127.0.0.1:8888/demo/item/1
  - curl -X PUT http://127.0.0.1:8888/demo/item/1 -H "Content-Type: application/json" -d "{\"name\":\"pear\"}"
  - curl http://127.0.0.1:8888/demo/item/1
```

- [ ] **Step 3: Commit**

```bash
git add backend/app/demo/api/internal/logic/getitemlogic.go docs/experiments/redis-demo-api/README.md
git commit -m "docs: add redis demo repro"
```

---

## Verification
- `go test ./...`
- Manual curl:
  - `GET /demo/item/1` -> returns item
  - second GET hits cache
  - `PUT /demo/item/1` changes name
  - next GET reflects updated name (cache invalidated)
  - stop redis, GET still returns data

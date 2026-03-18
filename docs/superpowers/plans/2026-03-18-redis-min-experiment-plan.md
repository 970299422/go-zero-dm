# Redis 最小实验实现计划

> **面向自动化执行的代理：** 必须使用 superpowers:subagent-driven-development（如可用）或 superpowers:executing-plans 来执行本计划。步骤使用 `- [ ]` 复选框语法跟踪。

**目标：** 在临时目录中用 Docker + redis-cli 完成一次可随时删除的 Redis 实验，观察 SET/GET/TTL/过期，并产出一句话总结。

**架构：** 单个 Redis 容器对外暴露 localhost:6379，通过 `docker exec` 运行 redis-cli。无应用接入，仅做命令级验证。

**技术栈：** Docker、redis-cli（容器内）、PowerShell

---

## 文件/产物清单
- 新建：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`（观察记录 + 总结）
- 新建：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`（可复用命令清单）

---

### 任务 1：准备临时工作区

**文件：**
- 新建：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`
- 新建：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`

- [ ] **步骤 1：创建临时目录**

运行：`New-Item -ItemType Directory -Force -Path $env:TEMP\redis-min-experiment`
预期：目录存在 `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment`

- [ ] **步骤 2：初始化记录文件**

运行：
```powershell
$readme = "# Redis 最小实验`n`n## 观察记录`n- `n`n## 一句话总结`n- `n"
Set-Content -Path $env:TEMP\redis-min-experiment\README.md -Value $readme -Encoding UTF8
Set-Content -Path $env:TEMP\redis-min-experiment\commands.txt -Value "" -Encoding UTF8
```
预期：`README.md` 与 `commands.txt` 已生成且包含模板内容。

---

### 任务 2：启动 Redis（Docker）

**文件：**
- 修改：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`

- [ ] **步骤 1：启动 Redis 容器**

运行：
```powershell
docker run --name redis-min-exp -p 6379:6379 -d redis:7
```
预期：打印容器 ID；`docker ps` 显示 `redis-min-exp` 运行中。

- [ ] **步骤 2：记录命令**

追加到 `commands.txt`：
```text
docker run --name redis-min-exp -p 6379:6379 -d redis:7
```

---

### 任务 3：最小 Redis 交互（SET/GET/TTL/过期）

**文件：**
- 修改：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`
- 修改：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`

- [ ] **步骤 1：SET 一个 key**

运行：
```powershell
docker exec -it redis-min-exp redis-cli SET user:info:1 "hello"
```
预期：`OK`

- [ ] **步骤 2：GET 读取**

运行：
```powershell
docker exec -it redis-min-exp redis-cli GET user:info:1
```
预期：`"hello"`

- [ ] **步骤 3：设置 TTL**

运行：
```powershell
docker exec -it redis-min-exp redis-cli EXPIRE user:info:1 5
```
预期：`(integer) 1`

- [ ] **步骤 4：观察 TTL 递减**

运行：
```powershell
docker exec -it redis-min-exp redis-cli TTL user:info:1
```
预期：出现一个小整数（如 4、3、2…），重复执行会递减。

- [ ] **步骤 5：验证过期 miss**

运行：
```powershell
Start-Sleep -Seconds 6

docker exec -it redis-min-exp redis-cli GET user:info:1
```
预期：`(nil)`

- [ ] **步骤 6：记录命令与现象**

把上面的 redis-cli 命令追加到 `commands.txt`，并在 `README.md` 的“观察记录”里写下 TTL 递减与过期行为。

---

### 任务 4：清理

**文件：**
- 修改：`C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`

- [ ] **步骤 1：停止并删除容器**

运行：
```powershell
docker stop redis-min-exp

docker rm redis-min-exp
```
预期：容器被删除；`docker ps -a` 不再出现。

- [ ] **步骤 2：写一句话总结**

在 `README.md` 里写一句总结，例如：
“Cache-Aside = 先读缓存，未命中再查源数据并写回，同时设置 TTL，保证缓存失效时主链路不被拖垮。”

---

## 验证
- Task 3 期间 `docker ps` 显示容器运行中。
- 过期前 `GET` 有值，过期后 `GET` 为 `(nil)`。
- `README.md` 包含观察记录 + 一句话总结。

## 备注
- 若 6379 被占用，改用 `-p 6380:6379`，并同步更新相关命令。

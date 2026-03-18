# Redis Minimal Experiment Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Run a disposable Redis experiment in a temporary directory using Docker and redis-cli to observe SET/GET/TTL/expiration and produce a one-sentence takeaway.

**Architecture:** A single Docker Redis container exposes localhost:6379. We interact via redis-cli in the same shell. No application integration; only command-level verification.

**Tech Stack:** Docker, redis-cli (from container), PowerShell

---

## File/Artifact Map
- Create: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md` (notes + takeaway)
- Create: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt` (exact commands run, for replay)

---

### Task 1: Prepare Temporary Workspace

**Files:**
- Create: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`
- Create: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`

- [ ] **Step 1: Create the temp directory**

Run: `New-Item -ItemType Directory -Force -Path $env:TEMP\redis-min-experiment`
Expected: Directory exists at `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment`

- [ ] **Step 2: Seed notes files**

Run:
```powershell
$readme = "# Redis Minimal Experiment`n`n## Observations`n- `n`n## One-sentence takeaway`n- `n"
Set-Content -Path $env:TEMP\redis-min-experiment\README.md -Value $readme -Encoding UTF8
Set-Content -Path $env:TEMP\redis-min-experiment\commands.txt -Value "" -Encoding UTF8
```
Expected: `README.md` and `commands.txt` exist with the stub content.

---

### Task 2: Start Redis (Docker)

**Files:**
- Modify: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`

- [ ] **Step 1: Start Redis container**

Run:
```powershell
docker run --name redis-min-exp -p 6379:6379 -d redis:7
```
Expected: Container id printed. `docker ps` shows `redis-min-exp` running.

- [ ] **Step 2: Record the command**

Append to `commands.txt`:
```text
docker run --name redis-min-exp -p 6379:6379 -d redis:7
```

---

### Task 3: Minimal Redis Interaction (SET/GET/TTL/Expire)

**Files:**
- Modify: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\commands.txt`
- Modify: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`

- [ ] **Step 1: SET a key**

Run:
```powershell
docker exec -it redis-min-exp redis-cli SET user:info:1 "hello"
```
Expected: `OK`

- [ ] **Step 2: GET the key**

Run:
```powershell
docker exec -it redis-min-exp redis-cli GET user:info:1
```
Expected: `"hello"`

- [ ] **Step 3: Set TTL**

Run:
```powershell
docker exec -it redis-min-exp redis-cli EXPIRE user:info:1 5
```
Expected: `(integer) 1`

- [ ] **Step 4: Observe TTL**

Run:
```powershell
docker exec -it redis-min-exp redis-cli TTL user:info:1
```
Expected: A small integer (e.g., 4, 3, 2...) that decreases on repeat.

- [ ] **Step 5: Confirm expiration miss**

Run:
```powershell
Start-Sleep -Seconds 6

docker exec -it redis-min-exp redis-cli GET user:info:1
```
Expected: `(nil)`

- [ ] **Step 6: Record commands and observations**

Append the five redis-cli commands into `commands.txt` and summarize the TTL behavior in `README.md` under Observations.

---

### Task 4: Cleanup

**Files:**
- Modify: `C:\Users\Dongm\AppData\Local\Temp\redis-min-experiment\README.md`

- [ ] **Step 1: Stop and remove container**

Run:
```powershell
docker stop redis-min-exp

docker rm redis-min-exp
```
Expected: Container removed; `docker ps -a` no longer lists it.

- [ ] **Step 2: Write one-sentence takeaway**

Add one concise sentence to `README.md`, for example:
"Cache-Aside = read cache first, miss goes to DB, then write cache with TTL so future reads are fast but never block the main flow when cache fails."

---

## Verification
- `docker ps` shows container running during Task 3.
- `redis-cli GET` returns value before TTL expires, then `(nil)` after expiry.
- `README.md` contains observations + a one-sentence takeaway.

## Notes
- If port 6379 is already in use, re-run Docker with `-p 6380:6379` and replace commands with `-p 6380:6379` and `redis-cli -p 6380` (or `docker exec` stays the same).

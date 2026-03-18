# Design: Redis Minimal Experiment (Docker + redis-cli)

Date: 2026-03-18
Owner: Dongm + Codex
Status: Draft

## Goal
Build a minimal, disposable Redis experiment in a temporary directory to understand:
- Basic read/write
- TTL and expiration
- Cache-aside intuition from direct observation

## Scope
In scope:
- Start Redis via Docker on localhost:6379
- Use redis-cli to run a short command sequence
- Observe TTL countdown and expired miss
- Write one-sentence takeaway

Out of scope:
- Application integration
- Cache consistency strategies
- Multi-key design, clustering, or persistence tuning

## Approach Options
1. Docker + redis-cli commands (recommended)
   - Fast feedback, minimal code
2. Docker + minimal Go script
   - Closer to project stack, slightly more setup
3. Docker + minimal JS script
   - Familiar to frontend, less aligned with Go stack

Chosen: Option 1

## Experiment Flow
1. Create a temporary working directory.
2. Start Redis in Docker (publish 6379).
3. Run these commands:
   - SET a key
   - GET the key
   - EXPIRE the key
   - TTL to observe remaining time
   - Wait for expiry, GET to confirm miss
4. Record one-sentence takeaway about cache-aside and expiration.

## Success Criteria
- Redis is reachable locally.
- TTL decreases and expires as expected.
- You can explain the flow in one sentence without referencing notes.

## Risks / Mitigations
- Docker not available: fall back to local Redis install or ask user to enable Docker.
- Port 6379 in use: bind to an alternate port and update commands.

## Deliverables
- A single command list and the resulting observations.
- One-sentence takeaway.

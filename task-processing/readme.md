# Distributed Task Processing with Kafka

A distributed task processing system where an API server accepts tasks and publishes them to Kafka, while multiple workers consume and process tasks in parallel. Workers auto-register with the API server and are monitored via health checks.

Built to understand how decoupling services through message queues enables independent scaling and fault tolerance in distributed systems.

**WIP:** Worker consumption from Kafka not yet implemented.

## Quick Start

```bash
# Start services
docker-compose up

# Submit a task
curl -X POST http://localhost:8081/api/v1/task \
  -H "Content-Type: application/json" \
  -d '{"name":"test-task","type":"report","body":{"key":"value"}}'

# Start a worker
SECRET_KEY=test go run cmd/workers/main.go

# List tasks
curl http://localhost:8081/api/v1/tasks

# Monitor Kafka
open http://localhost:8080
```

## How It Works

1. Client sends task to API (port 8081)
2. API validates, caches in-memory, publishes to Kafka
3. Workers consume from Kafka `tasks` topic
4. Workers auto-register with API using HMAC-SHA256 authentication
5. API health checks workers every 1 second at `/health`

## API

### POST /api/v1/task
Submit task. Requires `name`, `type`, and `body` fields.

```bash
curl -X POST http://localhost:8081/api/v1/task \
  -H "Content-Type: application/json" \
  -d '{"name":"my-task","type":"report","body":{"param":"value"}}'
```

Response (202):
```json
{"task added with success":"550e8400-e29b-41d4-a716-446655440000"}
```

### GET /api/v1/tasks
List all cached tasks.

### GET /api/v1/task/{taskID}
Get specific task by UUID.

### GET /api/v1/workers
List registered workers with status.

## Running Workers

Workers auto-detect local network IP and use random ports (8000-9000, excluding 8000/8080/8081):

```bash
SECRET_KEY=test go run cmd/workers/main.go
```

Each worker:
- Generates unique UUID
- Gets random available port
- Auto-discovers local IP
- Creates HMAC-SHA256 hash for auth
- Registers with API
- Exposes `/health` endpoint

Run multiple workers:

```bash
# Terminal 1
SECRET_KEY=test go run cmd/workers/main.go

# Terminal 2
SECRET_KEY=test go run cmd/workers/main.go
```

## Health Checks

API polls registered workers every 1 second at `http://{worker-addr}/health`. Workers respond with "healthy".

## Components

- **API Server** (8081): Task submission, worker registry, health monitoring
- **Kafka** (9092): Message broker, `tasks` and `failed_tasks` topics
- **Workers**: Random port assignment, auto-registration, Kafka consumers
- **Kafka UI** (8080): Monitor topics and messages

## Config

- `SECRET_KEY`: Worker auth (set to `test` in docker-compose.yaml)
- Kafka topics: `tasks` (1 partition), `failed_tasks`

## Project Structure

```
cmd/
  main.go           # API server
  workers/main.go   # Worker service
internal/
  server.go         # HTTP handlers
  task.go           # Task model
  task_cache.go     # In-memory storage
  workers/worker.go # Worker registration, health
  publisher/        # Kafka client
pkg/
  hash.go           # HMAC auth
  requests.go       # HTTP utils
```

## Status

**Working:**
- Task validation and submission
- Kafka publishing
- Worker registration with hash auth
- Random port assignment
- Health checks (1s interval)

**Not Implemented:**
- Worker Kafka consumption
- Task status updates

# Aggregation Dashboard

Production-style backend service for collecting, validating, processing, scheduling, and auditing multi-source data ingestion.

---

# Overview

This project implements an Aggregation Dashboard system that acts as a centralized ingestion and ETL platform.

The system is designed to:

- Collect data from multiple sources
- Store raw incoming data safely
- Process and validate data asynchronously
- Prevent duplicate processing (idempotency)
- Run scheduled polling jobs automatically
- Provide immutable audit logs
- Support hot-reload configuration
- Run fully containerized with Docker Compose

The implementation follows production-oriented backend engineering practices including:

- Clean architecture separation
- Repository pattern
- Async background processing
- Concurrency-safe execution
- Idempotent ETL
- Immutable audit logging
- Environment-based configuration
- Dockerized deployment

---

# Architecture

## High-Level Flow

```text
                ┌────────────────────┐
                │   Source A         │
                │ REST API Polling   │
                └─────────┬──────────┘
                          │
                ┌─────────▼──────────┐
                │   Source B         │
                │     Webhook        │
                └─────────┬──────────┘
                          │
                ┌─────────▼──────────┐
                │   Source C         │
                │   File Upload      │
                └─────────┬──────────┘
                          │
                          ▼
                ┌────────────────────┐
                │      raw_data      │
                │    status=PENDING  │
                └─────────┬──────────┘
                          │
                          ▼
                ┌────────────────────┐
                │    ETL Pipeline    │
                │                    │
                │ 1. Validation      │
                │ 2. Normalization   │
                │ 3. Processing      │
                └─────────┬──────────┘
                          │
             ┌────────────┴────────────┐
             ▼                         ▼
┌────────────────────┐     ┌────────────────────┐
│   processed_data   │     │ validation_errors  │
└────────────────────┘     └────────────────────┘
```

---

# Technology Stack

| Component | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Router | Chi |
| Database | PostgreSQL 16 |
| DB Driver | pgx/v5 |
| Scheduler | robfig/cron |
| Containerization | Docker Compose |
| Config Loader | YAML |
| Caching/Queue Ready | Redis |
| Migration Tool | golang-migrate |
| Testing | Go testing + testify |

---

# Why These Technologies?

## Why Go?

Go was chosen because it is highly suitable for:

- concurrent workloads
- background workers
- API services
- ETL pipelines
- infrastructure tooling

Go also provides:

- lightweight goroutines
- fast startup time
- strong standard library
- excellent performance for I/O-heavy services

---

## Why Chi Router?

Chi provides:

- lightweight routing
- idiomatic Go design
- middleware support
- clean API structure

It is commonly used in production Go services.

---

## Why PostgreSQL?

PostgreSQL was chosen because it supports:

- JSONB storage
- transactional consistency
- strong indexing support
- advanced SQL features
- UPSERT (`ON CONFLICT`)

These capabilities are important for:

- idempotent ETL
- audit logging
- flexible payload storage

---

## Why pgx Instead of database/sql?

pgx provides:

- better PostgreSQL-native support
- better performance
- better connection pooling
- richer PostgreSQL features

---

## Why robfig/cron?

robfig/cron is the standard cron scheduler library in the Go ecosystem.

It provides:

- cron expression support
- thread-safe scheduling
- dynamic job management
- production stability

---

# Project Structure

```text
.
├── cmd
│   └── api
│       └── main.go
│
├── internal
│   ├── audit
│   ├── collector
│   ├── config
│   ├── database
│   ├── handler
│   ├── models
│   ├── pipeline
│   ├── repository
│   ├── scheduler
│   ├── service
│   └── utils
│
├── migrations
│
├── Dockerfile
├── docker-compose.prod.yml
├── config.yaml
├── validation_schema.json
└── README.md
```

---

# Features

# Task 1 — Multi-Source Data Collector

## Source A — REST API Polling

Features:

- HTTP polling
- retry mechanism
- `If-Modified-Since` support
- 304 handling
- 5xx retry backoff
- 4xx fail-fast

Retry strategy:

```text
5s → 15s → 30s
```

Why?

Because transient server failures should be retried while client-side errors should fail immediately.

---

## Source B — Webhook Receiver

Features:

- HMAC-SHA256 validation
- timing attack protection using `hmac.Equal()`
- asynchronous processing
- idempotent insertion

Why `hmac.Equal()`?

Using `==` for signature comparison can expose timing attacks.

`hmac.Equal()` performs constant-time comparison and is the recommended secure approach.

---

## Source C — File Upload

Features:

- CSV upload
- JSON upload
- file size validation
- format auto-detection
- streaming-safe handling

File size limit:

```text
10MB
```

Large files are rejected before full memory loading to prevent memory abuse.

---

# Idempotency Strategy

The system prevents duplicate ingestion using:

- unique idempotency keys
- database uniqueness constraints
- deterministic hashing

Formula:

```text
SHA256(source_id + payload)
```

This ensures:

- duplicate webhook deliveries are ignored
- rerunning ETL does not duplicate processed records

---

# Task 2 — ETL Pipeline

The ETL pipeline processes raw incoming data asynchronously.

## Pipeline Stages

```text
raw_data
   ↓
validation
   ↓
normalization
   ↓
processed_data
```

---

## Validation

Validation rules are externalized into:

```text
validation_schema.json
```

Why?

Because validation logic should be configurable without recompiling the service.

Validation checks:

- required fields
- field types
- date formats

Invalid records are stored in:

```text
validation_errors
```

instead of being silently dropped.

---

## Normalization

Normalization includes:

- trimming whitespace
- converting dates to UTC RFC3339
- field renaming

Why normalization?

Because incoming data from multiple sources often has inconsistent formatting.

---

## Async Pipeline Execution

Pipeline execution is asynchronous using goroutines.

API response:

```http
202 Accepted
```

This prevents long-running ETL operations from blocking HTTP requests.

---

## Pipeline Concurrency Protection

The pipeline runner prevents concurrent execution.

Why?

Multiple simultaneous ETL runs can:

- waste resources
- increase database load
- introduce race conditions

Only one pipeline execution is allowed at a time.

---

## ETL Idempotency

Processed records use:

```sql
ON CONFLICT DO NOTHING
```

combined with deterministic idempotency keys.

This guarantees:

- rerunning ETL is safe
- no duplicated processed records

---

# Task 3 — Scheduler + Audit Log + Docker

## Scheduler System

The scheduler automatically polls configured sources.

Configuration is externalized into:

```text
config.yaml
```

Example:

```yaml
polling_sources:
  - source_id: source-alpha
    url: https://httpbin.org/anything
    schedule: "@every 30s"
```

---

## Hot Reload Configuration

The system supports:

```http
POST /config/reload
```

without restarting the application.

Why?

Production systems should minimize downtime.

Hot reload allows operators to:

- add jobs
- change schedules
- update polling targets

without restarting containers.

---

# Audit Logging

All critical operations are stored in:

```text
audit_log
```

Tracked events include:

- pipeline execution
- webhook ingestion
- file uploads
- config reloads

---

## Immutable Audit Log

The application database user cannot:

- DELETE
- UPDATE

records in `audit_log`.

Enforced by PostgreSQL permissions:

```sql
REVOKE ALL ON audit_log FROM app_user;

GRANT SELECT, INSERT
ON audit_log
TO app_user;
```

Why?

Audit systems must be tamper-resistant.

This ensures audit records remain trustworthy even if the application layer is compromised.

---

# Dockerized Deployment

The project uses:

```text
docker-compose.prod.yml
```

Services:

- API
- PostgreSQL
- Redis

---

## Why Multi-Container Architecture?

This approach provides:

- service isolation
- portability
- reproducible environments
- production parity

---

## Mounted Configuration

The following files are mounted as volumes:

- `config.yaml`
- `validation_schema.json`

Why?

So configuration can be updated without rebuilding Docker images.

---

# API Endpoints

## Health

```http
GET /health
```

---

## Webhook

```http
POST /webhooks/{source_id}
```

---

## File Upload

```http
POST /upload
```

---

## Run Pipeline

```http
POST /pipeline/run
```

---

## Pipeline Status

```http
GET /pipeline/status
```

---

## Scheduler Status

```http
GET /scheduler/status
```

---

## Reload Scheduler Config

```http
POST /config/reload
```

---

## Audit Logs

```http
GET /audit-log
```

Query parameters:

```text
?action=PIPELINE_RUN
&start_date=2026-05-17
&end_date=2026-05-18
```

---

# Running Locally

## Requirements

- Go 1.26+
- Docker
- PostgreSQL
- golang-migrate

---

## Install Dependencies

```bash
go mod tidy
```

---

## Run PostgreSQL + Redis

```bash
docker compose up -d
```

---

## Run Migration

```bash
migrate \
  -path migrations \
  -database "postgres://admin:admin123@localhost:5433/aggregation_dashboard?sslmode=disable" \
  up
```

---

## Run API

```bash
air
```

or:

```bash
go run ./cmd/api
```

---

# Running Production Stack

```bash
docker compose -f docker-compose.prod.yml up --build
```

---

# Testing

## Run Unit Tests

```bash
go test ./...
```

---

# Example Test Cases

## Webhook Valid Signature

```bash
curl -X POST http://localhost:8080/webhooks/source-001 \
  -H "X-Signature: sha256=..."
```

---

## Upload CSV

```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@test_data.csv"
```

---

## Run ETL Pipeline

```bash
curl -X POST http://localhost:8080/pipeline/run
```

---

## Check Scheduler Status

```bash
curl http://localhost:8080/scheduler/status
```

---

# Design Decisions

## Why Raw Data Is Stored First?

Incoming data is always stored before processing.

Benefits:

- replay capability
- debugging support
- failure recovery
- auditability

---

## Why Async Processing?

Long-running operations should not block HTTP requests.

This improves:

- responsiveness
- scalability
- user experience

---

## Why Repository Pattern?

Repositories separate:

- business logic
- database logic

Benefits:

- maintainability
- testability
- cleaner architecture

---

## Why Externalized Config?

Keeping schedules and validation rules outside source code allows:

- operational flexibility
- safer updates
- faster iteration

---

## Why Immutable Audit Logs?

Audit logs are only useful if they cannot be tampered with.

Database-level permission enforcement provides stronger guarantees than application-level restrictions.

---

# Future Improvements

Possible future enhancements:

- distributed worker queue
- Redis-backed job processing
- Prometheus metrics
- OpenTelemetry tracing
- dead letter queue
- authentication & RBAC
- Kafka ingestion
- S3 file storage
- Kubernetes deployment
- CI/CD pipeline

---

# Conclusion

This project demonstrates:

- backend architecture design
- asynchronous processing
- ETL implementation
- concurrency handling
- idempotent systems
- scheduler orchestration
- immutable auditing
- production-style Docker deployment

The implementation focuses not only on functionality, but also on:

- operational safety
- maintainability
- observability
- scalability
- production engineering practices


# LogFlow

A small log ingestion pipeline built in Go.

I'm building this project to understand how a backend handles a continuous stream of logs, especially around **concurrency, buffering, database writes, retries, and observability**.

This is a learning-focused project, not a production-ready logging system.

## What it does

- Accepts logs through an HTTP API
- Validates incoming logs
- Scrubs basic PII from messages
- Generates a UUID and timestamp for each log
- Buffers incoming logs using a Go channel
- Processes logs using worker goroutines
- Stores logs in PostgreSQL / Supabase
- Retries failed database inserts
- Keeps failed logs in an in-memory failed queue
- Supports retrieving logs
- Supports filtering logs by service, level, and time range
- Handles a full ingestion queue with backpressure
- Supports graceful worker shutdown

## Current Architecture

```text
HTTP Request
     │
     ▼
Validate + Scrub
     │
     ▼
UUID + Timestamp
     │
     ▼
Go Channel
 (buffer)
     │
     ▼
Worker Goroutines
     │
     ▼
PostgreSQL
     │
     └── retry on failure
```

## What I'm learning

The main purpose of this project is learning by building.

Some of the concepts I'm exploring:

- Go goroutines
- Channels and buffered channels
- Worker pools
- Backpressure
- `select`
- `sync.WaitGroup`
- Graceful shutdown
- PostgreSQL with `pgx`
- Redis
- Docker
- Prometheus and Grafana
- Performance testing

## Roadmap

- [x] HTTP log ingestion
- [x] Input validation
- [x] Basic PII scrubbing
- [x] UUID and timestamp generation
- [x] PostgreSQL storage
- [x] Log retrieval
- [x] Log filtering
- [x] Database retry mechanism
- [x] In-memory failed queue
- [x] Queue-full handling / backpressure
- [x] Multiple workers
- [x] Graceful shutdown
- [ ] Redis
- [ ] Performance benchmarking
- [ ] Prometheus + Grafana
- [ ] Docker
- [ ] Kubernetes
- [ ] CI/CD

## Why I built it

I wanted to understand what actually happens between an incoming HTTP request and a database write when many logs arrive at the same time.

Instead of only reading about goroutines, channels, queues, retries, and backpressure, I'm trying to understand them by building a small system and testing what happens when things go wrong.

---

**Status:** Work in progress 🚧

Built while learning Go and backend systems.
# Experiment 1 — 1,000 Request Workload

## Objective

The objective of Experiment 1 was to observe how LogFlow performs when the benchmark workload is increased from 100 requests to 1,000 requests while keeping the LogFlow worker configuration unchanged.

This experiment also demonstrates that repeated benchmark runs can produce different results, so a single measurement should not automatically be treated as a definitive performance characteristic.

## Configuration

### Benchmark

- Benchmark workers: **5**
- Requests per benchmark worker: **200**
- Total requests: **1,000**

### LogFlow

- LogFlow processing workers: **3**
- Redis queue: `logQueue`
- Database: PostgreSQL / Supabase

The benchmark sends HTTP `POST /logs` requests to LogFlow.

## Workload

Each benchmark worker sends 200 HTTP requests.

Therefore:

```text
5 benchmark workers × 200 requests
= 1,000 total requests
```

The benchmark measures the time required for the HTTP requests to complete.

Because LogFlow is asynchronous, benchmark completion represents HTTP-side acceptance into the ingestion pipeline rather than confirmation that every request has already been persisted to PostgreSQL.

## Run 1

```text
Requests:       1,000
Elapsed time:   512.81 ms
Throughput:     1,950.04 requests/sec
```

## Run 2

```text
Requests:       1,000
Elapsed time:   374.88 ms
Throughput:     2,667.54 requests/sec
```

## Observations

The two runs used the same benchmark configuration but produced different throughput values.

```text
Run 1 → 1,950.04 req/s
Run 2 → 2,667.54 req/s
```

This shows that benchmark results can vary between runs even when the configured workload remains unchanged.

The difference means that one benchmark run should not be treated as the definitive throughput of the system.

## Interpretation

Experiment 1 demonstrates that LogFlow can accept a workload of 1,000 HTTP requests while continuing to process logs asynchronously through Redis and the worker layer.

The benchmark measures HTTP request completion, not complete end-to-end persistence latency.

Therefore, the measured throughput should be interpreted as:

> Aggregate HTTP acceptance throughput observed under this specific workload and configuration.

It should **not** be interpreted as the maximum capacity of LogFlow.

## What We Learned

1. Increasing the workload from 100 to 1,000 requests is handled successfully.
2. Benchmark throughput can vary significantly between repeated runs.
3. HTTP acceptance and PostgreSQL persistence are separate stages in the asynchronous pipeline.
4. More controlled and repeated experiments are necessary before making strong performance claims.

## Limitations

This experiment does not identify the bottleneck in the system.

The benchmark also does not currently measure:

- p50 latency
- p95 latency
- p99 latency
- complete end-to-end persistence latency
- maximum sustainable throughput

Further experiments are required to investigate these measurements.
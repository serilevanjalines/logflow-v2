# Experiment 2 — 10,000 Request Workload

## Objective

The objective of Experiment 2 was to increase the benchmark workload substantially and observe how LogFlow behaves when processing 10,000 HTTP requests.

The LogFlow worker configuration remained unchanged from the previous experiment.

## Configuration

### Benchmark

- Benchmark workers: **5**
- Requests per benchmark worker: **2,000**
- Total requests: **10,000**

### LogFlow

- LogFlow processing workers: **3**
- Redis queue: `logQueue`
- Database: PostgreSQL / Supabase

Therefore:

```text
5 benchmark workers × 2,000 requests
= 10,000 total requests
```

## Result

```text
Requests:       10,000
Elapsed time:   3.864373948 seconds
Throughput:     2,587.74 requests/sec
```

## Initial Redis Observation

Immediately after the benchmark completed, Redis still contained a backlog:

```text
logQueue:       2,573
failedQueue:    1
```

At first, this could appear to indicate that some logs had not been processed.

However, LogFlow's architecture is asynchronous.

The benchmark measures the completion of the HTTP requests, while the LogFlow workers can continue processing the Redis queue after the benchmark itself has finished.

## Queue Drain

After allowing the workers to continue processing the backlog, the logs were eventually persisted to Supabase.

The database reached:

```text
12,100 benchmark logs
```

This was consistent with:

```text
2,100 previous benchmark logs
+
10,000 Experiment 2 logs
=
12,100 logs
```

Therefore, the initial Redis backlog did not represent permanently lost logs.

## Key Observation

This experiment clearly demonstrated the separation between:

```text
HTTP acceptance
```

and:

```text
complete processing / persistence
```

The benchmark can finish while LogFlow still has work remaining in Redis.

The pipeline behaves approximately as:

```text
Benchmark
    ↓
HTTP POST /logs
    ↓
LogFlow accepts request
    ↓
Redis logQueue
    ↓
Workers
    ↓
PostgreSQL
```

The benchmark timer stops before the entire pipeline necessarily completes.

## Interpretation

Experiment 2 demonstrates that Redis can absorb a burst of incoming work while the LogFlow workers process that work asynchronously.

The observed throughput of:

```text
2,587.74 requests/sec
```

represents the HTTP-side benchmark result for this particular workload and configuration.

It does **not** mean that PostgreSQL persisted all 10,000 logs within 3.86 seconds.

The remaining Redis backlog demonstrates why HTTP throughput and end-to-end processing throughput must be measured separately.

## What We Learned

1. LogFlow successfully accepted a 10,000-request workload.
2. Redis temporarily accumulated a substantial backlog.
3. The HTTP benchmark completed before all queued logs were processed.
4. Workers continued processing after the benchmark completed.
5. The expected logs eventually reached PostgreSQL.
6. An asynchronous ingestion pipeline requires separate measurements for request acceptance and complete processing.

## Limitations

This experiment does not tell us:

- how long the final log took to reach PostgreSQL
- the p50/p95/p99 request latency
- the p50/p95/p99 persistence latency
- the maximum sustainable ingestion rate
- the exact bottleneck responsible for queue buildup

These require additional measurements.
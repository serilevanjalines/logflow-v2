# Experiment 4 — HTTP Ingestion Latency Distribution

## Objective

The objective of Experiment 4 was to measure the **latency distribution of individual HTTP ingestion requests** in LogFlow.

Previous experiments primarily measured aggregate throughput:

```text
total requests / total elapsed time
```

Experiment 4 additionally records the duration of each individual HTTP request and calculates:

- p50 latency
- p95 latency
- p99 latency

The purpose is to understand how request latency is distributed rather than relying only on aggregate throughput.

---

## What Is Being Measured?

LogFlow is an asynchronous ingestion pipeline.

The HTTP request is completed after the log has been accepted into Redis. The HTTP handler does not wait for the worker to insert the log into PostgreSQL.

Therefore, Experiment 4 measures:

> **HTTP ingestion latency: the time between sending a `POST /logs` request and receiving LogFlow's response after the log has been accepted into Redis.**

It does **not** measure end-to-end persistence latency.

The measured path is:

```text
Benchmark
    ↓
POST /logs
    ↓
LogFlow HTTP handler
    ↓
Validation
    ↓
PII scrubbing
    ↓
UUID + timestamp generation
    ↓
Redis
    ↓
HTTP response
    ↓
Benchmark
```

PostgreSQL processing happens asynchronously after this point:

```text
Redis
    ↓
Worker
    ↓
PostgreSQL
```

---

## Configuration

### Benchmark

- Benchmark workers: **5**
- Requests per benchmark worker: **2,000**
- Total requests: **10,000**

### LogFlow

- LogFlow processing workers: **3**
- Redis queue: `logQueue`
- Database: PostgreSQL / Supabase

The workload was:

```text
5 benchmark workers × 2,000 requests
= 10,000 total requests
```

---

## Measurement Method

For every HTTP request, the benchmark records the time immediately before the request is sent and immediately after the HTTP response is received.

Conceptually:

```text
start timer
    ↓
HTTP POST /logs
    ↓
receive HTTP response
    ↓
stop timer
    ↓
store request duration
```

The benchmark therefore produces an individual latency measurement for each successful request.

After all benchmark workers finish:

```text
Individual request latencies
            ↓
          sort
            ↓
     ┌──────┼──────┐
     ↓      ↓      ↓
    p50    p95    p99
```

The latency measurements are sorted from smallest to largest before calculating the percentile values.

---

## Results

The experiment completed all 10,000 requests.

```text
Requests:        10,000
Total elapsed:   5.587057108 seconds
Throughput:      1,789.85 requests/sec
```

### Latency Distribution

```text
p50 latency:     2.397343 ms
p95 latency:     4.739370 ms
p99 latency:     7.637554 ms
```

| Metric | Result |
|---|---:|
| Total requests | 10,000 |
| Benchmark workers | 5 |
| Requests / worker | 2,000 |
| LogFlow workers | 3 |
| Total elapsed time | 5.587057108 s |
| Throughput | 1,789.85 req/s |
| p50 latency | 2.397343 ms |
| p95 latency | 4.739370 ms |
| p99 latency | 7.637554 ms |

---

## Understanding the Results

### p50 — 2.397343 ms

The p50 latency represents the middle of the measured latency distribution.

Approximately half of the measured requests completed at or below about:

```text
2.40 ms
```

This gives us an indication of the latency experienced by a typical request.

---

### p95 — 4.739370 ms

The p95 latency was:

```text
4.74 ms
```

Approximately 95% of the measured requests completed at or below this latency.

This gives us information about the slower portion of the request distribution.

---

### p99 — 7.637554 ms

The p99 latency was:

```text
7.64 ms
```

Approximately 99% of the measured requests completed at or below this latency.

This is useful for observing **tail latency** — the behavior of the slower requests that would not be visible from the median alone.

---

## Latency Distribution

The relationship between the measured percentiles is:

```text
Fastest
   │
   ├──────────── p50 ≈ 2.40 ms
   │
   ├──────────────────── p95 ≈ 4.74 ms
   │
   ├──────────────────────────── p99 ≈ 7.64 ms
   │
Slowest
```

The p99 latency is higher than the p50 latency, which shows that request latency is not identical for every request.

Some requests take longer than the typical request.

---

## Comparison With Previous Experiments

Experiments 0–3 primarily focused on aggregate throughput.

Experiment 4 adds a different perspective by measuring individual request latency.

```text
Experiments 0–3
        ↓
Aggregate throughput
        ↓
"How many requests can the benchmark complete per second?"

Experiment 4
        ↓
Per-request latency distribution
        ↓
"How long does an individual HTTP ingestion request take?"
```

This distinction is important because throughput and latency describe different aspects of system behavior.

---

## Important Comparison Limitation

The throughput observed in Experiment 4 was:

```text
1,789.85 requests/sec
```

This is lower than the throughput observed in Experiment 2:

```text
2,587.74 requests/sec
```

and Experiment 3:

```text
2,742.34 requests/sec
```

However, this experiment should **not** be used to conclude that adding latency measurement caused LogFlow's throughput to decrease.

Experiment 4 changed the benchmark implementation by recording individual request durations and sending those measurements through an additional channel.

Therefore, the Experiment 4 throughput result is not a perfectly controlled apples-to-apples comparison with Experiments 2 and 3.

The primary purpose of Experiment 4 is the latency distribution.

---

## What This Experiment Tells Us

Experiment 4 established that:

1. LogFlow successfully handled the 10,000-request workload.
2. Individual HTTP ingestion request latency can be measured.
3. The observed p50 latency was approximately **2.40 ms**.
4. The observed p95 latency was approximately **4.74 ms**.
5. The observed p99 latency was approximately **7.64 ms**.
6. Request latency varies across the workload.
7. Aggregate throughput alone does not describe the complete latency behavior of the system.

---

## What This Experiment Does NOT Tell Us

Experiment 4 does not tell us:

- how long a log takes to reach PostgreSQL
- PostgreSQL insertion latency
- Redis queue wait time
- worker processing latency
- end-to-end persistence latency
- the maximum sustainable throughput of LogFlow
- the exact bottleneck in the system

The HTTP response occurs before the worker necessarily completes PostgreSQL persistence.

Therefore:

```text
HTTP ingestion latency
        ≠
PostgreSQL persistence latency
```

---

## Key Engineering Lesson

An asynchronous system can have a fast HTTP response while work continues in the background.

For LogFlow:

```text
Client
  ↓
HTTP ingestion
  ↓
Redis
  ↓
HTTP response
  ↓
Client considers request accepted

             meanwhile...

Redis
  ↓
Worker
  ↓
PostgreSQL
```

Therefore, measuring only HTTP latency would not be enough to understand the complete performance of LogFlow.

Experiment 4 gives us visibility into the **front half of the pipeline**.

A future experiment could investigate the latency of the complete path from ingestion to successful PostgreSQL persistence.

---

## Conclusion

Experiment 4 extended the previous throughput experiments by measuring the distribution of individual HTTP ingestion request latencies.

Under the tested configuration:

```text
5 benchmark workers
2,000 requests / worker
10,000 total requests
3 LogFlow workers
```

the benchmark observed:

```text
Throughput:   1,789.85 requests/sec
p50:          2.397343 ms
p95:          4.739370 ms
p99:          7.637554 ms
```

These results describe **HTTP ingestion performance up to Redis acceptance**, not complete database persistence.

The experiment provides a foundation for more detailed performance analysis and future observability work.
# Experiment 3 — Increasing LogFlow Workers

## Objective

The objective of Experiment 3 was to investigate whether increasing the number of LogFlow processing workers improves throughput.

This was designed as a controlled comparison against Experiment 2.

The benchmark workload remained the same, while only the number of LogFlow workers was changed.

## Configuration

### Benchmark

The benchmark configuration remained unchanged:

- Benchmark workers: **5**
- Requests per benchmark worker: **2,000**
- Total requests: **10,000**

### LogFlow

Experiment 2:

```text
3 processing workers
```

Experiment 3:

```text
6 processing workers
```

The worker count was therefore doubled while the benchmark workload remained fixed.

## Workload

```text
5 benchmark workers × 2,000 requests
= 10,000 requests
```

## Result

Experiment 3 produced:

```text
Requests:       10,000
Elapsed time:   3.646517961 seconds
Throughput:     2,742.34 requests/sec
```

## Comparison With Experiment 2

| Configuration | Exp 2 | Exp 3 |
|---|---:|---:|
| Benchmark workers | 5 | 5 |
| Requests | 10,000 | 10,000 |
| LogFlow workers | 3 | 6 |
| Throughput | 2,587.74 req/s | 2,742.34 req/s |

The observed throughput increased from:

```text
2,587.74 req/s
```

to:

```text
2,742.34 req/s
```

This is an increase of approximately **6%**.

## Observation

Doubling the number of LogFlow workers did not double throughput.

The result therefore suggests that worker count is not the only factor determining the observed performance of the system.

However, this experiment alone does not identify the exact bottleneck.

## Possible Limiting Factors

Potential factors include:

- PostgreSQL
- PostgreSQL connection pool
- Redis
- HTTP/client overhead
- worker scheduling
- database latency
- benchmark design

These are possibilities, not confirmed bottlenecks.

Additional instrumentation is required before identifying one of them as the actual limiting factor.

## Interpretation

Experiment 3 provides an important result about concurrency.

Adding more workers can increase processing concurrency, but additional workers do not necessarily produce proportional throughput improvements.

In this experiment:

```text
3 workers
    ↓
2,587.74 req/s

6 workers
    ↓
2,742.34 req/s
```

The increase was modest compared with the 2× increase in worker count.

This indicates that another part of the pipeline may be limiting throughput, or that the benchmark is measuring a stage that is not strongly affected by the number of processing workers.

## What We Learned

1. Increasing worker count can improve observed throughput.
2. Doubling workers does not guarantee a 2× throughput improvement.
3. Worker concurrency is only one component of overall system performance.
4. The true bottleneck has not yet been identified.
5. Performance experiments should change one important variable at a time when possible.
6. More instrumentation is needed to understand where time is being spent.

## Limitations

This experiment does not prove that PostgreSQL, Redis, the connection pool, or any other component is the bottleneck.

It also does not measure:

- p50 latency
- p95 latency
- p99 latency
- end-to-end persistence latency
- queue wait time
- database insertion latency
- worker utilization

Therefore, the result should be treated as an observation rather than a final performance conclusion.

## Next Step

The next useful measurement is latency.

Instead of measuring only aggregate HTTP throughput, future experiments should measure latency distributions such as:

```text
p50
p95
p99
```

This can help determine how individual requests behave under load and provide a stronger basis for understanding LogFlow's performance.
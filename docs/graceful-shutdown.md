NORMAL
────────────────────────────

HTTP request
     ↓
PII scrub + UUID + timestamp
     ↓
logQueue
     ↓
workers
     ↓
PostgreSQL


SHUTDOWN
────────────────────────────

shutdown requested
     ↓
HTTP server stops accepting NEW requests
     ↓
already-accepted HTTP requests finish
     ↓
logQueue stops receiving
     ↓
workers finish their current DB operation
     ↓
workers drain remaining logQueue
     ↓
queue empty
     ↓
workers exit
     ↓
WaitGroup = all workers done
     ↓
pool.Close()
     ↓
program exits



FINAL EXECUTION:


MAIN GOROUTINE

Ctrl+C
   ↓
signalChan receives interrupt
   ↓
<-signalChan completes
   ↓
close(shutdown)
   ↓
wg.Wait()
   │
   │ waiting for workers
   ↓
   └───────────────┐
                   │
              WORKERS
                   │
                select
               /     \
              ↓       ↓
          logQueue  shutdown
              │       │
           process   exit
              │       │
              └───────┘
                   ↓
               wg.Done()
                   ↓
              wg.Wait() completes
                   ↓
               main returns
                   ↓
             pool.Close()
                   ↓
               program exits
1 .
redis-cli ping

pong

2.
╭─LogFlow on 🌊 main 
╰─❯ redis-cli
127.0.0.1:6379> 

3.
127.0.0.1:6379> 

╭─LogFlow on 🌊 main 
╰─❯ redis-cli
127.0.0.1:6379> SET test "hello"
OK
127.0.0.1:6379> GET test
"hello"
127.0.0.1:6379> 

4. LPUSH logQueue "log1" - Push this value onto the left side of the List.
   LPUSH logQueue "log2"

   LRANGE logQueue 0 -1 


5.

BRPop(
    context.Background(),
    0,
    "logQueue",
)

Break it apart:
BRPop(
   context.Background(),   ← context for this operation
   0,                      ← timeout
   "logQueue"              ← Redis list to consume from
)

The 0 is important.

It means:
Wait indefinitely.


┌─────────────────────┐
│     LogFlow         │
│      Go app         │
│                     │
│   redis.Client      │
└──────────┬──────────┘
           │
           │ Redis protocol
           ▼
┌─────────────────────┐
│   Redis / Valkey     │
│   localhost:6379     │
│                     │
│   data lives here   │
└─────────────────────┘


encoding part is


HTTP
  │
  │ JSON
  ▼
json.Decode
  │
  ▼
Go Log struct
  │
  │ json.Marshal
  ▼
JSON
  │
  ▼
Redis
  │
  │ RPOP / BRPOP
  ▼
JSON
  │
  │ json.Unmarshal
  ▼
Go Log struct
  │
  ▼
processLog()



And if PostgreSQL fails:


Go Log struct
      │
      │ json.Marshal
      ▼
   JSON
      │
      ▼
Redis failedQueue

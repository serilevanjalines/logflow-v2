1.Incoming path
JSON → Go
HTTP JSON
   ↓ Decode
Go struct

2.Put into Redis
Go → JSON
Go struct
   ↓ Marshal
JSON
   ↓
Redis

3.Worker receives it
JSON → Go
Redis
   ↓
JSON
   ↓ Unmarshal
Go struct

4.Failed log
Go → JSON again
Go struct
   ↓ Marshal
JSON
   ↓
failedQueue in Redis
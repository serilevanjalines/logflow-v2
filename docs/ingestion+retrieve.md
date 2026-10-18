A concurrent log-ingestion and retrieval service.

📥 Ingestion


Client
  ↓
POST /logs
  ↓
validate
  ↓
PII scrub
  ↓
ID + timestamp
  ↓
queue
  ↓
worker pool
  ↓
Supabase/PostgreSQL


📤 Retrieval


Client
  ↓
GET /logs
GET /logs/{id}
GET /logs?service=...
GET /logs?from=...&to=...
  ↓
PostgreSQL
  ↓
JSON

And both paths can operate concurrently.
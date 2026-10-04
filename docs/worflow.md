1.

Client
  │
  │ POST /logs
  ▼
Go HTTP Handler
  │
  ├── validate
  ├── scrub PII
  ├── generate UUID
  ├── generate timestamp
  │
  ▼
conn.Exec(...)
  │
  │ SQL + parameters
  ▼
Supabase PostgreSQL
  │
  ▼
logs table
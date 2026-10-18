POST /logs
      │
      ▼
HTTP handler goroutine
      │
      ├── decode JSON
      ├── validate
      ├── scrub PII
      ├── generate ID
      └── generate timestamp
              │
              ▼
       logQueue <- incomingLog
              │
              ▼
          CHANNEL
              │
              │ worker receives
              ▼
          logWorker()
              │
              ▼
        PostgreSQL INSERT

And importantly, the worker doesn't get created for every POST.
We create the worker once when the application starts:
go logWorker(conn)
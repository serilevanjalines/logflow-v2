Step 1 — Change the type

Your worker currently receives:
func logWorker(conn *pgx.Conn)

Conceptually, we're changing that to:
func logWorker(pool *pgxpool.Pool)

So the worker no longer says:
"Give me this particular database connection."

It says:
"Give me access to the database connection pool."


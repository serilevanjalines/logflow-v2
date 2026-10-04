Go API 1 ──┐
Go API 2 ──┤
Go API 3 ──┼──→ Connection Pooler ──→ PostgreSQL
Go API 4 ──┤
Go API 5 ──┘


go get github.com/jackc/pgx/v5 
We're choosing pgx because it's a Go PostgreSQL driver designed specifically for PostgreSQL.

go get github.com/joho/godotenv: dotenv



1. context.Background()

basically means:

"Start this operation with a basic, empty context that isn't currently cancelled."

We're using it because pgx.Connect expects a context.


2. ctx, cancel := context.WithTimeout(
    context.Background(),
    5*time.Second,
)
defer cancel()

Then PostgreSQL gets only 5 seconds to respond.

But don't worry about that yet.
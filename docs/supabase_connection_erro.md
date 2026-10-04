There are two places our environment variable can come from:

1. Shell environment
2. .env file

And the shell environment takes precedence.

So for LogFlow development, we can simply do:

unset DATABASE_URL
go run .

when necessary.
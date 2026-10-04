There are two places our environment variable can come from:

1. Shell environment
2. .env file

And the shell environment takes precedence.

So for LogFlow development, we can simply do:

unset DATABASE_URL
go run .

when necessary.


we're defer rows.close() because it contains
the streams of input and so , we need to release those resources to make it free

“I acquired/started using a resource → I should make sure it gets released.”
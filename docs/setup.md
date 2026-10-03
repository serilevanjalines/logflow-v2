1.check go version

go version
go version go1.26.8-X:nodwarf5 linux/amd64

2.go mod init LogFlow

That gives our Go code an import path and tells Go:
"This directory is the root of the logflow module."
Later, when we add dependencies such as a PostgreSQL driver, Go will record them here.

it tracks the dependencies in our project

go.sum :

and go.sum will contain checksums that help Go verify what it downloaded.


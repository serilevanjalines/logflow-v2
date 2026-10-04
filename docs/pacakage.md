The special thing about main is that Go treats it as an executable program when it contains:

Our whole LogFlow application is currently one main package:

main.go
db.go
pii.go
    ↓
package main
    ↓
one LogFlow program



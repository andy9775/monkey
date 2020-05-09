.PHONY: test repl bench

all: test

test:	
	@go test ./...

repl:
	@go run main.go repl

bench: 
	@go run main.go bench
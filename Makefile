.PHONY: test repl

all: test

test:	
	@go test lexer/*.go
repl:
	@go run main.go
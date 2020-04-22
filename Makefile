.PHONY: test repl

all: test

test:	
	@go test lexer/*.go
	@go test parser/*.go

repl:
	@go run main.go
.PHONY: test repl

all: test

test:	
	@go test lexer/*.go
	@go test parser/*.go
	@go test ast/*.go

repl:
	@go run main.go
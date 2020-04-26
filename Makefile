.PHONY: test repl bench

all: test

test:	
	@go test lexer/*.go
	@go test parser/*.go
	@go test ast/*.go
	@go test evaluator/*.go
	@go test object/*.go

repl:
	@go run main.go repl

bench: 
	@go run main.go bench
.PHONY: test repl bench

all: test

test:	
	@go test lexer/*.go
	@go test parser/*.go
	@go test ast/*.go
	@go test evaluator/*.go
	@go test object/*.go
	@go test code/*.go
	@go test compiler/*.go
	@go test vm/*.go

repl:
	@go run main.go repl

bench: 
	@go run main.go bench
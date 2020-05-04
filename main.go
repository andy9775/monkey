package main

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/andy9775/monkey/evaluator"
	"github.com/andy9775/monkey/lexer"
	"github.com/andy9775/monkey/object"
	"github.com/andy9775/monkey/parser"
	"github.com/andy9775/monkey/repl"
)

func main() {

	if os.Args[1] == "bench" {
		input := `
	let fib = fn(x) {
		if (x <= 1){
			return x;
		} 
		return fib(x - 1) + fib(x - 2);
	}
	puts(fib(30));
	`

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()

		num := 0

		start := time.Now()
		for num < 15 { // averages ~1.50 seconds vs python ~0.30 seconds
			evaluator.Eval(program, env)
			num++
		}
		end := time.Since(start)
		fmt.Printf("took: %s\n", end/time.Duration(num))
	} else if os.Args[1] == "repl" {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Printf("hello %s! This is the monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

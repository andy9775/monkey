package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/andy9775/interpreter/evaluator"
	"github.com/andy9775/interpreter/lexer"
	"github.com/andy9775/interpreter/object"
	"github.com/andy9775/interpreter/parser"
	"github.com/andy9775/interpreter/token"
)

// PROMPT is the text input console prompt
const PROMPT = ">> "

// Start begins the repl loop
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)        // print prompt and accept new input
		scanned := scanner.Scan() // read input

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line) // new lexer
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		// parse the tokens from the provided text
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

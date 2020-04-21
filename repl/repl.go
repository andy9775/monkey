package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/andy9775/interpreter/lexer"
	"github.com/andy9775/interpreter/token"
)

// PROMPT is the text input console prompt
const PROMPT = ">> "

// Start begins the repl loop
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)        // print prompt and accept new input
		scanned := scanner.Scan() // read input

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line) // new lexer

		// parse the tokens from the provided text
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/andy9775/monkey/compiler"
	"github.com/andy9775/monkey/lexer"
	"github.com/andy9775/monkey/object"
	"github.com/andy9775/monkey/parser"
	"github.com/andy9775/monkey/vm"
)

// PROMPT is the text input console prompt
const PROMPT = ">> "

// Start begins the repl loop
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

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

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed:\n%s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n%s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

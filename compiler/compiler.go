package compiler

import (
	"fmt"
	"sort"

	"github.com/andy9775/monkey/ast"
	"github.com/andy9775/monkey/code"
	"github.com/andy9775/monkey/object"
)

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants []object.Object

	// scopes are used to track instructions for the currrent call depth.
	// functions have their own scope
	scopes     []CompilationScope
	scopeIndex int

	symbolTable *SymbolTable
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// New returns a new instance of the compiler
func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	return &Compiler{
		constants: []object.Object{},

		scopes:     []CompilationScope{mainScope},
		scopeIndex: 0,

		symbolTable: NewSymbolTable(),
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

// Compile compiles the program and generates the bytecode
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

		c.emit(code.OpPop)
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.InfixExpression:
		if node.Operator == "<" { // swap order of operands
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}
		/*
			The flow of infix expressions:
			1. put the left element on the stack
			2. put the right element on the stack
			3. put the operator at the top of the stack
		*/
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// emit `OpJumpNotTruthy` with a bogus value for now
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		/*
			if the censequence emmits a OpPop we need to drop it else we
			risk getting rid of the evaluated result
		*/
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// emit `OpJumpNotTruthy` with a bogus value for now
		//
		// we have an alternative so the consequence should end with
		// an unconditional jump past the alternative
		jumpPos := c.emit(code.OpJump, 9999)

		// fix the truthy jump position
		afterConsequencePos := len(c.scopes[c.scopeIndex].instructions) // where to jump to
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			// Null is the default alternative
			c.emit(code.OpNull)
		} else {
			// compile the alternative branch
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		// change the unconditional jump position
		afterAlternativePos := len(c.scopes[c.scopeIndex].instructions)
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(node.Value) // compile the expression (lhs of assignment)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)

		if symbol.Scope == GlobalScope {
			// emit the correctly scoped instruction based on the symbol tables scope
			// this correctly binds locals
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok { // compile time error
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		// get a variable scopeed locally or globally
		if symbol.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, symbol.Index)
		} else {
			c.emit(code.OpGetLocal, symbol.Index)
		}

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}

		}
		// the last len(node.Elements) belong to the array
		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}

			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpression:
		err := c.Compile(node.Left) // e.g. left[1]
		if err != nil {
			return err
		}

		err = c.Compile(node.Index) // e.g. left[index]
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)

	case *ast.FunctionLiteral:
		// the functionliteral has it's own scope in which it compiles in
		c.enterScope()

		for _, p := range node.Parameters { // add the passed arguments to the local symbol table
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs((code.OpPop)) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			// last isn't a return value which means prevoius was not pop
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions // number of local bindings used by the function
		instructions := c.leaveScope()

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		c.emit(code.OpConstant, c.addConstant(compiledFn))

	case *ast.ReturnStatement:
		// the return value is an expression
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.CallExpression:
		err := c.Compile(node.Function) // identifier
		if err != nil {
			return err
		}

		for _, a := range node.Arguments { // handle the arguments
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments)) // tell the call how many arguments are needed
	}

	return nil
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.scopes[c.scopeIndex].instructions
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

// addConstant adds a constant value to the constants array and returns it's position
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit creates and adds code instructions to the instructions array
// and returns the starting position of the instruction
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.scopes[c.scopeIndex].instructions)
	c.scopes[c.scopeIndex].instructions = append(c.scopes[c.scopeIndex].instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.scopes[c.scopeIndex].instructions) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	c.scopes[c.scopeIndex].instructions = c.
		scopes[c.scopeIndex].
		instructions[:c.scopes[c.scopeIndex].lastInstruction.Position]
	c.scopes[c.scopeIndex].lastInstruction = c.scopes[c.scopeIndex].previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.scopes[c.scopeIndex].instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.scopes[c.scopeIndex].instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

// Bytecode contains the instructions generated by the compiler and the
// constants referenced throughout the program
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// Bytecode returns the bytecode for the application
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.scopes[c.scopeIndex].instructions,
		Constants:    c.constants,
	}
}

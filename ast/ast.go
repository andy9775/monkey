package ast

import (
	"bytes"
	"strings"

	"github.com/andy9775/interpreter/token"
)

// Node is a node in the abstract syntax tree
type Node interface {
	TokenLiteral() string // for debugging and testing

	String() string
}

// Statement represents a full statement
type Statement interface {
	Node
	statementNode()
}

// Expression is part of the statement which returns a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every ast
type Program struct {
	Statements []Statement // a full program
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns all of the statements representing the program (the original string of our program)
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ---------- let statement ----------

type LetStatement struct { // the full statement
	Token token.Token
	Name  *Identifier // variable name
	Value Expression  // what the value is
}

func (ls *LetStatement) statementNode()       {}                          // statement interface
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } // node interface

// String returns the string representation of the given let statement
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String()) // the identifier
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct { // left hand side (variable name)
	Token token.Token
	Value string // name of the identifier e.g. x
}

func (i *Identifier) expressionNode()      {}                         // expression interface
func (i *Identifier) TokenLiteral() string { return i.Token.Literal } // node interface

// String returns the string representation of the specific identifier
func (i *Identifier) String() string {
	return i.Value
}

// ----------- return statement -----------

// ReturnStatement represents the line `return <expression>`
type ReturnStatement struct {
	Token       token.Token // return token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}                          // statement interface
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal } // node interface

// String returns the string representation of the return statement
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ----------- expression statement -----------

// ExpressionStatement represents the line `x + 10`
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression  // the expression its self
}

func (es *ExpressionStatement) statementNode()       {}                          // node interface
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal } // node interface

// String reprents the string value of this expression
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// ----------- literals --------------

type IntegerLiteral struct {
	Token token.Token
	Value int64 // value of the expression e.g. 5
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// ---------- prefix ----------

type PrefixExpression struct {
	Token    token.Token // the prefix token e.g. !
	Operator string
	Right    Expression // right side of the token
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// ----------- infix --------------

type InfixExpression struct {
	Token    token.Token // the operator token e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// ---------- booleans ----------

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// ------------ if ------------

type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression  // brackets
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement // statements composing of this block
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ---------------- functions ----------------

type FunctionLiteral struct {
	Token      token.Token     // the fn token
	Parameters []*Identifier   // the list of function parameters (simple identifiers)
	Body       *BlockStatement // the body of a function is just a block
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // identifier or function literal (fn(a,b){...}(1,2))
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

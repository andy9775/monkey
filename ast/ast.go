package ast

import "github.com/andy9775/interpreter/token"

// Node is a node in the abstract syntax tree
type Node interface {
	TokenLiteral() string // for debugging and testing
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

// ---------- let statement ----------

type LetStatement struct { // the full statement
	Token token.Token
	Name  *Identifier // variable name
	Value Expression  // what the value is
}

func (ls *LetStatement) statementNode()       {}                          // statement interface
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } // node interface

type Identifier struct { // left hand side (variable name)
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}                         // expression interface
func (i *Identifier) TokenLiteral() string { return i.Token.Literal } // node interface

// ----------- return statement -----------

// ReturnStatement represents the line `return <expression>`
type ReturnStatement struct {
	Token       token.Token // return token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}                          // statement interface
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal } // node interface

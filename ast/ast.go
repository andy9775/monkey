package ast

import (
	"bytes"

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
	Value string // name of the identifier
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

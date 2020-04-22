package parser

import (
	"fmt"

	"github.com/andy9775/interpreter/ast"
	"github.com/andy9775/interpreter/lexer"
	"github.com/andy9775/interpreter/token"

	"strconv"
)

// identify the operator precedence from lowest to highest
// abs value doesn't matter but the relation to one another does
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// prefix and infix parsing functions for Pratt parser
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression /* left side of of operator being parsed */) ast.Expression
)

// Parser handles parsing the given text
type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixparseFns  map[token.TokenType]infixParseFn
}

// New returns a new instance of the parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// read two tokens so curr and peek are set
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns the array of error messages
func (p *Parser) Errors() []string {
	return p.errors
}

// ---------------- parse program ----------------

// ParseProgram parses the full program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // root node
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// -------------- parse statements --------------

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET: // current token is let
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// only have two statements, hence if we don't encounter either, it's an expression
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { // expression statements have optional semicolons
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) { // let identifier - identifier should be next token
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) { // let identifier = - next token should be an =
		return nil
	}

	for !p.currTokenIs(token.SEMICOLON) { // keep going till we hit a semi colon (skip expression)
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	for !p.currTokenIs(token.SEMICOLON) { // skip till end of line
		p.nextToken()
	}
	return stmt
}

// ----------- parse expressions -------------

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil { // not a prefix operator
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix() // parse the prefix operator

	return leftExp
}

// parseIdentifier returns an expression representing an identifier e.g `let x = 5;`
// 5 is the expression
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

// parseIntegerLiteral takes an integer and converts it into the IntegerLiteral AST node
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX) // set the right side

	return expression
}

// ----------------- helpers -----------------

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixparseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// check if the current token is of specific type
func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// check if the next token is of a specific type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// check next token and move to next if true
// used to assert that the next token is correct
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

package parser

import (
	"fmt"

	"github.com/andy9775/interpreter/ast"
	"github.com/andy9775/interpreter/lexer"
	"github.com/andy9775/interpreter/token"
)

// Parser handles parsing the given text
type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token
}

// New returns a new instance of the parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

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
		return nil
	}
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

// ----------------- helpers -----------------

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

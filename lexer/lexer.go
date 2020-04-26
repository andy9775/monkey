package lexer

import "github.com/andy9775/interpreter/token"

// Lexer iterates through the sourcecode and outputs tokens
type Lexer struct {
	input string

	// we use two positions in order to peek forward
	position     int // current position in input (points to current char (ch))
	readPosition int // current reading position in input (after current char)

	ch byte // current char under examination
}

// New creates a new instance of the lexer used to lex the input string
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // init the lexer
	return l
}

// NextToken reads the current character and returns the Token representing it
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' { // not equal
			ch := l.ch // get curr character and increment to next
			l.readChar()
			literal := string(ch) + string(l.ch) // set the two character literal
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else { // regular !
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			/*
				For any contigious series of letter which are parsed, we check to see if they
				match some keyword identifier (e.g. let, return, if, else, ...). If they do,
				its considered a keyword, otherwise it's a user defined identifier (x, y, a, b)
			*/
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}

		// don't know the type of token we have
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()

	return tok
}

// skipWhitespace keeps reading characters if we hit a whitespace since we don't care for it
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// read the specified identifier and return it
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { // keep going while we have a letter
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for { // read characters till we get to a closing quote
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

// readChar gets the next character and advances the pointer one step
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { // end of file
		l.ch = 0
	} else { // set the next character
		l.ch = l.input[l.readPosition]
	}
	// advance the pointers
	l.position = l.readPosition // the current position
	l.readPosition++            // where we are going next
}

// peek at the next character
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// readNumber keeps iterating through characters in order to get a full integer
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// ===================== helpers =====================

func isLetter(ch byte) bool {
	// outlines the letters (characters) that are accepted by our language as identifiers
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// helper function to return a new token of the specified type for the specified character
func newToken(TokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: TokenType, Literal: string(ch)}
}

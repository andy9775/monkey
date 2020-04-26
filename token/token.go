package token

// TokenType identifies the type of token encountered
// distinguish between (, ), 1, ;, etc - the different types
// using string isn't performent, but helps debugging -
// int or byte may be better for performance
type TokenType string

// Token represents the encountered token during the lexing stage
type Token struct {
	// specifies what was encountered, (,  ), [, ], let, 1, etc.
	Type TokenType
	// the literal value of the type e.g. 5
	Literal string
}

// ================ specify the different token types ================
const (
	ILLEGAL = "ILLEGAL" // token characeter we aren't familiar with
	EOF     = "EOF"     // we've reached the end of the file

	// identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y
	INT    = "INT"   // 1,2,3,4,5,....
	STRING = "STRING"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	LTE = "<="
	GTE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent matches the specified identifier to it's character representation
// returns either the keyword or a user specified identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT // all user defined identifiers
}

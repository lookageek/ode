package token

// fixed number of token types captured as constants
const (
	ILLEGAL     = "ILLEGAL"
	EOF         = "EOF"
	IDENT       = "IDENT"
	INT         = "INT"
	ASSIGN      = "="
	PLUS        = "+"
	COMMA       = ","
	SEMICOLON   = ";"
	LPAREN      = "("
	RPAREN      = ")"
	LBRACE      = "{"
	RBRACE      = "}"
	FUNCTION    = "FUNCTION"
	LET         = "LET"
	MINUS       = "-"
	DIVIDE      = "/"
	MULTIPLY    = "*"
	LESSTHAN    = "<"
	GREATERTHAN = ">"
	NEGATION    = "!"
	IF          = "IF"
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	ELSE        = "ELSE"
	RETURN      = "RETURN"
	EQ          = "=="
	NOTEQ       = "!="
	STRING      = "STRING"
	LBRACKET    = "["
	RBRACKET    = "]"
	COLON       = ":"
)

type TokenType string

// Token has two values - what type the token is, referred
// from the constant lookup set of all possible token types
// and the value of the token
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

// LookupIdent determines if the identifier is a keyword or a
// variable identifier
func LookupIdent(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}

	return IDENT
}

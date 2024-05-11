package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	//IDENTIFIERS AND LITERALS
	IDENT = "IDENT"
	INT = "INT"
	STRING = "STRING"
	
	//OPERATORS
	ASSIGN = "=" 
	PLUS = "+"
	MINUS = "-"
	MULTIPLY = "*"
	DIVIDE = "/"

	//DELIMITERS
	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRAK = "["
	RBRAK = "]"
)

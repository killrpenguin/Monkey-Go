package token

import (

)

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tkn, ok := keywords[ident]; ok {
		return tkn
	}
	return IDENT
}

const (
    ILLEGAL = "ILLEGAL"
    EOF     = "EOF"
    // Identifiers + literals
    IDENT = "IDENT"
    INT   = "INT"
	FLOAT = "FLOAT"
    // Operators
	ASSIGN   = "="
    PLUS     = "+"
    MINUS    = "-"
    BANG     = "!"
    ASTERISK = "*"
    SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="
	LTEQ     = "<="
	GTEQ     = ">="
    LT = "<"
    GT = ">"
    // Delimiters
    COMMA     = ","
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"
    // Keywords
    FUNCTION = "FUNCTION"
    LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)




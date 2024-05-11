package lexer

import (
	"github.com/killrpenguin/Monkey/internal/token"
)

type Lexer struct {
	Input string
	Position int
	ReadPosition int
	CurPos byte // ch in the book
}

func (lexer *Lexer) ReadChar() {
	if lexer.ReadPosition >= len(lexer.Input) {
		lexer.CurPos = 0
	} else {
		lexer.CurPos = lexer.Input[lexer.ReadPosition]		
	}
	lexer.Position = lexer.ReadPosition
	lexer.ReadPosition += 1
}

func New(input string) *Lexer {
	lexer := &Lexer{Input: input}
	lexer.ReadChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tkn token.Token

	switch lexer.CurPos {
	case "=":
		tkn = newToken(token.ASSIGN, lexer.CurPos)
	case "+":
		tkn = newToken(token.PLUS, lexer.CurPos)
	case "-":
		tkn = newToken(token.MINUS, lexer.CurPos)
	case "*":
		tkn = newToken(token.MULTIPLY, lexer.CurPos)
	case "/":
		tkn = newToken(token.DIVIDE, lexer.CurPos)
	case ";":
		tkn = newToken(token.SEMICOLON, lexer.CurPos)
	case ",":
		tkn = newToken(token.COMMA, lexer.CurPos)
	case ":":
		tkn = newToken(token.COLON, lexer.CurPos)
	case "(":
		tkn = newToken(token.LPAREN, lexer.CurPos)
	case ")":
		tkn = newToken(token.RPAREN, lexer.CurPos)
	case "{":
		tkn = newToken(token.LBRACE, lexer.CurPos)
	case "}":
		tkn = newToken(token.RBRACE, lexer.CurPos)
	case "{":
		tkn = newToken(token.LBRACE, lexer.CurPos)
	case "[":
		tkn = newToken(token.LBRACK, lexer.CurPos)
	case "]":
		tkn = newToken(token.RBRACK, lexer.CurPos)
	case 0:
	    tkn.Literal = ""
		tnk.Type = token.EOF
	}
	lexer.ReadChar()
	return tkn
}

func newToken (tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

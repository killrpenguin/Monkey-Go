package lexer

import (
	"Monkey/internal/token"
	"strings"
)

type Lexer struct {
	input string
	pos int
	readPos int
	ch byte	
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) readChar() {
	if lexer.readPos >= len(lexer.input) {
		lexer.ch = 0
	} else {
		lexer.ch = lexer.input[lexer.readPos]
	}
	lexer.pos = lexer.readPos
	lexer.readPos += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (lexer *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?'
}

func (lexer *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func (lexer *Lexer) readIdentifier() string {
	pos := lexer.pos
	for lexer.isLetter(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[pos:lexer.pos]
}

func (lexer *Lexer) readNumber() string {
	pos := lexer.pos
	for lexer.isDigit(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[pos:lexer.pos]
}

func (lexer *Lexer) eatWhitespace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
        lexer.readChar()
    }
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPos >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPos]
	}
}

func (lexer *Lexer) combineLiteral() string {
	ch := lexer.ch
	lexer.readChar()
	return string(ch) + string(lexer.ch)
}

func (lexer *Lexer) NextToken() token.Token {
	var tkn token.Token
	lexer.eatWhitespace()

	switch lexer.ch {
	case '=':
		if lexer.peekChar() == '=' {
			literal :=  lexer.combineLiteral()
			tkn = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tkn = newToken(token.ASSIGN, lexer.ch)
		}
    case '+':
        tkn = newToken(token.PLUS, lexer.ch)		
    case '-':
        tkn = newToken(token.MINUS, lexer.ch)
    case '!':
		if lexer.peekChar() == '=' {
			literal :=  lexer.combineLiteral()
			tkn = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tkn = newToken(token.BANG, lexer.ch)
		}
    case '*':
        tkn = newToken(token.ASTERISK, lexer.ch)		
    case '/':
        tkn = newToken(token.SLASH, lexer.ch)
    case '<':
		if lexer.peekChar() == '=' {
			literal :=  lexer.combineLiteral()
			tkn = token.Token{Type: token.LTEQ, Literal: literal}
		} else {
			tkn = newToken(token.LT, lexer.ch)
		}
    case '>':
		if lexer.peekChar() == '=' {
			literal :=  lexer.combineLiteral()
			tkn = token.Token{Type: token.GTEQ, Literal: literal}
		} else {
			tkn = newToken(token.GT, lexer.ch)
		}
    case ';':
        tkn = newToken(token.SEMICOLON, lexer.ch)
    case '(':
        tkn = newToken(token.LPAREN, lexer.ch)
    case ')':
        tkn = newToken(token.RPAREN, lexer.ch)
    case ',':
        tkn = newToken(token.COMMA, lexer.ch)
    case '{':
        tkn = newToken(token.LBRACE, lexer.ch)
    case '}':
        tkn = newToken(token.RBRACE, lexer.ch)
    case 0:
        tkn.Literal = ""
        tkn.Type = token.EOF
	default:
		if lexer.isLetter(lexer.ch) {
			tkn.Literal = lexer.readIdentifier()
			tkn.Type = token.LookupIdent(tkn.Literal)
			return tkn
		} else if lexer.isDigit(lexer.ch) {
			tkn.Literal = lexer.readNumber()
			if strings.Contains(tkn.Literal, ".") {
				tkn.Type = token.FLOAT
			} else {
				tkn.Type = token.INT
			}
			return tkn
		} else {
			tkn = newToken(token.ILLEGAL, lexer.ch)
		}
    }
	lexer.readChar()
	return tkn
}


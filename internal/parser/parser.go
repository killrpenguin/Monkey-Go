package parser

import (
	"Monkey/internal/token"
	"Monkey/internal/lexer"
	"Monkey/internal/ast"
	"fmt"
	"strings"
	"strconv"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ: EQUALS,
    token.NOT_EQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.SLASH: PRODUCT,
    token.ASTERISK: PRODUCT,
}

type Parser struct {
	lxr *lexer.Lexer
	errors []string

	curToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lxr: lexer,
		errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()
	
	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
    parser.registerInfix(token.PLUS, parser.parseInfixExpression)
    parser.registerInfix(token.MINUS, parser.parseInfixExpression)
    parser.registerInfix(token.SLASH, parser.parseInfixExpression)
    parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
    parser.registerInfix(token.EQ, parser.parseInfixExpression)
    parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
    parser.registerInfix(token.LT, parser.parseInfixExpression)
    parser.registerInfix(token.GT, parser.parseInfixExpression)
	
	return parser
}

func (parser *Parser) curTokenIs(tknType token.TokenType) bool {
	return parser.curToken.Type == tknType
}

func (parser *Parser) peekTokenIs(tknType token.TokenType) bool {
	return parser.peekToken.Type == tknType
}

func (parser *Parser) expectPeek(tknType token.TokenType) bool {
	if parser.peekTokenIs(tknType) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tknType)
		return false
	}
}

func (parser *Parser) peekPrecedence() int {
	if p, ok := precedences[parser.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (parser *Parser) curPrecedence() int {
	if p, ok := precedences[parser.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(tknType token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead",
		tknType, parser.peekToken.Type)
	parser.errors = append(parser.errors,msg)
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lxr.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.SingleStatement{}

	for parser.curToken.Type != token.EOF {
		parse_stmt := parser.parseStatement()
		if parse_stmt != nil {
			program.Statements = append(program.Statements, parse_stmt)
		}
		parser.nextToken()
	}
	return program
}

func (parser *Parser) parseStatement() ast.SingleStatement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.IDENT) {
		return nil		
	}
	stmt.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	
	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.curToken}

	parser.nextToken()

	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))
    stmt := &ast.ExpressionStatement{Token: parser.curToken}
    stmt.Expression = parser.parseExpression(LOWEST)

    if parser.peekTokenIs(token.SEMICOLON) {
        parser.nextToken()
    }
    return stmt
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (parser *Parser) noPrefixParseFnError(tkn token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found.", tkn)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	prefix := parser.prefixParseFns[parser.curToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.curToken.Type)
		return nil
	}
	
	leftExp := prefix()
	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		parser.nextToken()
		leftExp = infix(leftExp)
	}
	
	return leftExp
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token: parser.curToken,
		Operator: parser.curToken.Literal,
	}
	parser.nextToken()
	expression.Right = parser.parseExpression(PREFIX)
	return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token: parser.curToken,
		Operator: parser.curToken.Literal,
		Left: left,
	}

	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

// Tracing statement code from the book. Use the go test -v command below to invoke.
// go test -v -run TestOperatorPrecedenceParsing ./parser_test.go
func incIdent() {traceLevel = traceLevel + 1}
func decIdent() {traceLevel = traceLevel - 1}

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}
func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func trace(msg string) string {
	incIdent()
	tracePrint("Begin " + msg)
	return msg
}
func untrace(msg string) {
	tracePrint("End " + msg)
	decIdent()
}


package ast

import (
	"Monkey/internal/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type SingleStatement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}


type Program struct {
	Statements []SingleStatement
}

func (prog *Program) TokenLiteral() string {
	if len(prog.Statements) > 0 {
		return prog.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (prog *Program) String() string {
	var out bytes.Buffer
	for _, str := range prog.Statements {
		out.WriteString(str.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (letstmt *LetStatement) statementNode() {}
func (letstmt *LetStatement) TokenLiteral() string { return letstmt.Token.Literal }

func (letstmt *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(letstmt.TokenLiteral() + " ")
	out.WriteString(letstmt.Name.String())
	out.WriteString(" = ")

	if letstmt.Value != nil {
		out.WriteString(letstmt.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (retstmt *ReturnStatement) statementNode() {}
func (retstmt *ReturnStatement) TokenLiteral() string { return retstmt.Token.Literal }

func (retstmt *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(retstmt.TokenLiteral() + " ")

	if retstmt.ReturnValue != nil {
		out.WriteString(retstmt.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}
func (ident *Identifier) expressionNode() {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string { return ident.Value }

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}
func (ex *ExpressionStatement) statementNode() {}
func (ex *ExpressionStatement) TokenLiteral() string { return ex.Token.Literal }
func (ex *ExpressionStatement) String() string {
	if ex.Expression != nil {
		return ex.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return  il.Token.Literal }

type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	
	return out.String()
}

type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}


package ast

import "github.com/poolpOrg/julu/lexer"

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token lexer.Token // the token.IDENT token
	Value string
}

func NewIdentifier(token lexer.Token) *Identifier {
	return &Identifier{
		Token: token,
		Value: token.Literal,
	}
}

func (n *Identifier) expressionNode() {}
func (n *Identifier) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Identifier) String() string {
	return n.Value
}

type IntegerLiteral struct {
	Token lexer.Token // the token.INT token
	Value int64
}

func NewIntegerLiteral(token lexer.Token, value int64) *IntegerLiteral {
	return &IntegerLiteral{
		Token: token,
		Value: value,
	}
}
func (n *IntegerLiteral) expressionNode() {}
func (n *IntegerLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *IntegerLiteral) String() string {
	return n.Token.Literal
}

type PrefixExpression struct {
	Token    lexer.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func NewPrefixExpression(token lexer.Token, operator string) *PrefixExpression {
	return &PrefixExpression{
		Token:    token,
		Operator: operator,
	}
}
func (n *PrefixExpression) expressionNode() {}
func (n *PrefixExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *PrefixExpression) String() string {
	return "(" + n.Operator + n.Right.String() + ")"
}

type InfixExpression struct {
	Token    lexer.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func NewInfixExpression(token lexer.Token, left Expression) *InfixExpression {
	return &InfixExpression{
		Token:    token,
		Operator: token.Literal,
		Left:     left,
	}
}
func (n *InfixExpression) expressionNode() {}
func (n *InfixExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *InfixExpression) String() string {
	return "(" + n.Left.String() + " " + n.Operator + " " + n.Right.String() + ")"
}

type Boolean struct {
	Token lexer.Token
	Value bool
}

func NewBoolean(token lexer.Token, value bool) *Boolean {
	return &Boolean{
		Token: token,
		Value: value,
	}
}
func (n *Boolean) expressionNode() {}
func (n *Boolean) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Boolean) String() string {
	return n.Token.Literal
}

type IfExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func NewIfExpression(token lexer.Token) *IfExpression {
	return &IfExpression{
		Token: token,
	}
}
func (n *IfExpression) expressionNode() {}
func (n *IfExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *IfExpression) String() string {
	out := "if " + n.Condition.String() + " " + n.Consequence.String()
	if n.Alternative != nil {
		out += " else " + n.Alternative.String()
	}
	return out
}

type FunctionLiteral struct {
	Token      lexer.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func NewFunctionLiteral(token lexer.Token) *FunctionLiteral {
	return &FunctionLiteral{
		Token: token,
	}
}
func (n *FunctionLiteral) expressionNode() {}
func (n *FunctionLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *FunctionLiteral) String() string {
	params := []string{}
	for _, p := range n.Parameters {
		params = append(params, p.String())
	}
	return n.Token.Literal + "(" + "..." + ") " + n.Body.String()
}

type CallExpression struct {
	Token     lexer.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func NewCallExpression(token lexer.Token, function Expression) *CallExpression {
	return &CallExpression{
		Token:    token,
		Function: function,
	}
}
func (n *CallExpression) expressionNode() {}
func (n *CallExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *CallExpression) String() string {
	args := []string{}
	for _, a := range n.Arguments {
		args = append(args, a.String())
	}
	return n.Function.String() + "(" + n.Token.Literal + ")"
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func NewStringLiteral(token lexer.Token) *StringLiteral {
	return &StringLiteral{
		Token: token,
		Value: token.Literal,
	}
}
func (n *StringLiteral) expressionNode() {}
func (n *StringLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *StringLiteral) String() string {
	return "\"" + n.Token.Literal + "\""
}

type FStringLiteral struct {
	Token lexer.Token
	Value string
}

func NewFStringLiteral(token lexer.Token) *FStringLiteral {
	return &FStringLiteral{
		Token: token,
		Value: token.Literal,
	}
}
func (n *FStringLiteral) expressionNode() {}
func (n *FStringLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *FStringLiteral) String() string {
	return "f\"" + n.Token.Literal + "\""
}

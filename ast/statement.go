package ast

import (
	"fmt"
	"strings"

	"github.com/poolpOrg/julu/lexer"
)

type Statement interface {
	Node
	statementNode()
	Inspect(level int) string
}

type LetStatement struct {
	Token lexer.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func NewLetStatement(token lexer.Token) *LetStatement {
	return &LetStatement{
		Token: token,
	}
}
func (n *LetStatement) statementNode() {}
func (n *LetStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *LetStatement) String() string {
	return n.Token.Literal + " " + n.Name.String() + " = " + n.Value.String() + ";"
}
func (n *LetStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T: Name=%s\n", strings.Repeat(" ", level*2), n, n.Name.String())
	out += n.Value.Inspect(level + 1)
	return out
}

type ReturnStatement struct {
	Token       lexer.Token // the token.RETURN token
	ReturnValue Expression
}

func NewReturnStatement(token lexer.Token) *ReturnStatement {
	return &ReturnStatement{
		Token: token,
	}
}
func (n *ReturnStatement) statementNode() {}
func (n *ReturnStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *ReturnStatement) String() string {
	return n.Token.Literal + " " + n.ReturnValue.String() + ";"
}
func (n *ReturnStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	out += n.ReturnValue.Inspect(level + 1)
	return out
}

type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func NewExpressionStatement(token lexer.Token) *ExpressionStatement {
	return &ExpressionStatement{
		Token: token,
	}
}
func (n *ExpressionStatement) statementNode() {}
func (n *ExpressionStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *ExpressionStatement) String() string {
	if n.Expression != nil {
		return n.Expression.String()
	}
	return ""
}
func (n *ExpressionStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.TokenLiteral())
	if n.Expression != nil {
		out += n.Expression.Inspect(level + 1)
	}
	return out
}

type BlockStatement struct {
	Token      lexer.Token // the { token
	Statements []Statement
}

func NewBlockStatement(token lexer.Token) *BlockStatement {
	return &BlockStatement{
		Token: token,
	}
}
func (n *BlockStatement) statementNode() {}
func (n *BlockStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *BlockStatement) String() string {
	var out string

	if len(n.Statements) == 1 {
		out += "=> "
	} else {
		out += "{ "
	}
	for _, s := range n.Statements {
		out += s.String()
	}

	if len(n.Statements) != 1 {
		out += " }"
	}
	return out
}
func (n *BlockStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	for _, s := range n.Statements {
		out += s.Inspect(level + 1)
	}
	return out
}

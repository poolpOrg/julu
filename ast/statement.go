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
	Type  *Identifier
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

type MatchBlockStatement struct {
	Token lexer.Token // the { token
	Cases []CaseExpression
}

func NewMatchBlockStatement(token lexer.Token) *MatchBlockStatement {
	return &MatchBlockStatement{
		Token: token,
	}
}
func (n *MatchBlockStatement) statementNode() {}
func (n *MatchBlockStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *MatchBlockStatement) String() string {
	var out string

	if len(n.Cases) == 1 {
		out += "=> "
	} else {
		out += "{ "
	}
	for _, s := range n.Cases {
		out += s.String()
	}

	if len(n.Cases) != 1 {
		out += " }"
	}
	return out
}
func (n *MatchBlockStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	for _, s := range n.Cases {
		out += s.Inspect(level + 1)
	}
	return out
}

type BreakStatement struct {
	Token lexer.Token // the token.BREAK token
}

func NewBreakStatement(token lexer.Token) *BreakStatement {
	return &BreakStatement{
		Token: token,
	}
}
func (n *BreakStatement) statementNode() {}
func (n *BreakStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *BreakStatement) String() string {
	return n.Token.Literal
}
func (n *BreakStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	return out
}

type ContinueStatement struct {
	Token lexer.Token // the token.BREAK token
}

func NewContinueStatement(token lexer.Token) *ContinueStatement {
	return &ContinueStatement{
		Token: token,
	}
}
func (n *ContinueStatement) statementNode() {}
func (n *ContinueStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *ContinueStatement) String() string {
	return n.Token.Literal
}
func (n *ContinueStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	return out
}

type DoneStatement struct {
	Token lexer.Token // the token.BREAK token
}

func NewDoneStatement(token lexer.Token) *DoneStatement {
	return &DoneStatement{
		Token: token,
	}
}
func (n *DoneStatement) statementNode() {}
func (n *DoneStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *DoneStatement) String() string {
	return n.Token.Literal
}
func (n *DoneStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	return out
}

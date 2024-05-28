package ast

import "github.com/poolpOrg/julu/lexer"

type Statement interface {
	Node
	statementNode()
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

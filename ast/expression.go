package ast

import (
	"fmt"
	"strings"

	"github.com/poolpOrg/julu/lexer"
)

type Expression interface {
	Node
	expressionNode()
	Inspect(level int) string
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
func (n *Identifier) Inspect(level int) string {
	return fmt.Sprintf("%s%T: Name=%s\n", strings.Repeat(" ", level*2), n, n.String())
}

type IntegerLiteral struct {
	Token lexer.Token // the token.INT token
	Value int64
	Cast  *Identifier
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
func (n *IntegerLiteral) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type FloatLiteral struct {
	Token lexer.Token // the token.FLOAT token
	Value float64
	Cast  *Identifier
}

func NewFloatLiteral(token lexer.Token, value float64) *FloatLiteral {
	return &FloatLiteral{
		Token: token,
		Value: value,
	}
}
func (n *FloatLiteral) expressionNode() {}
func (n *FloatLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *FloatLiteral) String() string {
	return n.Token.Literal
}
func (n *FloatLiteral) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
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
func (n *PrefixExpression) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.TokenLiteral())
	out += n.Right.Inspect(level + 1)
	return out
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
func (n *InfixExpression) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.TokenLiteral())
	out += n.Left.Inspect(level + 1)
	out += n.Right.Inspect(level + 1)
	return out
}

type Boolean struct {
	Token lexer.Token
	Value bool
	Cast  *Identifier
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
func (n *Boolean) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type Null struct {
	Token lexer.Token
}

func NewNull(token lexer.Token) *Null {
	return &Null{
		Token: token,
	}
}
func (n *Null) expressionNode() {}
func (n *Null) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Null) String() string {
	return "null"
}
func (n *Null) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type IfExpression struct {
	Token                  lexer.Token // The 'if' token
	Condition              Expression
	Consequence            *BlockStatement
	ConditionalAlternative *IfExpression
	Alternative            *BlockStatement
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
func (n *IfExpression) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	out += strings.Repeat(" ", (level+1)*2) + "Condition:\n"
	out += n.Condition.Inspect(level + 2)
	if n.Consequence != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Consequence:\n"
		out += n.Consequence.Inspect(level + 2)
	}
	if n.ConditionalAlternative != nil {
		out += strings.Repeat(" ", (level+1)*2) + "ConditionalAlternative:\n"
		out += n.ConditionalAlternative.Inspect(level + 2)
	}
	if n.Alternative != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Alternative:\n"
		out += n.Alternative.Inspect(level + 2)
	}
	return out
}

type FunctionLiteral struct {
	Token      lexer.Token // The 'fn' token
	Name       *Identifier
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
func (n *FunctionLiteral) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T: Parameters=%s\n", strings.Repeat(" ", level*2), n, n.Parameters)
	out += n.Body.Inspect(level + 1)
	return out
}

type CallExpression struct {
	Token      lexer.Token // The '(' token
	Function   Expression  // Identifier or FunctionLiteral
	Parameters []Expression
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
	for _, a := range n.Parameters {
		args = append(args, a.String())
	}
	return n.Function.String() + "(" + n.Token.Literal + ")"
}
func (n *CallExpression) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T: Function=%s, Parameters=%s\n", strings.Repeat(" ", level*2), n, n.Function.String(), n.Parameters)
	return out
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
func (n *StringLiteral) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
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
func (n *FStringLiteral) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type ArrayLiteral struct {
	Token    lexer.Token // The '[' token
	Elements []Expression
}

func NewArrayLiteral(token lexer.Token) *ArrayLiteral {
	return &ArrayLiteral{
		Token: token,
	}
}
func (n *ArrayLiteral) expressionNode() {}
func (n *ArrayLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *ArrayLiteral) String() string {
	elements := []string{}
	for _, el := range n.Elements {
		elements = append(elements, el.String())
	}
	return "[" + "..." + "]"
}
func (n *ArrayLiteral) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type IndexExpression struct {
	Token lexer.Token // The '[' token
	Left  Expression
	Index Expression
}

func NewIndexExpression(token lexer.Token, left Expression) *IndexExpression {
	return &IndexExpression{
		Token: token,
		Left:  left,
	}
}
func (n *IndexExpression) expressionNode() {}
func (n *IndexExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *IndexExpression) String() string {
	return "(" + n.Left.String() + "[" + n.Index.String() + "])"
}
func (n *IndexExpression) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type HashLiteral struct {
	Token lexer.Token // The '{' token
	Pairs map[Expression]Expression
}

func NewHashLiteral(token lexer.Token) *HashLiteral {
	return &HashLiteral{
		Token: token,
	}
}
func (n *HashLiteral) expressionNode() {}
func (n *HashLiteral) TokenLiteral() string {
	return n.Token.Literal
}
func (n *HashLiteral) String() string {
	//pairs := []string{}
	//for key, value := range n.Pairs {
	//	pairs = append(pairs, key.String()+":"+value.String())
	//}
	return "{" + "..." + "}"
}
func (n *HashLiteral) Inspect(level int) string {
	return fmt.Sprintf("%s%T: %s\n", strings.Repeat(" ", level*2), n, n.String())
}

type LoopStatement struct {
	Token          lexer.Token // the { token
	WhileCondition Expression
	UntilCondition Expression
	Body           *BlockStatement
}

func NewLoopStatement(token lexer.Token) *LoopStatement {
	return &LoopStatement{
		Token: token,
	}
}

func (n *LoopStatement) expressionNode() {}
func (n *LoopStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *LoopStatement) String() string {
	return n.Token.Literal + " " + n.Body.String()
}
func (n *LoopStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	if n.WhileCondition != nil {
		out += strings.Repeat(" ", (level+1)*2) + "WhileCondition:\n"
		out += n.WhileCondition.Inspect(level + 2)
	}
	if n.UntilCondition != nil {
		out += strings.Repeat(" ", (level+1)*2) + "UntilCondition:\n"
		out += n.UntilCondition.Inspect(level + 2)
	}
	if n.WhileCondition != nil || n.UntilCondition != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Body:\n"
		out += n.Body.Inspect(level + 2)
	} else {
		out += n.Body.Inspect(level + 1)
	}
	return out
}

type ForStatement struct {
	Token    lexer.Token // the { token
	Variable Expression
	Iterable Expression
	Body     *BlockStatement
}

func NewForStatement(token lexer.Token) *ForStatement {
	return &ForStatement{
		Token: token,
	}
}

func (n *ForStatement) expressionNode() {}
func (n *ForStatement) TokenLiteral() string {
	return n.Token.Literal
}
func (n *ForStatement) String() string {
	return n.Token.Literal + " " + n.Body.String()
}
func (n *ForStatement) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)

	out += strings.Repeat(" ", (level+1)*2) + "Variable:\n"
	out += n.Variable.Inspect(level + 2)
	if n.Iterable != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Iterable:\n"
		out += n.Iterable.Inspect(level + 2)
	}
	return out
}

type MatchExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	MatchBlock  *MatchBlockStatement
	Alternative *BlockStatement
}

func NewMatchExpression(token lexer.Token) *MatchExpression {
	return &MatchExpression{
		Token: token,
	}
}
func (n *MatchExpression) expressionNode() {}
func (n *MatchExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *MatchExpression) String() string {
	out := "match " + n.Condition.String() + " "
	if n.Alternative != nil {
		out += " else " + n.Alternative.String()
	}
	return out
}
func (n *MatchExpression) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	out += strings.Repeat(" ", (level+1)*2) + "Condition:\n"
	out += n.Condition.Inspect(level + 2)

	if n.MatchBlock != nil {
		out += strings.Repeat(" ", (level+1)*2) + "MatchBlock:\n"
		out += n.MatchBlock.Inspect(level + 2)
	}

	if n.Alternative != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Alternative:\n"
		out += n.Alternative.Inspect(level + 2)
	}
	return out
}

type CaseExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	Guard       Expression
	Consequence *BlockStatement
}

func NewCaseExpression(token lexer.Token) *CaseExpression {
	return &CaseExpression{
		Token: token,
	}
}
func (n *CaseExpression) expressionNode() {}
func (n *CaseExpression) TokenLiteral() string {
	return n.Token.Literal
}
func (n *CaseExpression) String() string {
	out := "case " + n.Condition.String() + " " + n.Consequence.String()
	return out
}
func (n *CaseExpression) Inspect(level int) string {
	var out string
	out += fmt.Sprintf("%s%T\n", strings.Repeat(" ", level*2), n)
	out += strings.Repeat(" ", (level+1)*2) + "Condition:\n"
	out += n.Condition.Inspect(level + 2)
	if n.Guard != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Guard:\n"
		out += n.Guard.Inspect(level + 2)
	}
	if n.Consequence != nil {
		out += strings.Repeat(" ", (level+1)*2) + "Consequence:\n"
		out += n.Consequence.Inspect(level + 2)
	}
	return out
}

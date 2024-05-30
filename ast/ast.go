package ast

import "fmt"

type Node interface {
	TokenLiteral() string
	String() string
}

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: []Statement{},
	}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) Inspect() string {
	var out string
	out += fmt.Sprintf("%T:\n", p)
	for _, s := range p.Statements {
		out += s.Inspect(1)
	}
	return out
}

func (p *Program) String() string {
	var out string

	for _, s := range p.Statements {
		out += s.String()
	}

	return out
}

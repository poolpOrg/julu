package parser

import (
	"github.com/poolpOrg/julu/lexer"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	COMPARISON  // &&, ||
	SUM         // +
	PRODUCT     // *
	BITWISE     // &, |, ^
	BITSHIFT    // <<, >>
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[lexer.TokenType]int{
	lexer.EQUALS:           EQUALS,
	lexer.NOT_EQUALS:       EQUALS,
	lexer.LESSER_THAN:      LESSGREATER,
	lexer.GREATER_THAN:     LESSGREATER,
	lexer.LESSER_OR_EQUAL:  LESSGREATER,
	lexer.GREATER_OR_EQUAL: LESSGREATER,

	lexer.LOGICAL_AND: COMPARISON,
	lexer.LOGICAL_OR:  COMPARISON,

	lexer.BITWISE_AND:      BITWISE,
	lexer.BITWISE_OR:       BITWISE,
	lexer.BITWISE_XOR:      BITWISE,
	lexer.RSHIFT:           BITSHIFT,
	lexer.LSHIFT:           BITSHIFT,
	lexer.CIRCULAR_RSHIFT:  BITSHIFT,
	lexer.CIRCULAR_LSHIFT:  BITSHIFT,
	lexer.ADD:              SUM,
	lexer.SUB:              SUM,
	lexer.MUL:              PRODUCT,
	lexer.DIV:              PRODUCT,
	lexer.MOD:              PRODUCT,
	lexer.LEFT_PARENTHESIS: CALL,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

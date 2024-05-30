package parser_test

import (
	"testing"

	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/parser"
)

type MockParser struct {
	curToken  lexer.Token
	peekToken lexer.Token
}

func (m *MockParser) peekPrecedence() int {
	precedences := parser.GetPrecedenceTable()
	if p, ok := precedences[m.peekToken.Type]; ok {
		return p
	}
	return parser.LOWEST
}

func (m *MockParser) curPrecedence() int {
	precedences := parser.GetPrecedenceTable()
	if p, ok := precedences[m.curToken.Type]; ok {
		return p
	}
	return parser.LOWEST
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		tokenType lexer.TokenType
		expected  int
	}{
		{tokenType: lexer.EQUALS, expected: parser.EQUALS},
		{tokenType: lexer.NOT_EQUALS, expected: parser.EQUALS},
		{tokenType: lexer.LESSER_THAN, expected: parser.LESSGREATER},
		{tokenType: lexer.GREATER_THAN, expected: parser.LESSGREATER},
		{tokenType: lexer.LESSER_OR_EQUAL, expected: parser.LESSGREATER},
		{tokenType: lexer.GREATER_OR_EQUAL, expected: parser.LESSGREATER},
		{tokenType: lexer.LOGICAL_AND, expected: parser.COMPARISON},
		{tokenType: lexer.LOGICAL_OR, expected: parser.COMPARISON},
		{tokenType: lexer.BITWISE_AND, expected: parser.BITWISE},
		{tokenType: lexer.BITWISE_OR, expected: parser.BITWISE},
		{tokenType: lexer.BITWISE_XOR, expected: parser.BITWISE},
		{tokenType: lexer.RSHIFT, expected: parser.BITSHIFT},
		{tokenType: lexer.LSHIFT, expected: parser.BITSHIFT},
		{tokenType: lexer.CIRCULAR_RSHIFT, expected: parser.BITSHIFT},
		{tokenType: lexer.CIRCULAR_LSHIFT, expected: parser.BITSHIFT},
		{tokenType: lexer.ADD, expected: parser.SUM},
		{tokenType: lexer.SUB, expected: parser.SUM},
		{tokenType: lexer.MUL, expected: parser.PRODUCT},
		{tokenType: lexer.DIV, expected: parser.PRODUCT},
		{tokenType: lexer.MOD, expected: parser.PRODUCT},
		{tokenType: lexer.LEFT_PARENTHESIS, expected: parser.CALL},
		{tokenType: lexer.IDENTIFIER, expected: parser.LOWEST}, // Example of token not in the map
	}

	for _, tt := range tests {
		mockParser := &MockParser{
			curToken:  lexer.Token{Type: tt.tokenType},
			peekToken: lexer.Token{Type: tt.tokenType},
		}

		curPrec := mockParser.curPrecedence()
		if curPrec != tt.expected {
			t.Fatalf("expected current precedence %d, got %d for token type %s", tt.expected, curPrec, tt.tokenType)
		}

		peekPrec := mockParser.peekPrecedence()
		if peekPrec != tt.expected {
			t.Fatalf("expected peek precedence %d, got %d for token type %s", tt.expected, peekPrec, tt.tokenType)
		}
	}
}

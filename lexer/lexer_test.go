package lexer_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/poolpOrg/julu/lexer"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []lexer.Token
	}{
		{
			input: `+ - * / % ++ -- += -= *= /= %= 
			< <= << <<= <<< > >= >> >>= >>> == != = => &
			&= && | |= || ^ ^= ~ ( ) { } [ ] ; : , . 
			// comment
			/*
			multiline
			comment
			*/
			let x = 5; if (x > 1) { return x + 1; }
			`,
			expected: []lexer.Token{
				{Type: lexer.ADD, Literal: "+"},
				{Type: lexer.SUB, Literal: "-"},
				{Type: lexer.MUL, Literal: "*"},
				{Type: lexer.DIV, Literal: "/"},
				{Type: lexer.MOD, Literal: "%"},
				{Type: lexer.INCR, Literal: "++"},
				{Type: lexer.DECR, Literal: "--"},
				{Type: lexer.ADD_AND_ASSIGN, Literal: "+="},
				{Type: lexer.SUB_AND_ASSIGN, Literal: "-="},
				{Type: lexer.MUL_AND_ASSIGN, Literal: "*="},
				{Type: lexer.DIV_AND_ASSIGN, Literal: "/="},
				{Type: lexer.MOD_AND_ASSIGN, Literal: "%="},
				{Type: lexer.LESSER_THAN, Literal: "<"},
				{Type: lexer.LESSER_OR_EQUAL, Literal: "<="},
				{Type: lexer.LSHIFT, Literal: "<<"},
				{Type: lexer.LSHIFT_ASSIGN, Literal: "<<="},
				{Type: lexer.CIRCULAR_LSHIFT, Literal: "<<<"},
				{Type: lexer.GREATER_THAN, Literal: ">"},
				{Type: lexer.GREATER_OR_EQUAL, Literal: ">="},
				{Type: lexer.RSHIFT, Literal: ">>"},
				{Type: lexer.RSHIFT_ASSIGN, Literal: ">>="},
				{Type: lexer.CIRCULAR_RSHIFT, Literal: ">>>"},
				{Type: lexer.EQUALS, Literal: "=="},
				{Type: lexer.NOT_EQUALS, Literal: "!="},
				{Type: lexer.ASSIGN, Literal: "="},
				{Type: lexer.ARROW, Literal: "=>"},
				{Type: lexer.BITWISE_AND, Literal: "&"},
				{Type: lexer.BITWISE_AND_ASSIGN, Literal: "&="},
				{Type: lexer.LOGICAL_AND, Literal: "&&"},
				{Type: lexer.BITWISE_OR, Literal: "|"},
				{Type: lexer.BITWISE_OR_ASSIGN, Literal: "|="},
				{Type: lexer.LOGICAL_OR, Literal: "||"},
				{Type: lexer.BITWISE_XOR, Literal: "^"},
				{Type: lexer.BITWISE_XOR_ASSIGN, Literal: "^="},
				{Type: lexer.BITWISE_NOT, Literal: "~"},
				{Type: lexer.LEFT_PARENTHESIS, Literal: "("},
				{Type: lexer.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: lexer.LEFT_CURLY_BRACKET, Literal: "{"},
				{Type: lexer.RIGHT_CURLY_BRACKET, Literal: "}"},
				{Type: lexer.LEFT_SQUARE_BRACKET, Literal: "["},
				{Type: lexer.RIGHT_SQUARE_BRACKET, Literal: "]"},
				{Type: lexer.SEMICOLON, Literal: ";"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.DOT, Literal: "."},
				{Type: lexer.LET, Literal: "let"},
				{Type: lexer.IDENTIFIER, Literal: "x"},
				{Type: lexer.ASSIGN, Literal: "="},
				{Type: lexer.INTEGER, Literal: "5"},
				{Type: lexer.SEMICOLON, Literal: ";"},
				{Type: lexer.IF, Literal: "if"},
				{Type: lexer.LEFT_PARENTHESIS, Literal: "("},
				{Type: lexer.IDENTIFIER, Literal: "x"},
				{Type: lexer.GREATER_THAN, Literal: ">"},
				{Type: lexer.INTEGER, Literal: "1"},
				{Type: lexer.RIGHT_PARENTHESIS, Literal: ")"},
				{Type: lexer.LEFT_CURLY_BRACKET, Literal: "{"},
				{Type: lexer.RETURN, Literal: "return"},
				{Type: lexer.IDENTIFIER, Literal: "x"},
				{Type: lexer.ADD, Literal: "+"},
				{Type: lexer.INTEGER, Literal: "1"},
				{Type: lexer.SEMICOLON, Literal: ";"},
				{Type: lexer.RIGHT_CURLY_BRACKET, Literal: "}"},
				{Type: lexer.EOF, Literal: ""},
			},
		},
		{
			input: `'a' "string" 123 123.456 0b101 0o77 0xFA f"formatted string"`,
			expected: []lexer.Token{
				{Type: lexer.RUNE, Literal: "a"},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.INTEGER, Literal: "123"},
				{Type: lexer.FLOAT, Literal: "123.456"},
				{Type: lexer.INTEGER, Literal: "0b101"},
				{Type: lexer.INTEGER, Literal: "0o77"},
				{Type: lexer.INTEGER, Literal: "0xFA"},
				{Type: lexer.FSTRING, Literal: "formatted string"},
				{Type: lexer.EOF, Literal: ""},
			},
		},
		{
			input: `// Single line comment
			/*
			Multiline
			Comment
			*/
			"escaped\nstring" 'c' 1.23 invalid_token`,
			expected: []lexer.Token{
				{Type: lexer.STRING, Literal: "escaped\nstring"},
				{Type: lexer.RUNE, Literal: "c"},
				{Type: lexer.FLOAT, Literal: "1.23"},
				{Type: lexer.IDENTIFIER, Literal: "invalid_token"},
				{Type: lexer.EOF, Literal: ""},
			},
		},
		{
			input: `let y = 0o077; let z = 0x0F0;`,
			expected: []lexer.Token{
				{Type: lexer.LET, Literal: "let"},
				{Type: lexer.IDENTIFIER, Literal: "y"},
				{Type: lexer.ASSIGN, Literal: "="},
				{Type: lexer.INTEGER, Literal: "0o077"},
				{Type: lexer.SEMICOLON, Literal: ";"},
				{Type: lexer.LET, Literal: "let"},
				{Type: lexer.IDENTIFIER, Literal: "z"},
				{Type: lexer.ASSIGN, Literal: "="},
				{Type: lexer.INTEGER, Literal: "0x0F0"},
				{Type: lexer.SEMICOLON, Literal: ";"},
				{Type: lexer.EOF, Literal: ""},
			},
		},
	}

	for tid, tt := range tests {
		l := lexer.New(bufio.NewReader(strings.NewReader(tt.input)))
		for _, expectedToken := range tt.expected {
			tok := l.Lex()
			if tok.Type != expectedToken.Type {
				t.Fatalf("test %d: expected token type %q, got %q", tid, expectedToken.Type, tok.Type)
			}
			if tok.Literal != expectedToken.Literal {
				t.Fatalf("test %d: expected token literal %q, got %q", tid, expectedToken.Literal, tok.Literal)
			}
		}
	}
}

func TestLexNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected lexer.Token
	}{
		{input: "123", expected: lexer.Token{Type: lexer.INTEGER, Literal: "123"}},
		{input: "123.456", expected: lexer.Token{Type: lexer.FLOAT, Literal: "123.456"}},
		{input: "0b101", expected: lexer.Token{Type: lexer.INTEGER, Literal: "0b101"}},
		{input: "0o77", expected: lexer.Token{Type: lexer.INTEGER, Literal: "0o77"}},
		{input: "0xFA", expected: lexer.Token{Type: lexer.INTEGER, Literal: "0xFA"}},
	}

	for _, tt := range tests {
		l := lexer.New(bufio.NewReader(strings.NewReader(tt.input)))
		tok := l.Lex()
		if tok.Type != tt.expected.Type {
			t.Fatalf("expected token type %q, got %q", tt.expected.Type, tok.Type)
		}
		if tok.Literal != tt.expected.Literal {
			t.Fatalf("expected token literal %q, got %q", tt.expected.Literal, tok.Literal)
		}
	}
}

func TestLexIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected lexer.Token
	}{
		{input: "let", expected: lexer.Token{Type: lexer.LET, Literal: "let"}},
		{input: "x", expected: lexer.Token{Type: lexer.IDENTIFIER, Literal: "x"}},
		{input: "if", expected: lexer.Token{Type: lexer.IF, Literal: "if"}},
		{input: "y", expected: lexer.Token{Type: lexer.IDENTIFIER, Literal: "y"}},
	}

	for _, tt := range tests {
		l := lexer.New(bufio.NewReader(strings.NewReader(tt.input)))
		tok := l.Lex()
		if tok.Type != tt.expected.Type {
			t.Fatalf("expected token type %q, got %q", tt.expected.Type, tok.Type)
		}
		if tok.Literal != tt.expected.Literal {
			t.Fatalf("expected token literal %q, got %q", tt.expected.Literal, tok.Literal)
		}
	}
}

func TestLexString(t *testing.T) {
	tests := []struct {
		input    string
		expected lexer.Token
	}{
		{input: `"string"`, expected: lexer.Token{Type: lexer.STRING, Literal: "string"}},
		{input: "`raw string`", expected: lexer.Token{Type: lexer.STRING, Literal: "raw string"}},
		{input: `f"formatted string"`, expected: lexer.Token{Type: lexer.FSTRING, Literal: "formatted string"}},
	}

	for _, tt := range tests {
		l := lexer.New(bufio.NewReader(strings.NewReader(tt.input)))
		tok := l.Lex()
		if tok.Type != tt.expected.Type {
			t.Fatalf("expected token type %q, got %q", tt.expected.Type, tok.Type)
		}
		if tok.Literal != tt.expected.Literal {
			t.Fatalf("expected token literal %q, got %q", tt.expected.Literal, tok.Literal)
		}
	}
}

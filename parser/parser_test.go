package parser_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/poolpOrg/julu/ast"
	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/parser"
)

func newParser(input string) *parser.Parser {
	l := lexer.New(bufio.NewReader(strings.NewReader(input)))
	return parser.New(l)
}

func TestParseLetStatement(t *testing.T) {
	input := `let x = 5;`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.LetStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Fatalf("stmt.Name.Value not 'x'. got=%s", stmt.Name.Value)
	}

	if stmt.Value.String() != "5" {
		t.Fatalf("stmt.Value not '5'. got=%s", stmt.Value.String())
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := `return 5;`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ReturnStatement. got=%T", program.Statements[0])
	}

	if stmt.ReturnValue.String() != "5" {
		t.Fatalf("stmt.ReturnValue not '5'. got=%s", stmt.ReturnValue.String())
	}
}

func TestParseExpressionStatement(t *testing.T) {
	input := `5;`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if stmt.Expression.String() != "5" {
		t.Fatalf("stmt.Expression not '5'. got=%s", stmt.Expression.String())
	}
}

func TestParsePrefixExpression(t *testing.T) {
	input := `-5;`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T", stmt.Expression)
	}

	if exp.Operator != "-" {
		t.Fatalf("exp.Operator is not '-'. got=%s", exp.Operator)
	}

	if exp.Right.String() != "5" {
		t.Fatalf("exp.Right not '5'. got=%s", exp.Right.String())
	}
}

func TestParseInfixExpression(t *testing.T) {
	input := `5 + 5;`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.InfixExpression. got=%T", stmt.Expression)
	}

	if exp.Operator != "+" {
		t.Fatalf("exp.Operator is not '+'. got=%s", exp.Operator)
	}

	if exp.Left.String() != "5" {
		t.Fatalf("exp.Left not '5'. got=%s", exp.Left.String())
	}

	if exp.Right.String() != "5" {
		t.Fatalf("exp.Right not '5'. got=%s", exp.Right.String())
	}
}

func TestParseGroupedExpression(t *testing.T) {
	input := `(5 + 5);`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if stmt.Expression.String() != "(5 + 5)" {
		t.Fatalf("stmt.Expression not '(5 + 5)'. got=%s", stmt.Expression.String())
	}
}

func TestParseIfExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if exp.Condition.String() != "(x < y)" {
		t.Fatalf("exp.Condition not '(x < y)'. got=%s", exp.Condition.String())
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence does not contain 1 statement. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if consequence.Expression.String() != "x" {
		t.Fatalf("consequence.Expression not 'x'. got=%s", consequence.Expression.String())
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("exp.Alternative does not contain 1 statement. got=%d", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if alternative.Expression.String() != "y" {
		t.Fatalf("alternative.Expression not 'y'. got=%s", alternative.Expression.String())
	}
}

func TestParseFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters does not contain 2 parameters. got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("fn.Parameters[0] not 'x'. got=%s", fn.Parameters[0].String())
	}

	if fn.Parameters[1].String() != "y" {
		t.Fatalf("fn.Parameters[1] not 'y'. got=%s", fn.Parameters[1].String())
	}

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("fn.Body.Statements does not contain 1 statement. got=%d", len(fn.Body.Statements))
	}

	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fn.Body.Statements[0] is not *ast.ExpressionStatement. got=%T", fn.Body.Statements[0])
	}

	if bodyStmt.Expression.String() != "(x + y)" {
		t.Fatalf("bodyStmt.Expression not '(x + y)'. got=%s", bodyStmt.Expression.String())
	}
}

func TestParseCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	p := newParser(input)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if exp.Function.String() != "add" {
		t.Fatalf("exp.Function not 'add'. got=%s", exp.Function.String())
	}

	if len(exp.Parameters) != 3 {
		t.Fatalf("exp.Arguments does not contain 3 arguments. got=%d", len(exp.Parameters))
	}

	if exp.Parameters[0].String() != "1" {
		t.Fatalf("exp.Arguments[0] not '1'. got=%s", exp.Parameters[0].String())
	}

	if exp.Parameters[1].String() != "(2 * 3)" {
		t.Fatalf("exp.Arguments[1] not '(2 * 3)'. got=%s", exp.Parameters[1].String())
	}

	if exp.Parameters[2].String() != "(4 + 5)" {
		t.Fatalf("exp.Arguments[2] not '(4 + 5)'. got=%s", exp.Parameters[2].String())
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

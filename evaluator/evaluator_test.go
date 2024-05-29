package evaluator_test

import (
	"testing"

	"github.com/poolpOrg/julu/ast"
	"github.com/poolpOrg/julu/evaluator"
	"github.com/poolpOrg/julu/object"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected int64
	}{
		{input: &ast.IntegerLiteral{Value: 5}, expected: 5},
		{input: &ast.IntegerLiteral{Value: 10}, expected: 10},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(tt.input, env)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected bool
	}{
		{input: &ast.Boolean{Value: true}, expected: true},
		{input: &ast.Boolean{Value: false}, expected: false},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(tt.input, env)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected interface{}
	}{
		{
			input: &ast.PrefixExpression{
				Operator: "!",
				Right:    &ast.Boolean{Value: true},
			},
			expected: false,
		},
		{
			input: &ast.PrefixExpression{
				Operator: "!",
				Right:    &ast.Boolean{Value: false},
			},
			expected: true,
		},
		{
			input: &ast.PrefixExpression{
				Operator: "-",
				Right:    &ast.IntegerLiteral{Value: 5},
			},
			expected: -5,
		},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(tt.input, env)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected interface{}
	}{
		{
			input: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 5},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Value: 5},
			},
			expected: 10,
		},
		{
			input: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 5},
				Operator: "-",
				Right:    &ast.IntegerLiteral{Value: 5},
			},
			expected: 0,
		},
		{
			input: &ast.InfixExpression{
				Left:     &ast.Boolean{Value: true},
				Operator: "==",
				Right:    &ast.Boolean{Value: true},
			},
			expected: true,
		},
		{
			input: &ast.InfixExpression{
				Left:     &ast.Boolean{Value: true},
				Operator: "!=",
				Right:    &ast.Boolean{Value: false},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(tt.input, env)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}
}

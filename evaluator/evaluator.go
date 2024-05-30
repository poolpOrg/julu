package evaluator

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/poolpOrg/julu/ast"
	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/object"
	"github.com/poolpOrg/julu/parser"
)

var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
	VOID     = &object.Void{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.MatchExpression:
		return evalMatchExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.DoneStatement:
		return &object.ReturnValue{Value: VOID}

	case *ast.FunctionLiteral:
		if node.Name != nil {
			funcObj := &object.Function{Name: node.Name, Parameters: node.Parameters, Body: node.Body, Env: env}
			env.Set(node.Name.Value, funcObj)
			return funcObj
		} else {
			return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Parameters, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)

	case *ast.LoopStatement:
		return evalLoopStatement(node, env)

	case *ast.BreakStatement:
		return BREAK
	case *ast.ContinueStatement:
		return CONTINUE

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.Null:
		return NULL

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.FStringLiteral:
		return evalFStringLiteral(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	default:
		return newError("unknown node type: %T", node)
	}
	return nil
}

func EvalFunctionObject(fn object.Object, env *object.Environment) object.Object {
	if fn.Type() != object.FUNCTION_OBJ {
		return newError("not a function: %s", fn.Type())
	}

	fn = fn.(*object.Function)
	// evaluate function parameters and match with prototype
	//args := evalExpressions(fn.Parameters, env)
	//if len(args) == 1 && isError(args[0]) {
	//	return args[0]
	//}
	return applyFunction(fn, []object.Object{})
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE: //!true = false
		return FALSE
	case FALSE: //!false = true
		return TRUE
	case NULL: //!null = true
		return TRUE
	default: //!5 = false
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}

	case "<<":
		return &object.Integer{Value: leftVal << rightVal}
	case ">>":
		return &object.Integer{Value: leftVal >> rightVal}

	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)

		// for now, "is" is just a synonym for "==", we'll deal later with identity
	case "==", "is":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	case "&&", "and":
		return nativeBoolToBooleanObject(leftVal != 0 && rightVal != 0)

	case "||", "or":
		return nativeBoolToBooleanObject(leftVal != 0 || rightVal != 0)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}

	case "-":
		return &object.String{Value: strings.ReplaceAll(leftVal, rightVal, "")}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==", "is":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "&&", "and":
		return nativeBoolToBooleanObject(leftVal != "" && rightVal != "")

	case "||", "or":
		return nativeBoolToBooleanObject(leftVal != "" || rightVal != "")

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.ConditionalAlternative != nil {
		return evalIfExpression(ie.ConditionalAlternative, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalMatchExpression(ie *ast.MatchExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	for _, match := range ie.MatchBlock.Cases {
		if match.Condition != nil {
			caseCondition := Eval(match.Condition, env)
			if isError(caseCondition) {
				return caseCondition
			}
			if caseCondition.Type() == condition.Type() {
				if isEqual(condition, caseCondition) {
					return Eval(match.Consequence, env)
				}
			} else if caseCondition.Type() == object.BOOLEAN_OBJ {
				if isTruthy(caseCondition) {
					return Eval(match.Consequence, env)
				}
			} else if isTruthy(caseCondition) {
				return Eval(match.Consequence, env)
			}
		}
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	msg := fmt.Sprintf("[%d:%d] identifier not found: %s", node.Token.Position().Line(), node.Token.Position().Column(), node.Value)

	return newError(msg)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)

	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func evalHashIndexExpression(str, index object.Object) object.Object {
	hashObject := str.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalFStringLiteral(node *ast.FStringLiteral, env *object.Environment) object.Object {
	newEnv := object.NewEnclosedEnvironment(env)

	strCopy := node.Value
	start := 0
	for {
		openBrace := strings.Index(strCopy[start:], "{")
		if openBrace == -1 {
			break
		}
		openBrace += start
		closeBrace := strings.Index(strCopy[openBrace:], "}")
		if closeBrace == -1 {
			return newError("unmatched '{' found in the string")
		}
		closeBrace += openBrace

		placeholder := strCopy[openBrace+1 : closeBrace]
		p := parser.New(lexer.New(bufio.NewReader(strings.NewReader(placeholder))))
		x := p.Parse()
		if len(p.Errors()) > 0 {
			return newError("error parsing placeholder: %s", p.Errors())
		}

		value := Eval(x, newEnv)
		if isError(value) {
			return value
		}

		strCopy = strCopy[:openBrace] + value.Inspect() + strCopy[closeBrace+1:]
		start = openBrace + len(value.Inspect())
	}

	return &object.String{Value: strCopy}
}

func evalLoopStatement(loop *ast.LoopStatement, env *object.Environment) object.Object {
	for {
		if loop.WhileCondition != nil {
			condition := Eval(loop.WhileCondition, env)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		if loop.UntilCondition != nil {
			condition := Eval(loop.UntilCondition, env)
			if isError(condition) {
				return condition
			}
			if isTruthy(condition) {
				break
			}
		}

		shouldBreak := false
		shouldContinue := false
		for _, statement := range loop.Body.Statements {
			stmtResult := Eval(statement, env)
			if stmtResult != nil {
				if isError(stmtResult) {
					return stmtResult
				}
				if stmtResult.Type() == object.BREAK_OBJ {
					shouldBreak = true
					break
				}
				if stmtResult.Type() == object.CONTINUE_OBJ {
					shouldContinue = true
					break
				}
				if stmtResult.Type() == object.RETURN_VALUE_OBJ {
					return stmtResult
				}
			}
		}
		if shouldBreak {
			break
		}
		if shouldContinue {
			continue
		}
	}
	return VOID
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isEqual(obj1 object.Object, obj2 object.Object) bool {
	if obj1.Type() != obj2.Type() {
		return false
	}
	if obj1.Inspect() != obj2.Inspect() {
		return false
	}
	return true
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

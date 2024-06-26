package evaluator

import (
	"fmt"
	"time"

	"github.com/poolpOrg/julu/object"
)

var builtins = map[string]*object.Builtin{
	"type": {
		Fn: builtin_type,
	},

	"len": {
		Fn: builtin_len,
	},
	"println": {
		Fn: builtin_println,
	},

	// TEMPORARY: This is a temporary function to test the evaluator.
	"sleep": {
		Fn: builtin_sleep,
	},
}

func builtin_type(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
	}
	return &object.String{Value: string(args[0].Type())}
}

func builtin_len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
	}
}

func builtin_println(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return nil
}

// TEMPORARY: This is a temporary function to test the evaluator.

func builtin_sleep(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return &object.Error{Message: fmt.Sprintf("argument to `sleep` must be INTEGER, got %s", args[0].Type())}
	}

	time.Sleep(time.Duration(args[0].(*object.Integer).Value) * time.Second)
	return nil
}

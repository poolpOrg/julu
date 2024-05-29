package evaluator

import (
	"fmt"

	"github.com/poolpOrg/julu/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: builtin_len,
	},
	"println": {
		Fn: builtin_println,
	},
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

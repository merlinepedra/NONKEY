package evaluator

import (
	"os"

	"github.com/kasworld/nonkey/enum/objecttype"
	"github.com/kasworld/nonkey/interpreter/asti"
	"github.com/kasworld/nonkey/interpreter/object"
)

// os.getenv() -> ( Hash )
func builtinOsEnvironment(node asti.NodeI, env *object.Environment, args ...object.ObjectI) object.ObjectI {

	osenv := os.Environ()
	newHash := make(map[object.HashKey]object.HashPair)

	//
	// If we get a match then the output is an array
	// First entry is the match, any additional parts
	// are the capture-groups.
	//
	for i := 1; i < len(osenv); i++ {

		// Capture groups start at index 0.
		k := &object.String{Value: osenv[i]}
		v := &object.String{Value: os.Getenv(osenv[i])}

		newHashPair := object.HashPair{Key: k, Value: v}
		newHash[k.HashKey()] = newHashPair
	}

	return &object.Hash{Pairs: newHash}
}

// os.getenv( "PATH" ) -> string
func builtinOsGetEnv(node asti.NodeI, env *object.Environment, args ...object.ObjectI) object.ObjectI {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != objecttype.STRING {
		return object.NewError(node, "argument must be a string, got=%s",
			args[0].Type())
	}
	input := args[0].(*object.String).Value
	return &object.String{Value: os.Getenv(input)}

}

// os.setenv( "PATH", "/home/skx/bin:/usr/bin" );
func builtinOsSetEnv(node asti.NodeI, env *object.Environment, args ...object.ObjectI) object.ObjectI {
	if len(args) != 2 {
		return object.NewError(node, "wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != objecttype.STRING {
		return object.NewError(node, "argument must be a string, got=%s",
			args[0].Type())
	}
	if args[1].Type() != objecttype.STRING {
		return object.NewError(node, "argument must be a string, got=%s",
			args[1].Type())
	}
	name := args[0].(*object.String).Value
	value := args[1].(*object.String).Value
	os.Setenv(name, value)
	return object.NULL
}

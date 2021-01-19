package evaluator

import (
	"path/filepath"

	"github.com/kasworld/nonkey/interpreter/asti"
	"github.com/kasworld/nonkey/interpreter/object"
)

// array = directory.glob( "/etc/*.conf" )
func builtinDirectoryGlob(node asti.NodeI, env *object.Environment, args ...object.ObjectI) object.ObjectI {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got=%d, want=1",
			len(args))
	}
	pattern := args[0].(*object.String).Value

	entries, err := filepath.Glob(pattern)
	if err != nil {
		return object.NULL
	}

	// Create an array to hold the results and populate it
	l := len(entries)
	result := make([]object.ObjectI, l)
	for i, txt := range entries {
		result[i] = &object.String{Value: txt}
	}
	return &object.Array{Elements: result}
}

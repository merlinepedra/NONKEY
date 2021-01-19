package evaluator

import (
	"os"

	"github.com/kasworld/nonkey/interpreter/asti"
	"github.com/kasworld/nonkey/interpreter/object"
	"github.com/kasworld/version"
)

//
// Implemention of "version()" function.
//
func builtinVersion(node asti.NodeI, env *object.Environment, args ...object.ObjectI) object.ObjectI {
	return &object.String{Value: version.GetVersion()}
}

//
// Implemention of "args()" function.
//
func builtinArgs(node asti.NodeI, env *object.Environment, args ...object.ObjectI) object.ObjectI {
	l := len(os.Args[1:])
	result := make([]object.ObjectI, l)
	for i, txt := range os.Args[1:] {
		result[i] = &object.String{Value: txt}
	}
	return &object.Array{Elements: result}
}

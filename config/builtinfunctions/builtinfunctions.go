package builtinfunctions

import "github.com/kasworld/nonkey/interpreter/object"

// The built-in functions / standard-library methods are stored here.
var BuiltinFunctions = map[string]*object.Builtin{}

// Register registers a built-in function.  This is used to register
// our "standard library" functions.
func Register(name string, fun object.BuiltinFunction) {
	BuiltinFunctions[name] = &object.Builtin{Fn: fun}
}

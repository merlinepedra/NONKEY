package object

import (
	"github.com/kasworld/nonkey/enum/objecttype"
	"github.com/kasworld/nonkey/interpreter/asti"
)

// BuiltinFunction holds the type of a built-in function.
type BuiltinFunction func(node asti.NodeI, env *Environment, args ...ObjectI) ObjectI

// Builtin wraps func and implements ObjectI interface.
type Builtin struct {
	// Value holds the function we wrap.
	Fn BuiltinFunction
}

// Type returns the type of this object.
func (b *Builtin) Type() objecttype.ObjectType {
	return objecttype.BUILTIN
}

// Inspect returns a string-representation of the given object.
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (b *Builtin) InvokeMethod(method string, env Environment, args ...ObjectI) ObjectI {
	if method == "methods" {
		names := []string{"methods"}

		result := make([]ObjectI, len(names))
		for i, txt := range names {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	}
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Builtin) ToInterface() interface{} {
	return "<BUILTIN>"
}

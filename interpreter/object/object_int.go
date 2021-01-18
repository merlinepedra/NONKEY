package object

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kasworld/nonkey/enum/objecttype"
)

// Integer wraps int64 and implements ObjectI and Hashable interfaces.
type Integer struct {
	// Value holds the integer value this object wraps
	Value int64
}

// Inspect returns a string-representation of the given object.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type returns the type of this object.
func (i *Integer) Type() objecttype.ObjectType {
	return objecttype.INTEGER
}

// HashKey returns a hash key for the given object.
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (i *Integer) InvokeMethod(method string, env Environment, args ...ObjectI) ObjectI {
	if method == "chr" {
		return &String{Value: string(rune(i.Value))}
	}
	if method == "methods" {
		static := []string{"chr", "methods"}
		dynamic := env.Names("integer.")

		var names []string
		names = append(names, static...)

		for _, e := range dynamic {
			bits := strings.Split(e, ".")
			names = append(names, bits[1])
		}
		sort.Strings(names)

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
func (i *Integer) ToInterface() interface{} {
	return i.Value
}

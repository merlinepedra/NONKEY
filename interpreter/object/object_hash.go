package object

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/kasworld/nonkey/enum/objecttype"
)

// HashKey is the structure used for hash-keys
type HashKey struct {
	// Type holds the type of the object.
	Type objecttype.ObjectType

	// Value holds the actual hash-key value.
	Value uint64
}

// HashPair is a structure which is used to store hash-entries
type HashPair struct {
	// Key holds our hash-key key.
	Key ObjectI

	// Value holds our hash-key value.
	Value ObjectI
}

// Hash wrap map[HashKey]HashPair and implements ObjectI interface.
type Hash struct {
	// Pairs holds the key/value pairs of the hash we wrap
	Pairs map[HashKey]HashPair

	// offset holds our iteration-offset.
	offset int
}

// Type returns the type of this object.
func (h *Hash) Type() objecttype.ObjectType {
	return objecttype.HASH
}

// Inspect returns a string-representation of the given object.
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	fmt.Fprintf(&out, "{%v}", strings.Join(pairs, ", "))
	return out.String()
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (h *Hash) InvokeMethod(method string, env Environment, args ...ObjectI) ObjectI {
	if method == "keys" {
		ents := len(h.Pairs)
		array := make([]ObjectI, ents)

		// Now copy the keys into it.
		i := 0
		for _, ent := range h.Pairs {
			array[i] = ent.Key
			i++
		}

		return &Array{Elements: array}
	}
	if method == "methods" {
		static := []string{"keys", "methods"}
		dynamic := env.Names("hash.")

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

// Reset implements the Iterable interface, and allows the contents
// of the array to be reset to allow re-iteration.
func (h *Hash) Reset() {
	h.offset = 0
}

// Next implements the Iterable interface, and allows the contents
// of our array to be iterated over.
func (h *Hash) Next() (ObjectI, ObjectI, bool) {
	if h.offset < len(h.Pairs) {
		idx := 0

		for _, pair := range h.Pairs {
			if h.offset == idx {
				h.offset++
				return pair.Key, pair.Value, true
			}
			idx++
		}
	}

	return nil, &Integer{Value: 0}, false
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (h *Hash) ToInterface() interface{} {
	return "<HASH>"
}

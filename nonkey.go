// Monkey is a scripting language implemented in golang, based upon
// the book "Write an Interpreter in Go", written by Thorsten Ball.
//
// This implementation adds a number of tweaks, improvements, and new
// features.  For example we support file-based I/O, regular expressions,
// the ternary operator, and more.
//
// For full details please consult the project homepage https://github.com/skx/monkey/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kasworld/nonkey/evaluator"
	"github.com/kasworld/nonkey/object"
	"github.com/kasworld/nonkey/repl"
	"github.com/kasworld/nonkey/runmon"
)

// This version-string will be updated via travis for generated binaries.
var version = "master/unreleased"

//
// Implemention of "version()" function.
//
func builtinVersion(env *object.Environment, args ...object.Object) object.Object {
	return &object.String{Value: version}
}

//
// Implemention of "args()" function.
//
func builtinArgs(env *object.Environment, args ...object.Object) object.Object {
	l := len(os.Args[1:])
	result := make([]object.Object, l)
	for i, txt := range os.Args[1:] {
		result[i] = &object.String{Value: txt}
	}
	return &object.Array{Elements: result}
}

func main() {
	eval := flag.String("eval", "", "Code to execute.")
	vers := flag.Bool("version", false, "Show our version and exit.")
	flag.Parse()

	// show version
	if *vers {
		fmt.Printf("monkey %s\n", version)
		os.Exit(1)
	}

	evaluator.RegisterBuiltin("version", builtinVersion)
	evaluator.RegisterBuiltin("args", builtinArgs)

	env := object.NewEnvironment()
	env = runmon.RunFile("data/stdlib.mon", env)

	if *eval != "" { // run 1 line
		runmon.RunString(*eval, env)
		os.Exit(1)
	} else {
		if len(flag.Args()) > 0 { // run file
			runmon.RunFile(os.Args[1], env)
		} else { // repl line by line
			repl.Start(os.Stdin, os.Stdout, env)
		}
	}
}

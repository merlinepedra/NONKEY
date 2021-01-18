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

	"github.com/kasworld/nonkey/config/builtinfunctions"
	"github.com/kasworld/nonkey/interpreter/object"
	"github.com/kasworld/nonkey/interpreter/repl"
	"github.com/kasworld/nonkey/interpreter/runmon"
)

// This version-string will be updated via travis for generated binaries.
var version = "master/unreleased"

//
// Implemention of "version()" function.
//
func builtinVersion(env *object.Environment, args ...object.ObjectI) object.ObjectI {
	return &object.String{Value: version}
}

//
// Implemention of "args()" function.
//
func builtinArgs(env *object.Environment, args ...object.ObjectI) object.ObjectI {
	l := len(os.Args[1:])
	result := make([]object.ObjectI, l)
	for i, txt := range os.Args[1:] {
		result[i] = &object.String{Value: txt}
	}
	return &object.Array{Elements: result}
}

func main() {
	eval := flag.String("eval", "", "Code to execute.")
	vers := flag.Bool("version", false, "Show our version and exit.")
	autoload := flag.String("autoload", "", "autoload filename")
	flag.Parse()

	// show version
	if *vers {
		fmt.Printf("monkey %s\n", version)
		os.Exit(1)
	}

	builtinfunctions.Register("version", builtinVersion)
	builtinfunctions.Register("args", builtinArgs)

	env := object.NewEnvironment()
	if *autoload != "" {
		fmt.Printf("autoload %v\n", *autoload)
		env = runmon.RunFile(*autoload, env)
	}

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

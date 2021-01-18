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

	"github.com/kasworld/nonkey/interpreter/object"
	"github.com/kasworld/nonkey/interpreter/repl"
	"github.com/kasworld/nonkey/interpreter/runmon"
	"github.com/kasworld/version"
)

var Ver = "dev_notbuild"

func init() {
	version.Set(Ver)
}

func main() {
	eval := flag.String("eval", "", "Code to execute.")
	vers := flag.Bool("version", false, "Show our version and exit.")
	autoload := flag.String("autoload", "", "autoload filename")
	flag.Parse()

	// show version
	if *vers {
		fmt.Printf("monkey %s\n", Ver)
		os.Exit(1)
	}

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

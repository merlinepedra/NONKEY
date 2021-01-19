package runmon

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kasworld/nonkey/interpreter/evaluator"
	"github.com/kasworld/nonkey/interpreter/lexer"
	"github.com/kasworld/nonkey/interpreter/object"
	"github.com/kasworld/nonkey/interpreter/parser"
)

func RunFile(filename string, env *object.Environment) *object.Environment {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to load %v %v\n", filename, err)
		return env
	}
	return RunString(string(input), env)
}

func RunString(input string, env *object.Environment) *object.Environment {
	initL := lexer.New(input)
	initP := parser.New(initL)
	initProg := initP.ParseProgram()

	if len(initP.Errors()) != 0 {
		for _, v := range initP.Errors() {
			fmt.Fprintf(os.Stderr, "%v\n", v)
		}
		return env
	}

	evaluator.Eval(initProg, env)
	return env
}

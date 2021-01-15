package runmon

import (
	"fmt"
	"io/ioutil"

	"github.com/kasworld/nonkey/evaluator"
	"github.com/kasworld/nonkey/lexer"
	"github.com/kasworld/nonkey/object"
	"github.com/kasworld/nonkey/parser"
)

func RunFile(filename string, env *object.Environment) *object.Environment {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("fail to load %v %v", filename, err)
		return env
	}
	return RunString(string(input), env)
}

func RunString(input string, env *object.Environment) *object.Environment {
	initL := lexer.New(input)
	initP := parser.New(initL)
	initProg := initP.ParseProgram()
	evaluator.Eval(initProg, env)
	return env
}

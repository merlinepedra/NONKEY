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
	l := lexer.New(input)
	p := parser.New(l)
	prg := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, v := range p.Errors() {
			fmt.Fprintf(os.Stderr, "%v\n", v)
		}
		return env
	}

	evaluated := evaluator.Eval(prg, env)
	if evaluated != nil {
		if erro, ok := evaluated.(*object.Error); ok {
			fmt.Fprintf(os.Stderr, "%v\n", evaluated.Inspect())
			fmt.Fprintf(os.Stderr, "%v\n", l.GetLineStr(erro.Node.GetToken().Line))
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", evaluated.Inspect())
		}
	}

	return env
}

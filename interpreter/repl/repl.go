package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kasworld/nonkey/interpreter/evaluator"
	"github.com/kasworld/nonkey/interpreter/lexer"
	"github.com/kasworld/nonkey/interpreter/object"
	"github.com/kasworld/nonkey/interpreter/parser"
	"github.com/kasworld/version"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, env *object.Environment) {

	fmt.Fprintf(out, "welcome to nonkey version:%v\n", version.GetVersion())

	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if erro, ok := evaluated.(*object.Error); ok {
				fmt.Fprintf(out, "%v\n", evaluated.Inspect())
				fmt.Fprintf(out, "%v\n", l.GetLineStr(erro.Node.GetToken().Line))
			} else {
				fmt.Fprintf(out, "%v\n", evaluated.Inspect())
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []parser.Error) {
	io.WriteString(out, "parser errors\n")
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%v\n", msg)
	}
}

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

	io.WriteString(out, "welcome to nonkey version:")
	io.WriteString(out, version.GetVersion())
	io.WriteString(out, "\n")

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
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []parser.Error) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.String()+"\n")
	}
}

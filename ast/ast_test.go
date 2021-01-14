package ast

import (
	"testing"

	"github.com/kasworld/nonkey/token"
	"github.com/kasworld/nonkey/tokentype"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: tokentype.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: tokentype.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: tokentype.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

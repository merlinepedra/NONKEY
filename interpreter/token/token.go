// Package token contains constants which are used when lexing a program
// written in the monkey language, as done by the parser.
package token

import (
	"fmt"

	"github.com/kasworld/nonkey/enum/tokentype"
)

// Token struct represent the lexer token
type Token struct {
	Type    tokentype.TokenType
	Literal string

	Line int // token pos line in source code
	Pos  int // token pos of source code
}

func (tk Token) String() string {
	return fmt.Sprintf("Token[%v %v at line %v pos %v]",
		tk.Type, tk.Literal, tk.Line+1, tk.Pos+1)
}

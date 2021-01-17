// Package token contains constants which are used when lexing a program
// written in the monkey language, as done by the parser.
package token

import (
	"github.com/kasworld/nonkey/enum/tokentype"
)

// Token struct represent the lexer token
type Token struct {
	Type    tokentype.TokenType
	Literal string
	Pos     int // token post of source code
}

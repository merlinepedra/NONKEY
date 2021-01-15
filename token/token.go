// Package token contains constants which are used when lexing a program
// written in the monkey language, as done by the parser.
package token

import "github.com/kasworld/nonkey/tokentype"

// Token struct represent the lexer token
type Token struct {
	Type    tokentype.TokenType
	Literal string
	Pos     int // token post of source code
}

// reversed keywords
var keywords = map[string]tokentype.TokenType{
	"case":     tokentype.CASE,
	"const":    tokentype.CONST,
	"default":  tokentype.DEFAULT,
	"else":     tokentype.ELSE,
	"false":    tokentype.FALSE,
	"fn":       tokentype.FUNCTION,
	"for":      tokentype.FOR,
	"foreach":  tokentype.FOREACH,
	"function": tokentype.DEFINE_FUNCTION,
	"if":       tokentype.IF,
	"in":       tokentype.IN,
	"let":      tokentype.LET,
	"return":   tokentype.RETURN,
	"switch":   tokentype.SWITCH,
	"true":     tokentype.TRUE,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) tokentype.TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return tokentype.IDENT
}

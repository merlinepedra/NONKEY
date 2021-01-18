package tokentype

import "github.com/kasworld/nonkey/enum/precedence"

// Keywords reversed keywords
var Keywords = map[string]TokenType{
	"case":     CASE,
	"const":    CONST,
	"default":  DEFAULT,
	"else":     ELSE,
	"false":    FALSE,
	"fn":       FUNCTION,
	"for":      FOR,
	"foreach":  FOREACH,
	"function": DEFINE_FUNCTION,
	"if":       IF,
	"in":       IN,
	"let":      LET,
	"return":   RETURN,
	"switch":   SWITCH,
	"true":     TRUE,
}

// LookupKeyword used to determinate whether identifier is keyword nor not
func LookupKeyword(identifier string) TokenType {
	if tok, ok := Keywords[identifier]; ok {
		return tok
	}
	return IDENT
}

// Token2Precedences each token precedence
var Token2Precedences = [TokenType_Count]precedence.Precedence{
	QUESTION:     precedence.TERNARY,
	ASSIGN:       precedence.ASSIGN,
	DOTDOT:       precedence.DOTDOT,
	EQ:           precedence.EQUALS,
	NOT_EQ:       precedence.EQUALS,
	LT:           precedence.LESSGREATER,
	LT_EQUALS:    precedence.LESSGREATER,
	GT:           precedence.LESSGREATER,
	GT_EQUALS:    precedence.LESSGREATER,
	CONTAINS:     precedence.REGEXP_MATCH,
	NOT_CONTAINS: precedence.REGEXP_MATCH,

	PLUS:            precedence.SUM,
	PLUS_EQUALS:     precedence.SUM,
	MINUS:           precedence.SUM,
	MINUS_EQUALS:    precedence.SUM,
	SLASH:           precedence.PRODUCT,
	SLASH_EQUALS:    precedence.PRODUCT,
	ASTERISK:        precedence.PRODUCT,
	ASTERISK_EQUALS: precedence.PRODUCT,
	POW:             precedence.POWER,
	MOD:             precedence.MOD,
	AND:             precedence.COND,
	OR:              precedence.COND,
	LPAREN:          precedence.CALL,
	PERIOD:          precedence.CALL,
	LBRACKET:        precedence.INDEX,
}

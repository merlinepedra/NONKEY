package tokentype

import "github.com/kasworld/nonkey/enum/precedence"

func (tk TokenType) IsKeyword() bool {
	return attrib[tk].keyword
}

func (tk TokenType) IsOperator() bool {
	return attrib[tk].operator
}

func (tk TokenType) Literal() string {
	return attrib[tk].literal
}

var attrib = [TokenType_Count]struct {
	keyword  bool
	operator bool
	literal  string
}{
	ILLEGAL:         {false, false, "ILLEGAL"},
	REGEXP:          {false, false, "REGEXP"},
	EOF:             {false, false, "EOF"},
	IDENT:           {false, false, "IDENT"},
	CASE:            {true, false, "case"},
	CONST:           {true, false, "const"},
	DEFAULT:         {true, false, "default"},
	DEFINE_FUNCTION: {true, false, "function"},
	ELSE:            {true, false, "else"},
	STRING:          {true, false, "string"},
	SWITCH:          {true, false, "switch"},
	TRUE:            {true, false, "true"},
	FALSE:           {true, false, "false"},
	FLOAT:           {true, false, "float"},
	FOR:             {true, false, "for"},
	FOREACH:         {true, false, "foreach"},
	FUNCTION:        {true, false, "fn"},
	IF:              {true, false, "if"},
	IN:              {true, false, "in"},
	INT:             {true, false, "int"},
	LET:             {true, false, "let"},
	RETURN:          {true, false, "return"},
	AND:             {false, true, "&&"},
	ASSIGN:          {false, true, "="},
	ASTERISK:        {false, true, "*"},
	ASTERISK_EQUALS: {false, true, "*="},
	BACKTICK:        {false, true, "`"},
	BANG:            {false, true, "!"},
	COLON:           {false, true, ":"},
	COMMA:           {false, true, ","},
	CONTAINS:        {false, true, "~="},
	DOTDOT:          {false, true, ".."},
	EQ:              {false, true, "=="},
	GT:              {false, true, ">"},
	GT_EQUALS:       {false, true, ">="},
	LBRACE:          {false, true, "{"},
	LBRACKET:        {false, true, "["},
	LPAREN:          {false, true, "("},
	LT:              {false, true, "<"},
	LT_EQUALS:       {false, true, "<="},
	MINUS:           {false, true, "-"},
	MINUS_EQUALS:    {false, true, "-="},
	MINUS_MINUS:     {false, true, "--"},
	MOD:             {false, true, "%"},
	NOT_CONTAINS:    {false, true, "!~"},
	NOT_EQ:          {false, true, "!="},
	OR:              {false, true, "||"},
	PERIOD:          {false, true, "."},
	PLUS:            {false, true, "+"},
	PLUS_EQUALS:     {false, true, "+="},
	PLUS_PLUS:       {false, true, "++"},
	POW:             {false, true, "**"},
	QUESTION:        {false, true, "?"},
	RBRACE:          {false, true, "}"},
	RBRACKET:        {false, true, "]"},
	RPAREN:          {false, true, ")"},
	SEMICOLON:       {false, true, ";"},
	SLASH:           {false, true, "/"},
	SLASH_EQUALS:    {false, true, "/="},
}

// Keywords reversed keywords
var Keywords = map[string]TokenType{}

func init() {
	// build keyword map
	for i, v := range attrib {
		if v.keyword {
			Keywords[v.literal] = TokenType(i)
		}
	}
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

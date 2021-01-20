package tokentype

import "github.com/kasworld/nonkey/enum/precedence"

func (tk TokenType) IsKeyword() bool {
	return attrib[tk].keyword
}

func (tk TokenType) Literal() string {
	return attrib[tk].literal
}

var attrib = [TokenType_Count]struct {
	keyword bool
	literal string
}{
	ILLEGAL: {false, "ILLEGAL"},
	REGEXP:  {false, "REGEXP"},
	EOF:     {false, "EOF"},
	EOL:     {false, "EOL"},
	IDENT:   {false, "IDENT"},

	// keyword
	CASE:            {true, "case"},
	CONST:           {true, "const"},
	DEFAULT:         {true, "default"},
	DEFINE_FUNCTION: {true, "function"},
	ELSE:            {true, "else"},
	STRING:          {true, "string"},
	SWITCH:          {true, "switch"},
	TRUE:            {true, "true"},
	FALSE:           {true, "false"},
	FLOAT:           {true, "float"},
	FOR:             {true, "for"},
	FOREACH:         {true, "foreach"},
	FUNCTION:        {true, "fn"},
	IF:              {true, "if"},
	IN:              {true, "in"},
	INT:             {true, "int"},
	LET:             {true, "let"},
	RETURN:          {true, "return"},

	BACKTICK:    {false, "`"},
	BANG:        {false, "!"},
	COLON:       {false, ":"},
	COMMA:       {false, ","},
	LBRACE:      {false, "{"},
	MINUS_MINUS: {false, "--"},
	PLUS_PLUS:   {false, "++"},
	RBRACE:      {false, "}"},
	RBRACKET:    {false, "]"},
	RPAREN:      {false, ")"},
	SEMICOLON:   {false, ";"},

	// precedence
	QUESTION:        {false, "?"},
	ASSIGN:          {false, "="},
	DOTDOT:          {false, ".."},
	EQ:              {false, "=="},
	NOT_EQ:          {false, "!="},
	LT:              {false, "<"},
	LT_EQUALS:       {false, "<="},
	GT:              {false, ">"},
	GT_EQUALS:       {false, ">="},
	CONTAINS:        {false, "~="},
	NOT_CONTAINS:    {false, "!~"},
	PLUS:            {false, "+"},
	PLUS_EQUALS:     {false, "+="},
	MINUS:           {false, "-"},
	MINUS_EQUALS:    {false, "-="},
	SLASH:           {false, "/"},
	SLASH_EQUALS:    {false, "/="},
	ASTERISK:        {false, "*"},
	ASTERISK_EQUALS: {false, "*="},
	POW:             {false, "**"},
	MOD:             {false, "%"},
	AND:             {false, "&&"},
	OR:              {false, "||"},
	LPAREN:          {false, "("},
	PERIOD:          {false, "."},
	LBRACKET:        {false, "["},
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

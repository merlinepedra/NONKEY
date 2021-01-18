package lexer

import (
	"testing"

	"github.com/kasworld/nonkey/enum/tokentype"
)

func TestNextToken1(t *testing.T) {
	input := "%=+(){},;?|| &&`/bin/ls`++--***=.."

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.MOD, "%"},
		{tokentype.ASSIGN, "="},
		{tokentype.PLUS, "+"},
		{tokentype.LPAREN, "("},
		{tokentype.RPAREN, ")"},
		{tokentype.LBRACE, "{"},
		{tokentype.RBRACE, "}"},
		{tokentype.COMMA, ","},
		{tokentype.SEMICOLON, ";"},
		{tokentype.QUESTION, "?"},
		{tokentype.OR, "||"},
		{tokentype.AND, "&&"},
		{tokentype.BACKTICK, "/bin/ls"},
		{tokentype.PLUS_PLUS, "++"},
		{tokentype.MINUS_MINUS, "--"},
		{tokentype.POW, "**"},
		{tokentype.ASTERISK_EQUALS, "*="},
		{tokentype.DOTDOT, ".."},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken2(t *testing.T) {
	input := `let five=5;
let ten =10;
let add = fn(x, y){
  x+y;
};
let result = add(five, ten);
!- *5;
5<10>5;

if(5<10){
	return true;
}else{
	return false;
}
10 == 10;
10 != 9;
"foobar"
"foo bar"
[1,2];
{"foo":"bar"}
1.2
0.5
0.3
世界
for
2 >= 1
1 <= 3
empty?
`
	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.LET, "let"},
		{tokentype.IDENT, "five"},
		{tokentype.ASSIGN, "="},
		{tokentype.INT, "5"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.LET, "let"},
		{tokentype.IDENT, "ten"},
		{tokentype.ASSIGN, "="},
		{tokentype.INT, "10"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.LET, "let"},
		{tokentype.IDENT, "add"},
		{tokentype.ASSIGN, "="},
		{tokentype.FUNCTION, "fn"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "x"},
		{tokentype.COMMA, ","},
		{tokentype.IDENT, "y"},
		{tokentype.RPAREN, ")"},
		{tokentype.LBRACE, "{"},
		{tokentype.IDENT, "x"},
		{tokentype.PLUS, "+"},
		{tokentype.IDENT, "y"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.RBRACE, "}"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.LET, "let"},
		{tokentype.IDENT, "result"},
		{tokentype.ASSIGN, "="},
		{tokentype.IDENT, "add"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "five"},
		{tokentype.COMMA, ","},
		{tokentype.IDENT, "ten"},
		{tokentype.RPAREN, ")"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.BANG, "!"},
		{tokentype.MINUS, "-"},
		{tokentype.ASTERISK, "*"},
		{tokentype.INT, "5"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.INT, "5"},
		{tokentype.LT, "<"},
		{tokentype.INT, "10"},
		{tokentype.GT, ">"},
		{tokentype.INT, "5"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.INT, "5"},
		{tokentype.LT, "<"},
		{tokentype.INT, "10"},
		{tokentype.RPAREN, ")"},
		{tokentype.LBRACE, "{"},
		{tokentype.RETURN, "return"},
		{tokentype.TRUE, "true"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.RBRACE, "}"},
		{tokentype.ELSE, "else"},
		{tokentype.LBRACE, "{"},
		{tokentype.RETURN, "return"},
		{tokentype.FALSE, "false"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.RBRACE, "}"},
		{tokentype.INT, "10"},
		{tokentype.EQ, "=="},
		{tokentype.INT, "10"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.INT, "10"},
		{tokentype.NOT_EQ, "!="},
		{tokentype.INT, "9"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.STRING, "foobar"},
		{tokentype.STRING, "foo bar"},
		{tokentype.LBRACKET, "["},
		{tokentype.INT, "1"},
		{tokentype.COMMA, ","},
		{tokentype.INT, "2"},
		{tokentype.RBRACKET, "]"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.LBRACE, "{"},
		{tokentype.STRING, "foo"},
		{tokentype.COLON, ":"},
		{tokentype.STRING, "bar"},
		{tokentype.RBRACE, "}"},
		{tokentype.FLOAT, "1.2"},
		{tokentype.FLOAT, "0.5"},
		{tokentype.FLOAT, "0.3"},
		{tokentype.IDENT, "世界"},
		{tokentype.FOR, "for"},
		{tokentype.INT, "2"},
		{tokentype.GT_EQUALS, ">="},
		{tokentype.INT, "1"},
		{tokentype.INT, "1"},
		{tokentype.LT_EQUALS, "<="},
		{tokentype.INT, "3"},
		{tokentype.IDENT, "empty?"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestUnicodeLexer(t *testing.T) {
	input := `世界`
	l := New(input)
	tok := l.NextToken()
	if tok.Type != tokentype.IDENT {
		t.Fatalf("token type wrong, expected=%q, got=%q", tokentype.IDENT, tok.Type)
	}
	if tok.Literal != "世界" {
		t.Fatalf("token literal wrong, expected=%q, got=%q", "世界", tok.Literal)
	}
}

func TestSimpleComment(t *testing.T) {
	input := `=+// This is a comment
// This is still a comment
# I like comments
let a = 1; # This is a comment too.
// This is a final
// comment on two-lines`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.ASSIGN, "="},
		{tokentype.PLUS, "+"},
		{tokentype.LET, "let"},
		{tokentype.IDENT, "a"},
		{tokentype.ASSIGN, "="},
		{tokentype.INT, "1"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestMultiLineComment(t *testing.T) {
	input := `=+/* This is a comment

We're still in a comment
let c = 2; */
let a = 1;
// This isa comment
// This is still a comment.
/* Now a multi-line again
   Which is two-lines
 */`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.ASSIGN, "="},
		{tokentype.PLUS, "+"},
		{tokentype.LET, "let"},
		{tokentype.IDENT, "a"},
		{tokentype.ASSIGN, "="},
		{tokentype.INT, "1"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIntegers(t *testing.T) {
	input := `10 0x10 0xF0 0xFE 0b0101 0xFF 0b101 0xFF;`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.INT, "10"},
		{tokentype.INT, "0x10"},
		{tokentype.INT, "0xF0"},
		{tokentype.INT, "0xFE"},
		{tokentype.INT, "0b0101"},
		{tokentype.INT, "0xFF"},
		{tokentype.INT, "0b101"},
		{tokentype.INT, "0xFF"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test that the shebang-line is handled specially.
func TestShebang(t *testing.T) {
	input := `#!/bin/monkey
10;`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.INT, "10"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestMoreHandling does nothing real, but it bumps our coverage!
func TestMoreHandling(t *testing.T) {
	input := `#!/bin/monkey
1 += 1;
2 -= 2;
3 /= 3;
x */ 3;

let t = true;
let f = false;

if ( t && f ) { puts( "What?" ); }
if ( t || f ) { puts( "What?" ); }

let a = 1;
a++;

let b = a % 1;
b--;
b -= 2;

if ( a<3 ) { puts( "Blah!"); }
if ( a>3 ) { puts( "Blah!"); }

let b = 3;
b**b;
b *= 3;
if ( b <= 3  ) { puts "blah\n" }
if ( b >= 3  ) { puts "blah\n" }

let a = "steve";
let a = "steve\n";
let a = "steve\t";
let a = "steve\r";
let a = "steve\\";
let a = "steve\"";
let c = 3.113;
.;`

	l := New(input)
	tok := l.NextToken()
	for tok.Type != tokentype.EOF {

		tok = l.NextToken()
	}
}

// TestStdLib ensures that identifiers are parsed correctly for the
// case where we need to support the legacy-names.
func TestStdLib(t *testing.T) {
	input := `
os.getenv
os.setenv
os.environment
directory.glob
math.abs
math.random
math.sqrt
string.interpolate
string.toupper
string.tolower
string.trim
string.reverse
string.split
moi.kissa
`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.IDENT, "os.getenv"},
		{tokentype.IDENT, "os.setenv"},
		{tokentype.IDENT, "os.environment"},
		{tokentype.IDENT, "directory.glob"},
		{tokentype.IDENT, "math.abs"},
		{tokentype.IDENT, "math.random"},
		{tokentype.IDENT, "math.sqrt"},
		{tokentype.IDENT, "string.interpolate"},
		{tokentype.IDENT, "string.toupper"},
		{tokentype.IDENT, "string.tolower"},
		{tokentype.IDENT, "string.trim"},
		{tokentype.IDENT, "string.reverse"},
		{tokentype.IDENT, "string.split"},
		{tokentype.IDENT, "moi"},
		{tokentype.PERIOD, "."},
		{tokentype.IDENT, "kissa"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%v, got=%q", i, tt, tok)
		}
	}
}

// TestDotMethod ensures that identifiers are parsed correctly for the
// case where we need to split at periods.
func TestDotMethod(t *testing.T) {
	input := `
foo.bar();
moi.kissa();
a?.b?();
`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.IDENT, "foo"},
		{tokentype.PERIOD, "."},
		{tokentype.IDENT, "bar"},
		{tokentype.LPAREN, "("},
		{tokentype.RPAREN, ")"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.IDENT, "moi"},
		{tokentype.PERIOD, "."},
		{tokentype.IDENT, "kissa"},
		{tokentype.LPAREN, "("},
		{tokentype.RPAREN, ")"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.IDENT, "a?"},
		{tokentype.PERIOD, "."},
		{tokentype.IDENT, "b?"},
		{tokentype.LPAREN, "("},
		{tokentype.RPAREN, ")"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%v, got=%q", i, tt, tok)
		}
	}
}

// TestIntDotMethod ensures that identifiers are parsed correctly for the
// case where they immediately follow int/float valies.
func TestIntDotMethod(t *testing.T) {
	input := `
3.foo();
3.14.bar();
`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.INT, "3"},
		{tokentype.PERIOD, "."},
		{tokentype.IDENT, "foo"},
		{tokentype.LPAREN, "("},
		{tokentype.RPAREN, ")"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.FLOAT, "3.14"},
		{tokentype.PERIOD, "."},
		{tokentype.IDENT, "bar"},
		{tokentype.LPAREN, "("},
		{tokentype.RPAREN, ")"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%v, got=%q", i, tt, tok)
		}
	}
}

// TestRegexp ensures a simple regexp can be parsed.
func TestRegexp(t *testing.T) {
	input := `if ( f ~= /steve/i )
if ( f ~= /steve/m )
if ( f ~= /steve/mi )
if ( f !~ /steve/mi )
if ( f ~= /steve/miiiiiiiiiiiiiiiiimmmmmmmmmmmmmiiiii )`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "f"},
		{tokentype.CONTAINS, "~="},
		{tokentype.REGEXP, "(?i)steve"},
		{tokentype.RPAREN, ")"},
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "f"},
		{tokentype.CONTAINS, "~="},
		{tokentype.REGEXP, "(?m)steve"},
		{tokentype.RPAREN, ")"},
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "f"},
		{tokentype.CONTAINS, "~="},
		{tokentype.REGEXP, "(?mi)steve"},
		{tokentype.RPAREN, ")"},
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "f"},
		{tokentype.NOT_CONTAINS, "!~"},
		{tokentype.REGEXP, "(?mi)steve"},
		{tokentype.RPAREN, ")"},
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "f"},
		{tokentype.CONTAINS, "~="},
		{tokentype.REGEXP, "(?mi)steve"},
		{tokentype.RPAREN, ")"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestIllegalRegexp is designed to look for an unterminated/illegal regexp
func TestIllegalRegexp(t *testing.T) {
	input := `if ( f ~= /steve )`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.IF, "if"},
		{tokentype.LPAREN, "("},
		{tokentype.IDENT, "f"},
		{tokentype.CONTAINS, "~="},
		{tokentype.REGEXP, "unterminated regular expression"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestDiv is designed to test that a division is recognized; that it is
// not confused with a regular-expression.
func TestDiv(t *testing.T) {
	input := `a = b / c;
a = 3/4;
`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.IDENT, "a"},
		{tokentype.ASSIGN, "="},
		{tokentype.IDENT, "b"},
		{tokentype.SLASH, "/"},
		{tokentype.IDENT, "c"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.IDENT, "a"},
		{tokentype.ASSIGN, "="},
		{tokentype.INT, "3"},
		{tokentype.SLASH, "/"},
		{tokentype.INT, "4"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestDotDot is designed to ensure we get a ".." not an integer value.
func TestDotDot(t *testing.T) {
	input := `a = 1..10;`

	tests := []struct {
		expectedType    tokentype.TokenType
		expectedLiteral string
	}{
		{tokentype.IDENT, "a"},
		{tokentype.ASSIGN, "="},
		{tokentype.INT, "1"},
		{tokentype.DOTDOT, ".."},
		{tokentype.INT, "10"},
		{tokentype.SEMICOLON, ";"},
		{tokentype.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestGetLineStr test get source code line
func TestGetLineStr(t *testing.T) {
	code := `// A simple test-function for switch-statements.
	function test( name ) {
	
	  // Did we match?
	  m = false;
	
	  switch( name ) {
		case /^steve$/ , /^STEVE$/i {
		   printf("Hello Steve - we matched you via a regexp\n");
		   m = true;
		}
		case "St" + "even" {
		printf("Hello SteveN, you were matched via an expression\n" );
			m = true;
		}
		case 3, 6, 9 {
			printf("Hello multiple of three, we matched you literally: %d\n", int(name));
			m = true;
		}
		default {
		printf("Default case: %s\n", string(name) );
		}
	  }
	
	  // Show we matched, if we did.
	  if ( m ) { printf( "\tMatched!\n"); }
	}
	
	// Test the switch statement
	test( "Steve" );   // Regexp match
	test( "steve" );   // Regexp match
	test( "Steven" );  // Literal match
	test( 3 );         // Literal match
	
	// Unhandled/Default cases
	test( "Bob" );
	test( false );
	
	// Try some other numbers - only one will match
	foreach number in 1..10 {
	  test(number);
	}
	a = "으하하";
	printf( "All done\n" );`

	l := New(code)
	for ; l.ch != rune(0); l.readChar() {
	}
	for i := range l.codeLineBegins {
		println(l.GetLineStr(i))
	}
}

// Package lexer contains the code to lex input-programs into a stream
// of tokens, such that they may be parsed.
package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kasworld/nonkey/enum/tokentype"
	"github.com/kasworld/nonkey/interpreter/token"
)

// Lexer holds our object-state.
type Lexer struct {
	// for debug,error message
	curLine        int
	curPosInLine   int
	codeLineBegins []int // line begin pos

	// The current character position
	position int

	// The next character position
	readPosition int

	//The current character
	ch rune

	// A rune slice of our input string
	characters []rune

	// Previous token.
	prevToken token.Token
}

// New a Lexer instance from string input.
func New(input string) *Lexer {
	l := &Lexer{characters: []rune(input)}
	l.codeLineBegins = []int{0}
	l.readChar()
	return l
}

func (l *Lexer) GetLineStr(line int) string {
	lineBegin := l.codeLineBegins[line]
	if len(l.codeLineBegins) > line+1 {
		lineEnd := l.codeLineBegins[line+1]
		return string(l.characters[lineBegin:lineEnd])
	} else {
		return string(l.characters[lineBegin:])
	}
}

// CurrentLine return current line in source code
func (l *Lexer) CurrentLine() int {
	return l.curLine
}

// CurrentPosInLine return current pos in line in source code
func (l *Lexer) CurrentPosInLine() int {
	return l.curPosInLine
}

// read one forward character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.characters) {
		l.ch = rune(0)
		l.codeLineBegins = append(l.codeLineBegins, l.position+1)
	} else {
		l.ch = l.characters[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	// for debug, error message
	l.curPosInLine++
	if l.ch == rune('\n') {
		l.curLine++
		l.curPosInLine = 0
		l.codeLineBegins = append(l.codeLineBegins, l.position)
	}
}

// NextToken to read next token, skipping the white space.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	// skip single-line comments
	if l.ch == rune('#') ||
		(l.ch == rune('/') && l.peekChar() == rune('/')) {
		l.skipComment()
		return (l.NextToken())
	}

	// multi-line comments
	if l.ch == rune('/') && l.peekChar() == rune('*') {
		l.skipMultiLineComment()
	}

	switch l.ch {
	case rune('&'):
		if l.peekChar() == rune('&') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.AND, string(ch)+string(l.ch))
		}
	case rune('|'):
		if l.peekChar() == rune('|') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.OR, string(ch)+string(l.ch))
		}

	case rune('='):
		tok = l.newToken(tokentype.ASSIGN, string(l.ch))
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.ASSIGN, string(l.ch))
		}
	case rune(';'):
		tok = l.newToken(tokentype.SEMICOLON, string(l.ch))
	case rune('?'):
		tok = l.newToken(tokentype.QUESTION, string(l.ch))
	case rune('('):
		tok = l.newToken(tokentype.LPAREN, string(l.ch))
	case rune(')'):
		tok = l.newToken(tokentype.RPAREN, string(l.ch))
	case rune(','):
		tok = l.newToken(tokentype.COMMA, string(l.ch))
	case rune('.'):
		if l.peekChar() == rune('.') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.DOTDOT, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.PERIOD, string(l.ch))
		}
	case rune('+'):
		if l.peekChar() == rune('+') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.PLUS_PLUS, string(ch)+string(l.ch))
		} else if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.PLUS_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.PLUS, string(l.ch))
		}
	case rune('%'):
		tok = l.newToken(tokentype.MOD, string(l.ch))
	case rune('{'):
		tok = l.newToken(tokentype.LBRACE, string(l.ch))
	case rune('}'):
		tok = l.newToken(tokentype.RBRACE, string(l.ch))
	case rune('-'):
		if l.peekChar() == rune('-') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.MINUS_MINUS, string(ch)+string(l.ch))
		} else if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.MINUS_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.MINUS, string(l.ch))
		}
	case rune('/'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.SLASH_EQUALS, string(ch)+string(l.ch))
		} else {
			// slash is mostly division, but could
			// be the start of a regular expression

			// We exclude:
			//   a[b] / c       -> RBRACKET
			//   ( a + b ) / c   -> RPAREN
			//   a / c           -> IDENT
			//   3.2 / c         -> FLOAT
			//   1 / c           -> IDENT
			//
			if l.prevToken.Type == tokentype.RBRACKET ||
				l.prevToken.Type == tokentype.RPAREN ||
				l.prevToken.Type == tokentype.IDENT ||
				l.prevToken.Type == tokentype.INT ||
				l.prevToken.Type == tokentype.FLOAT {

				tok = l.newToken(tokentype.SLASH, string(l.ch))
			} else {
				str, err := l.readRegexp()
				if err == nil {
					tok.Type = tokentype.REGEXP
					tok.Literal = str
				} else {
					fmt.Printf("%s\n", err.Error())
					tok.Type = tokentype.REGEXP
					tok.Literal = str
				}
			}
		}
	case rune('*'):
		if l.peekChar() == rune('*') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.POW, string(ch)+string(l.ch))
		} else if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.ASTERISK_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.ASTERISK, string(l.ch))
		}
	case rune('<'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.LT_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.LT, string(l.ch))
		}
	case rune('>'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.GT_EQUALS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(tokentype.GT, string(l.ch))
		}
	case rune('~'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.CONTAINS, string(ch)+string(l.ch))
		}

	case rune('!'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = l.newToken(tokentype.NOT_EQ, string(ch)+string(l.ch))
		} else {
			if l.peekChar() == rune('~') {
				ch := l.ch
				l.readChar()
				tok = l.newToken(tokentype.NOT_CONTAINS, string(ch)+string(l.ch))

			} else {
				tok = l.newToken(tokentype.BANG, string(l.ch))
			}
		}
	case rune('"'):
		tok.Type = tokentype.STRING
		tok.Literal = l.readString()
	case rune('`'):
		tok.Type = tokentype.BACKTICK
		tok.Literal = l.readBacktick()
	case rune('['):
		tok = l.newToken(tokentype.LBRACKET, string(l.ch))
	case rune(']'):
		tok = l.newToken(tokentype.RBRACKET, string(l.ch))
	case rune(':'):
		tok = l.newToken(tokentype.COLON, string(l.ch))
	case rune(0):
		tok.Literal = ""
		tok.Type = tokentype.EOF
	default:

		if isDigit(l.ch) {
			tok = l.readDecimal()
			l.prevToken = tok
			return tok

		}
		tok.Literal = l.readIdentifier()
		tok.Type = tokentype.LookupKeyword(tok.Literal)
		l.prevToken = tok

		return tok
	}
	l.readChar()
	l.prevToken = tok
	return tok
}

// return new token
func (l *Lexer) newToken(tokenType tokentype.TokenType, s string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: s,
		Line:    l.curLine,
		Pos:     l.curPosInLine,
	}
}

// readIdentifier is designed to read an identifier (name of variable,
// function, etc).
//
// However there is a complication due to our historical implementation
// of the standard library.  We really want to stop identifiers if we hit
// a period, to allow method-calls to work on objects.
//
// So with input like this:
//
//   a.blah();
//
// Our identifier should be "a" (then we have a period, then a second
// identifier "blah", followed by opening & closing parenthesis).
//
// However we also have to cover the case of:
//
//    string.toupper( "blah" );
//    os.getenv( "PATH" );
//    ..
//
// So we have a horrid implementation..
func (l *Lexer) readIdentifier() string {

	//
	// Functions which are permitted to have dots in their name.
	//
	valid := map[string]bool{
		"directory.glob":     true,
		"math.abs":           true,
		"math.random":        true,
		"math.sqrt":          true,
		"os.environment":     true,
		"os.getenv":          true,
		"os.setenv":          true,
		"string.interpolate": true,
	}

	//
	// Types which will have valid methods.
	//
	types := []string{
		"string.",
		"array.",
		"integer.",
		"float.",
		"hash.",
		"object."}

	id := ""

	//
	// Save our position, in case we need to jump backwards in
	// our scanning.  Yeah.
	//
	position := l.position
	rposition := l.readPosition

	//
	// Build up our identifier, handling only valid characters.
	//
	// NOTE: This WILL consider the period valid, allowing the
	// parsing of "foo.bar", "os.getenv", "blah.blah.blah", etc.
	//
	for isIdentifier(l.ch) {
		id += string(l.ch)
		l.readChar()
	}

	//
	// Now we to see if our identifier had a period inside it.
	//
	if strings.Contains(id, ".") {

		// Is it a known-good function?
		ok := valid[id]

		// If not see if it has a type-prefix, which will
		// let the definition succeed.
		if !ok {
			for _, i := range types {
				if strings.HasPrefix(id, i) {
					ok = true
				}
			}
		}

		//
		// Not permitted?  Then we abort.
		//
		// We reset our lexer-state to the position just ahead
		// of the period.  This will then lead to a syntax
		// error.
		//
		// Which probably means our lexer should abort instead.
		//
		// For the moment we'll leave as-is.
		//
		if !ok {

			//
			// OK first of all we truncate our identifier
			// at the position before the "."
			//
			offset := strings.Index(id, ".")
			id = id[:offset]

			//
			// Now we have to move backwards - as a quickie
			// We'll reset our position and move forwards
			// the length of the bits we went too-far.
			l.position = position
			l.readPosition = rposition
			for offset > 0 {
				l.readChar()
				offset--
			}
		}
	}

	// And now our pain is over.
	return id
}

// skip white space
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// skip comment (until the end of the line).
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != rune(0) {
		l.readChar()
	}
	l.skipWhitespace()
}

// Consume all tokens until we've had the close of a multi-line
// comment.
func (l *Lexer) skipMultiLineComment() {
	found := false

	for !found {
		// break at the end of our input.
		if l.ch == rune(0) {
			found = true
		}

		// otherwise keep going until we find "*/".
		if l.ch == '*' && l.peekChar() == '/' {
			found = true

			// Our current position is "*", so skip
			// forward to consume the "/".
			l.readChar()
		}

		l.readChar()
	}

	l.skipWhitespace()
}

// read number - this handles 0x1234 and 0b101010101 too.
func (l *Lexer) readNumber() string {
	str := ""

	// We usually just accept digits.
	accept := "0123456789"

	// But if we have `0x` as a prefix we accept hexadecimal instead.
	if l.ch == '0' && l.peekChar() == 'x' {
		accept = "0x123456789abcdefABCDEF"
	}

	// If we have `0b` as a prefix we accept binary digits only.
	if l.ch == '0' && l.peekChar() == 'b' {
		accept = "b01"
	}

	for strings.Contains(accept, string(l.ch)) {
		str += string(l.ch)
		l.readChar()
	}
	return str
}

// read decimal
func (l *Lexer) readDecimal() token.Token {

	//
	// Read an integer-number.
	//
	integer := l.readNumber()

	//
	// Now we either expect:
	//
	//   .[digits]  -> Which converts us from an int to a float.
	//
	//   .blah      -> Which is a method-call on a raw number.
	//
	if l.ch == rune('.') && isDigit(l.peekChar()) {
		//
		// OK here we think we've got a float.
		//
		l.readChar()
		fraction := l.readNumber()
		return l.newToken(tokentype.FLOAT, integer+"."+fraction)
	}
	return l.newToken(tokentype.INT, integer)
}

// read string
func (l *Lexer) readString() string {
	out := ""

	for {
		l.readChar()
		if l.ch == '"' {
			break
		}

		//
		// Handle \n, \r, \t, \", etc.
		//
		if l.ch == '\\' {
			l.readChar()

			if l.ch == rune('n') {
				l.ch = '\n'
			}
			if l.ch == rune('r') {
				l.ch = '\r'
			}
			if l.ch == rune('t') {
				l.ch = '\t'
			}
			if l.ch == rune('"') {
				l.ch = '"'
			}
			if l.ch == rune('\\') {
				l.ch = '\\'
			}
		}
		out = out + string(l.ch)
	}

	return out
}

// read a regexp, including flags.
func (l *Lexer) readRegexp() (string, error) {
	out := ""

	for {
		l.readChar()

		if l.ch == rune(0) {
			return "unterminated regular expression", fmt.Errorf("unterminated regular expression")
		}
		if l.ch == '/' {

			// consume the terminating "/".
			l.readChar()

			// prepare to look for flags
			flags := ""

			// two flags are supported:
			//   i -> Ignore-case
			//   m -> Multiline
			//
			for l.ch == rune('i') || l.ch == rune('m') {

				// save the char - unless it is a repeat
				if !strings.Contains(flags, string(l.ch)) {

					// we're going to sort the flags
					tmp := strings.Split(flags, "")
					tmp = append(tmp, string(l.ch))
					flags = strings.Join(tmp, "")

				}

				// read the next
				l.readChar()
			}

			// convert the regexp to go-lang
			if len(flags) > 0 {
				out = "(?" + flags + ")" + out
			}
			break
		}
		out = out + string(l.ch)
	}

	return out, nil
}

// read the end of a backtick-quoted string
func (l *Lexer) readBacktick() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '`' {
			break
		}
	}
	out := string(l.characters[position:l.position])
	return out
}

// peek character
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.readPosition]
}

// determinate ch is identifier or not
func isIdentifier(ch rune) bool {

	if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '.' || ch == '?' || ch == '$' || ch == '_' {
		return true
	}

	return false
}

// is white space
func isWhitespace(ch rune) bool {
	return ch == rune(' ') || ch == rune('\t') || ch == rune('\n') || ch == rune('\r')
}

// is Digit
func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}

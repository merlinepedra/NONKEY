// Package parser is used to parse input-programs written in monkey
// and convert them to an abstract-syntax tree.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kasworld/nonkey/enum/precedence"
	"github.com/kasworld/nonkey/enum/tokentype"
	"github.com/kasworld/nonkey/interpreter/ast"
	"github.com/kasworld/nonkey/interpreter/lexer"
	"github.com/kasworld/nonkey/interpreter/token"
)

// prefix Parse function
// infix parse function
// postfix parse function
type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func() ast.Expression
)

// Parser object
type Parser struct {
	// l is our lexer
	l *lexer.Lexer

	// prevToken holds the previous token from our lexer.
	// (used for "++" + "--")
	prevToken token.Token

	// curToken holds the current token from our lexer.
	curToken token.Token

	// peekToken holds the next token which will come from the lexer.
	peekToken token.Token

	// errors holds parsing-errors.
	errors []string

	// prefixParseFns holds a map of parsing methods for
	// prefix-based syntax.
	prefixParseFns [tokentype.TokenType_Count]prefixParseFn

	// infixParseFns holds a map of parsing methods for
	// infix-based syntax.
	infixParseFns [tokentype.TokenType_Count]infixParseFn

	// postfixParseFns holds a map of parsing methods for
	// postfix-based syntax.
	postfixParseFns [tokentype.TokenType_Count]postfixParseFn

	// are we inside a ternary expression?
	//
	// Nested ternary expressions are illegal :)
	tern bool
}

// New returns our new parser-object.
func New(l *lexer.Lexer) *Parser {

	// Create the parser, and prime the pump
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	// Register prefix-functions
	p.prefixParseFns = [tokentype.TokenType_Count]prefixParseFn{
		tokentype.BACKTICK:        p.parseBacktickLiteral,
		tokentype.BANG:            p.parsePrefixExpression,
		tokentype.DEFINE_FUNCTION: p.parseFunctionDefinition,
		tokentype.EOF:             p.parsingBroken,
		tokentype.FALSE:           p.parseBoolean,
		tokentype.FLOAT:           p.parseFloatLiteral,
		tokentype.FOR:             p.parseForLoopExpression,
		tokentype.FOREACH:         p.parseForEach,
		tokentype.FUNCTION:        p.parseFunctionLiteral,
		tokentype.IDENT:           p.parseIdentifier,
		tokentype.IF:              p.parseIfExpression,
		tokentype.ILLEGAL:         p.parsingBroken,
		tokentype.INT:             p.parseIntegerLiteral,
		tokentype.LBRACE:          p.parseHashLiteral,
		tokentype.LBRACKET:        p.parseArrayLiteral,
		tokentype.LPAREN:          p.parseGroupedExpression,
		tokentype.MINUS:           p.parsePrefixExpression,
		tokentype.REGEXP:          p.parseRegexpLiteral,
		tokentype.STRING:          p.parseStringLiteral,
		tokentype.SWITCH:          p.parseSwitchStatement,
		tokentype.TRUE:            p.parseBoolean,
	}

	// Register infix functions
	p.infixParseFns = [tokentype.TokenType_Count]infixParseFn{
		tokentype.AND:             p.parseInfixExpression,
		tokentype.ASSIGN:          p.parseAssignExpression,
		tokentype.ASTERISK:        p.parseInfixExpression,
		tokentype.ASTERISK_EQUALS: p.parseAssignExpression,
		tokentype.CONTAINS:        p.parseInfixExpression,
		tokentype.DOTDOT:          p.parseInfixExpression,
		tokentype.EQ:              p.parseInfixExpression,
		tokentype.GT:              p.parseInfixExpression,
		tokentype.GT_EQUALS:       p.parseInfixExpression,
		tokentype.LBRACKET:        p.parseIndexExpression,
		tokentype.LPAREN:          p.parseCallExpression,
		tokentype.LT:              p.parseInfixExpression,
		tokentype.LT_EQUALS:       p.parseInfixExpression,
		tokentype.MINUS:           p.parseInfixExpression,
		tokentype.MINUS_EQUALS:    p.parseAssignExpression,
		tokentype.MOD:             p.parseInfixExpression,
		tokentype.NOT_CONTAINS:    p.parseInfixExpression,
		tokentype.NOT_EQ:          p.parseInfixExpression,
		tokentype.OR:              p.parseInfixExpression,
		tokentype.PERIOD:          p.parseMethodCallExpression,
		tokentype.PLUS:            p.parseInfixExpression,
		tokentype.PLUS_EQUALS:     p.parseAssignExpression,
		tokentype.POW:             p.parseInfixExpression,
		tokentype.QUESTION:        p.parseTernaryExpression,
		tokentype.SLASH:           p.parseInfixExpression,
		tokentype.SLASH_EQUALS:    p.parseAssignExpression,
	}

	// Register postfix functions.
	p.postfixParseFns = [tokentype.TokenType_Count]postfixParseFn{
		tokentype.MINUS_MINUS: p.parsePostfixExpression,
		tokentype.PLUS_PLUS:   p.parsePostfixExpression,
	}

	// All done
	return p
}

// Errors return stored errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError raises an error if the next token is not the expected type.
func (p *Parser) peekError(t tokentype.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead around line %d", t, p.curToken.Type, p.l.GetLine())
	p.errors = append(p.errors, msg)
}

// nextToken moves to our next token from the lexer.
func (p *Parser) nextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram used to parse the whole program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != tokentype.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement parses a single statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tokentype.LET:
		return p.parseLetStatement()
	case tokentype.CONST:
		return p.parseConstStatement()
	case tokentype.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses a let-statement.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(tokentype.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(tokentype.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(precedence.LOWEST)
	for !p.curTokenIs(tokentype.SEMICOLON) {

		if p.curTokenIs(tokentype.EOF) {
			p.errors = append(p.errors, "unterminated let statement")
			return nil
		}

		p.nextToken()
	}
	return stmt
}

// parseConstStatement parses a constant declaration.
func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}
	if !p.expectPeek(tokentype.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(tokentype.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(precedence.LOWEST)
	for !p.curTokenIs(tokentype.SEMICOLON) {

		if p.curTokenIs(tokentype.EOF) {
			p.errors = append(p.errors, "unterminated const statement")
			return nil
		}

		p.nextToken()
	}
	return stmt
}

// parseReturnStatement parses a return-statement.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(precedence.LOWEST)
	for !p.curTokenIs(tokentype.SEMICOLON) {

		if p.curTokenIs(tokentype.EOF) {
			p.errors = append(p.errors, "unterminated return statement")
			return nil
		}

		p.nextToken()
	}
	return stmt
}

// no prefix parse function error
func (p *Parser) noPrefixParseFnError(t tokentype.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found around line %d", t, p.l.GetLine())
	p.errors = append(p.errors, msg)
}

// parse Expression Statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(precedence.LOWEST)
	for p.peekTokenIs(tokentype.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence1 precedence.Precedence) ast.Expression {
	postfix := p.postfixParseFns[p.curToken.Type]
	if postfix != nil {
		return (postfix())
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(tokentype.SEMICOLON) && precedence1 < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// parsingBroken is hit if we see an EOF in our input-stream
// this means we're screwed
func (p *Parser) parsingBroken() ast.Expression {
	return nil
}

// parseIdentifier parses an identifier.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral parses an integer literal.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	var value int64
	var err error

	if strings.HasPrefix(p.curToken.Literal, "0b") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 2, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0x") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 16, 64)
	} else {
		value, err = strconv.ParseInt(p.curToken.Literal, 10, 64)
	}

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer around line %d", p.curToken.Literal, p.l.GetLine())
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// parseFloatLiteral parses a float-literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	flo := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float around line %d", p.curToken.Literal, p.l.GetLine())
		p.errors = append(p.errors, msg)
		return nil
	}
	flo.Value = value
	return flo
}

// parseSwitchStatement handles a switch statement
func (p *Parser) parseSwitchStatement() ast.Expression {

	// switch
	expression := &ast.SwitchExpression{Token: p.curToken}

	// look for (xx)
	if !p.expectPeek(tokentype.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Value = p.parseExpression(precedence.LOWEST)
	if expression.Value == nil {
		return nil
	}
	if !p.expectPeek(tokentype.RPAREN) {
		return nil
	}

	// Now we have a block containing blocks.
	if !p.expectPeek(tokentype.LBRACE) {
		return nil
	}
	p.nextToken()

	// Process the block which we think will contain
	// various case-statements
	for !p.curTokenIs(tokentype.RBRACE) {

		if p.curTokenIs(tokentype.EOF) {
			p.errors = append(p.errors, "unterminated switch statement")
			return nil
		}
		tmp := &ast.CaseExpression{Token: p.curToken}

		// Default will be handled specially
		if p.curTokenIs(tokentype.DEFAULT) {

			// We have a default-case here.
			tmp.Default = true

		} else if p.curTokenIs(tokentype.CASE) {

			// skip "case"
			p.nextToken()

			// Here we allow "case default" even though
			// most people would prefer to write "default".
			if p.curTokenIs(tokentype.DEFAULT) {
				tmp.Default = true
			} else {

				// parse the match-expression.
				tmp.Expr = append(tmp.Expr, p.parseExpression(precedence.LOWEST))
				for p.peekTokenIs(tokentype.COMMA) {

					// skip the comma
					p.nextToken()

					// setup the expression.
					p.nextToken()

					tmp.Expr = append(tmp.Expr, p.parseExpression(precedence.LOWEST))

				}
			}
		}

		if !p.expectPeek(tokentype.LBRACE) {

			msg := fmt.Sprintf("expected token to be '{', got %s instead", p.curToken.Type)
			p.errors = append(p.errors, msg)
			fmt.Printf("error\n")
			return nil
		}

		// parse the block
		tmp.Block = p.parseBlockStatement()

		if !p.curTokenIs(tokentype.RBRACE) {
			msg := fmt.Sprintf("Syntax Error: expected token to be '}', got %s instead", p.curToken.Type)
			p.errors = append(p.errors, msg)
			fmt.Printf("error\n")
			return nil

		}
		p.nextToken()

		// save the choice away
		expression.Choices = append(expression.Choices, tmp)

	}

	// ensure we're at the the closing "}"
	if !p.curTokenIs(tokentype.RBRACE) {
		return nil
	}

	// More than one default is a bug
	count := 0
	for _, c := range expression.Choices {
		if c.Default {
			count++
		}
	}
	if count > 1 {
		msg := fmt.Sprintf("A switch-statement should only have one default block")
		p.errors = append(p.errors, msg)
		return nil

	}
	return expression

}

// parseBoolean parses a boolean token.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(tokentype.TRUE)}
}

// parsePrefixExpression parses a prefix-based expression.
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(precedence.PREFIX)
	return expression
}

// parsePostfixExpression parses a postfix-based expression.
func (p *Parser) parsePostfixExpression() ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.prevToken,
		Operator: p.curToken.Literal,
	}
	return expression
}

// parseInfixExpression parses an infix-based expression.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	curPrecedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(curPrecedence)
	return expression
}

// parseTernaryExpression parses a ternary expression
func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {

	if p.tern {
		msg := fmt.Sprintf("nested ternary expressions are illegal, around line %d", p.l.GetLine())
		p.errors = append(p.errors, msg)
		return nil
	}

	p.tern = true
	defer func() { p.tern = false }()

	expression := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}
	p.nextToken() //skip the '?'
	curPrecedence := p.curPrecedence()
	expression.IfTrue = p.parseExpression(curPrecedence)

	if !p.expectPeek(tokentype.COLON) { //skip the ":"
		return nil
	}

	// Get to next token, then parse the else part
	p.nextToken()
	expression.IfFalse = p.parseExpression(curPrecedence)

	p.tern = false
	return expression
}

// parseGroupedExpression parses a grouped-expression.
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(precedence.LOWEST)
	if !p.expectPeek(tokentype.RPAREN) {
		return nil
	}
	return exp
}

// parseIfCondition parses an if-expression.
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(tokentype.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(precedence.LOWEST)
	if !p.expectPeek(tokentype.RPAREN) {
		return nil
	}
	if !p.expectPeek(tokentype.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(tokentype.ELSE) {
		p.nextToken()
		if !p.expectPeek(tokentype.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

// parseForLoopExpression parses a for-loop.
func (p *Parser) parseForLoopExpression() ast.Expression {
	expression := &ast.ForLoopExpression{Token: p.curToken}
	if !p.expectPeek(tokentype.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(precedence.LOWEST)
	if !p.expectPeek(tokentype.RPAREN) {
		return nil
	}
	if !p.expectPeek(tokentype.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	return expression
}

// parseForEach parses 'foreach x X { .. block .. }`
func (p *Parser) parseForEach() ast.Expression {
	expression := &ast.ForeachStatement{Token: p.curToken}

	// get the id
	p.nextToken()
	expression.Ident = p.curToken.Literal

	// If we find a "," we then get a second identifier too.
	if p.peekTokenIs(tokentype.COMMA) {

		//
		// Generally we have:
		//
		//    foreach IDENT in THING { .. }
		//
		// If we have two arguments the first becomes
		// the index, and the second becomes the IDENT.
		//

		// skip the comma
		p.nextToken()

		if !p.peekTokenIs(tokentype.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("second argument to foreach must be ident, got %v", p.peekToken))
			return nil
		}
		p.nextToken()

		//
		// Record the updated values.
		//
		expression.Index = expression.Ident
		expression.Ident = p.curToken.Literal

	}

	// The next token, after the ident(s), should be `in`.
	if !p.expectPeek(tokentype.IN) {
		return nil
	}
	p.nextToken()

	// get the thing we're going to iterate  over.
	expression.Value = p.parseExpression(precedence.LOWEST)
	if expression.Value == nil {
		return nil
	}

	// parse the block
	p.nextToken()
	expression.Body = p.parseBlockStatement()

	return expression
}

// parseBlockStatement parsea a block.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(tokentype.RBRACE) {

		// Don't loop forever
		if p.curTokenIs(tokentype.EOF) {
			p.errors = append(p.errors,
				"unterminated block statement")
			return nil
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// parseFunctionLiteral parses a function-literal.
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(tokentype.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(tokentype.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionDefinition parses the definition of a function.
func (p *Parser) parseFunctionDefinition() ast.Expression {
	p.nextToken()
	lit := &ast.FunctionDefineLiteral{Token: p.curToken}
	if !p.expectPeek(tokentype.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(tokentype.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionParameters parses the parameters used for a function.
func (p *Parser) parseFunctionParameters() (map[string]ast.Expression, []*ast.Identifier) {

	// Any default parameters.
	m := make(map[string]ast.Expression)

	// The argument-definitions.
	identifiers := make([]*ast.Identifier, 0)

	// Is the next parameter ")" ?  If so we're done. No args.
	if p.peekTokenIs(tokentype.RPAREN) {
		p.nextToken()
		return m, identifiers
	}
	p.nextToken()

	// Keep going until we find a ")"
	for !p.curTokenIs(tokentype.RPAREN) {

		if p.curTokenIs(tokentype.EOF) {
			p.errors = append(p.errors, "unterminated function parameters")
			return nil, nil
		}

		// Get the identifier.
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
		p.nextToken()

		// If there is "=xx" after the name then that's
		// the default parameter.
		if p.curTokenIs(tokentype.ASSIGN) {
			p.nextToken()
			// Save the default value.
			m[ident.Value] = p.parseExpressionStatement().Expression
			p.nextToken()
		}

		// Skip any comma.
		if p.curTokenIs(tokentype.COMMA) {
			p.nextToken()
		}
	}

	return m, identifiers
}

// parseStringLiteral parses a string-literal.
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseRegexpLiteral parses a regular-expression.
func (p *Parser) parseRegexpLiteral() ast.Expression {

	flags := ""

	val := p.curToken.Literal
	if strings.HasPrefix(val, "(?") {
		val = strings.TrimPrefix(val, "(?")

		i := 0
		for i < len(val) {

			if val[i] == ')' {

				val = val[i+1:]
				break
			} else {
				flags += string(val[i])
			}

			i++
		}
	}
	return &ast.RegexpLiteral{Token: p.curToken, Value: val, Flags: flags}
}

// parseBacktickLiteral parses a backtick-expression.
func (p *Parser) parseBacktickLiteral() ast.Expression {
	return &ast.BacktickLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseArrayLiteral parses an array literal.
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(tokentype.RBRACKET)
	return array
}

// parsearray elements literal
func (p *Parser) parseExpressionList(end tokentype.TokenType) []ast.Expression {
	list := make([]ast.Expression, 0)
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(precedence.LOWEST))
	for p.peekTokenIs(tokentype.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(precedence.LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// parseInfixExpression parsea an array index expression.
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(precedence.LOWEST)
	if !p.expectPeek(tokentype.RBRACKET) {
		return nil
	}
	return exp
}

// parseAssignExpression parses a bare assignment, without a `let`.
func (p *Parser) parseAssignExpression(name ast.Expression) ast.Expression {
	stmt := &ast.AssignStatement{Token: p.curToken}
	if n, ok := name.(*ast.Identifier); ok {
		stmt.Name = n
	} else {
		msg := fmt.Sprintf("expected assign token to be IDENT, got %s instead around line %d", name.TokenLiteral(), p.l.GetLine())
		p.errors = append(p.errors, msg)
	}

	oper := p.curToken
	p.nextToken()

	//
	// An assignment is generally:
	//
	//    variable = value
	//
	// But we cheat and reuse the implementation for:
	//
	//    i += 4
	//
	// In this case we record the "operator" as "+="
	//
	switch oper.Type {
	case tokentype.PLUS_EQUALS:
		stmt.Operator = "+="
	case tokentype.MINUS_EQUALS:
		stmt.Operator = "-="
	case tokentype.SLASH_EQUALS:
		stmt.Operator = "/="
	case tokentype.ASTERISK_EQUALS:
		stmt.Operator = "*="
	default:
		stmt.Operator = "="
	}
	stmt.Value = p.parseExpression(precedence.LOWEST)
	return stmt
}

// parseCallExpression parses a function-call expression.
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(tokentype.RPAREN)
	return exp
}

// parseHashLiteral parses a hash literal.
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	for !p.peekTokenIs(tokentype.RBRACE) {
		p.nextToken()
		key := p.parseExpression(precedence.LOWEST)
		if !p.expectPeek(tokentype.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(precedence.LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(tokentype.RBRACE) && !p.expectPeek(tokentype.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(tokentype.RBRACE) {
		return nil
	}
	return hash
}

// parseMethodCallExpression parses an object-based method-call.
func (p *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	methodCall := &ast.ObjectCallExpression{Token: p.curToken, Object: obj}
	p.nextToken()
	name := p.parseIdentifier()
	p.nextToken()
	methodCall.Call = p.parseCallExpression(name)
	return methodCall
}

// curTokenIs tests if the current token has the given type.
func (p *Parser) curTokenIs(t tokentype.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs tests if the next token has the given type.
func (p *Parser) peekTokenIs(t tokentype.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek validates the next token is of the given type,
// and advances if so.  If it is not an error is stored.
func (p *Parser) expectPeek(t tokentype.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

// peekPrecedence looks up the next token precedence.
func (p *Parser) peekPrecedence() precedence.Precedence {
	if prd, ok := tokentype.Token2Precedences[p.peekToken.Type]; ok {
		return prd
	}
	return precedence.LOWEST
}

// curPrecedence looks up the current token precedence.
func (p *Parser) curPrecedence() precedence.Precedence {
	if prd, ok := tokentype.Token2Precedences[p.curToken.Type]; ok {
		return prd
	}
	return precedence.LOWEST
}

// Package ast contains the definitions of the abstract-syntax tree
// that our parse produces, and our interpreter executes.
package ast

import (
	"bytes"

	"github.com/kasworld/nonkey/interpreter/asti"
	"github.com/kasworld/nonkey/interpreter/token"
)

// LetStatement holds a let-statemnt
type LetStatement struct {
	// Token holds the token
	Token token.Token

	// Name is the name of the variable to which we're assigning
	Name *Identifier

	// Value is the thing we're storing in the variable.
	Value asti.ExpressionI
}

func (ls *LetStatement) StatementNode() {}

// TokenLiteral returns the literal token.
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// String returns this object as a string.
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.TokenLiteral())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ConstStatement is the same as let-statement, but the value
// can't be changed later.
type ConstStatement struct {
	// Token is the token
	Token token.Token

	// Name is the name of the variable we're setting
	Name *Identifier

	// Value contains the value which is to be set
	Value asti.ExpressionI
}

func (ls *ConstStatement) StatementNode() {}

// TokenLiteral returns the literal token.
func (ls *ConstStatement) TokenLiteral() string { return ls.Token.Literal }

// String returns this object as a string.
func (ls *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.TokenLiteral())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement stores a return-statement
type ReturnStatement struct {
	// Token contains the literal token.
	Token token.Token

	// ReturnValue is the value whichis to be returned.
	ReturnValue asti.ExpressionI
}

func (rs *ReturnStatement) StatementNode() {}

// TokenLiteral returns the literal token.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns this object as a string.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.TokenLiteral())
	}
	out.WriteString(";")
	return out.String()
}

// BlockStatement holds a group of statements, which are treated
// as a block.  (For example the body of an `if` expression.)
type BlockStatement struct {
	// Token holds the actual token
	Token token.Token

	// Statements contain the set of statements within the block
	Statements []asti.StatementI
}

func (bs *BlockStatement) StatementNode() {}

// TokenLiteral returns the literal token.
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String returns this object as a string.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ForeachStatement holds a foreach-statement.
type ForeachStatement struct {
	// Token is the actual token
	Token token.Token

	// Index is the variable we'll set with the index, for the blocks' scope
	//
	// This is optional.
	Index string

	// Ident is the variable we'll set with each item, for the blocks' scope
	Ident string

	// Value is the thing we'll range over.
	Value asti.ExpressionI

	// Body is the block we'll execute.
	Body *BlockStatement
}

func (fes *ForeachStatement) ExpressionNode() {}

// TokenLiteral returns the literal token.
func (fes *ForeachStatement) TokenLiteral() string { return fes.Token.Literal }

// String returns this object as a string.
func (fes *ForeachStatement) String() string {
	var out bytes.Buffer
	out.WriteString("foreach ")
	out.WriteString(fes.Ident)
	out.WriteString(" ")
	out.WriteString(fes.Value.String())
	out.WriteString(fes.Body.String())
	return out.String()
}

// AssignStatement is generally used for a (let-less) assignment,
// such as "x = y", however we allow an operator to be stored ("=" in that
// example), such that we can do self-operations.
//
// Specifically "x += y" is defined as an assignment-statement with
// the operator set to "+=".  The same applies for "+=", "-=", "*=", and
// "/=".
type AssignStatement struct {
	Token    token.Token
	Name     *Identifier
	Operator string
	Value    asti.ExpressionI
}

func (as *AssignStatement) ExpressionNode() {}

// TokenLiteral returns the literal token.
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }

// String returns this object as a string.
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(as.Operator)
	out.WriteString(as.Value.String())
	return out.String()
}

// Package ast contains the definitions of the abstract-syntax tree
// that our parse produces, and our interpreter executes.
package ast

import (
	"bytes"
	"fmt"

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

// GetToken returns the token.
func (ls *LetStatement) GetToken() token.Token { return ls.Token }

// String returns this object as a string.
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%v %v = %v;",
		ls.GetToken().Literal,
		ls.Name.GetToken().Literal,
		ls.Value,
	)
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

// GetToken returns the token.
func (ls *ConstStatement) GetToken() token.Token { return ls.Token }

// String returns this object as a string.
func (ls *ConstStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%v %v = %v;",
		ls.GetToken().Literal,
		ls.Name.GetToken().Literal,
		ls.Value,
	)
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

// GetToken returns the token.
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }

// String returns this object as a string.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%v %v;",
		rs.GetToken(),
		rs.ReturnValue.GetToken(),
	)
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

// GetToken returns the token.
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }

// String returns this object as a string.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (es *ExpressionStatement) StatementNode() {}

// GetToken returns the token.
func (es *ExpressionStatement) GetToken() token.Token { return es.Token }

// String returns this object as a string.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

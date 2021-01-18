// Package ast contains the definitions of the abstract-syntax tree
// that our parse produces, and our interpreter executes.
package ast

import (
	"bytes"

	"github.com/kasworld/nonkey/interpreter/asti"
)

// Program represents a complete program.
type Program struct {
	// Statements is the set of statements which the program is comprised
	// of.
	Statements []asti.StatementI
}

// TokenLiteral returns the literal token of our program.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns this object as a string.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

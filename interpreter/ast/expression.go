// Package ast contains the definitions of the abstract-syntax tree
// that our parse produces, and our interpreter executes.
package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kasworld/nonkey/enum/tokentype"
	"github.com/kasworld/nonkey/interpreter/asti"
	"github.com/kasworld/nonkey/interpreter/token"
)

// Identifier holds a single identifier.
type Identifier struct {
	// Token is the literal token
	Token token.Token

	// Value is the name of the identifier
	Value string
}

func (i *Identifier) ExpressionNode() {}

// GetToken returns the token.
func (i *Identifier) GetToken() token.Token { return i.Token }

// String returns this object as a string.
func (i *Identifier) String() string {
	return i.Value
}

// ExpressionStatement is an expression
type ExpressionStatement struct {
	// Token is the literal token
	Token token.Token

	// Expression holds the expression
	Expression asti.ExpressionI
}

// IntegerLiteral holds an integer
type IntegerLiteral struct {
	// Token is the literal token
	Token token.Token

	// Value holds the integer.
	Value int64
}

func (il *IntegerLiteral) ExpressionNode() {}

// GetToken returns the token.
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }

// String returns this object as a string.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// FloatLiteral holds a floating-point number
type FloatLiteral struct {
	// Token is the literal token
	Token token.Token

	// Value holds the floating-point number.
	Value float64
}

func (fl *FloatLiteral) ExpressionNode() {}

// GetToken returns the token.
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }

// String returns this object as a string.
func (fl *FloatLiteral) String() string { return fl.Token.Literal }

// PrefixExpression holds a prefix-based expression
type PrefixExpression struct {
	// Token holds the token.  e.g. "!"
	Token token.Token

	// Operator holds the operator being invoked (e.g. "!" ).
	Operator tokentype.TokenType

	// Right holds the thing to be operated upon
	Right asti.ExpressionI
}

func (pe *PrefixExpression) ExpressionNode() {}

// GetToken returns the token.
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }

// String returns this object as a string.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "(%v%v)", tokentype.Attrib[pe.Operator].String, pe.Right)
	return out.String()
}

// InfixExpression stores an infix expression.
type InfixExpression struct {
	// Token holds the literal expression
	Token token.Token

	// Left holds the left-most argument
	Left asti.ExpressionI

	// Operator holds the operation to be carried out (e.g. "+", "-" )
	Operator tokentype.TokenType

	// Right holds the right-most argument
	Right asti.ExpressionI
}

func (ie *InfixExpression) ExpressionNode() {}

// GetToken returns the token.
func (ie *InfixExpression) GetToken() token.Token { return ie.Token }

// String returns this object as a string.
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "(%v %v %v)", ie.Left, tokentype.Attrib[ie.Operator].String, ie.Right)
	return out.String()
}

// PostfixExpression holds a postfix-based expression
type PostfixExpression struct {
	// Token holds the token we're operating upon
	Token token.Token
	// Operator holds the postfix token, e.g. ++
	Operator tokentype.TokenType
}

func (pe *PostfixExpression) ExpressionNode() {}

// GetToken returns the token.
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }

// String returns this object as a string.
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "(%v%v)", pe.Token.Literal, tokentype.Attrib[pe.Operator].String)
	return out.String()
}

// Boolean holds a boolean type
type Boolean struct {
	// Token holds the actual token
	Token token.Token

	// Value stores the bools' value: true, or false.
	Value bool
}

func (b *Boolean) ExpressionNode() {}

// GetToken returns the token.
func (b *Boolean) GetToken() token.Token { return b.Token }

// String returns this object as a string.
func (b *Boolean) String() string { return b.Token.Literal }

// IfExpression holds an if-statement
type IfExpression struct {
	// Token is the actual token
	Token token.Token

	// Condition is the thing that is evaluated to determine
	// which block should be executed.
	Condition asti.ExpressionI

	// Consequence is the set of statements executed if the
	// condition is true.
	Consequence *BlockStatement

	// Alternative is the set of statements executed if the
	// condition is not true (optional).
	Alternative *BlockStatement
}

func (ie *IfExpression) ExpressionNode() {}

// GetToken returns the token.
func (ie *IfExpression) GetToken() token.Token { return ie.Token }

// String returns this object as a string.
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "if %v %v", ie.Condition, ie.Consequence)
	if ie.Alternative != nil {
		fmt.Fprintf(&out, "else %v", ie.Alternative)
	}
	return out.String()
}

// TernaryExpression holds a ternary-expression.
type TernaryExpression struct {
	// Token is the actual token.
	Token token.Token

	// Condition is the thing that is evaluated to determine
	// which expression should be returned
	Condition asti.ExpressionI

	// IfTrue is the expression to return if the condition is true.
	IfTrue asti.ExpressionI

	// IFFalse is the expression to return if the condition is not true.
	IfFalse asti.ExpressionI
}

func (te *TernaryExpression) ExpressionNode() {}

// GetToken returns the token.
func (te *TernaryExpression) GetToken() token.Token { return te.Token }

// String returns this object as a string.
func (te *TernaryExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "(%v ? %v : %v)", te.Condition, te.IfTrue, te.IfFalse)
	return out.String()
}

// ForLoopExpression holds a for-loop
type ForLoopExpression struct {
	// Token is the actual token
	Token token.Token

	// Condition is the expression used to determine if the loop
	// is still running.
	Condition asti.ExpressionI

	// Consequence is the set of statements to be executed for the
	// loop body.
	Consequence *BlockStatement
}

func (fle *ForLoopExpression) ExpressionNode() {}

// GetToken returns the token.
func (fle *ForLoopExpression) GetToken() token.Token { return fle.Token }

// String returns this object as a string.
func (fle *ForLoopExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "for (%v) {%v}", fle.Condition, fle.Consequence)
	return out.String()
}

// FunctionLiteral holds a function-definition
//
// See-also FunctionDefineLiteral.
type FunctionLiteral struct {
	// Token is the actual token
	Token token.Token

	// Parameters is the list of parameters the function receives.
	Parameters []*Identifier

	// Defaults holds any default values for arguments which aren't
	// specified
	Defaults map[string]asti.ExpressionI

	// Body contains the set of statements within the function.
	Body *BlockStatement
}

func (fl *FunctionLiteral) ExpressionNode() {}

// GetToken returns the token.
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }

// String returns this object as a string.
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	fmt.Fprintf(&out, "%v(%v) %v", fl.GetToken().Literal, strings.Join(params, ", "), fl.Body)
	return out.String()

}

// FunctionDefineLiteral holds a function-definition.
//
// See-also FunctionLiteral.
type FunctionDefineLiteral struct {
	// Token holds the token
	Token token.Token

	// Paremeters holds the function parameters.
	Parameters []*Identifier

	// Defaults holds any default-arguments.
	Defaults map[string]asti.ExpressionI

	// Body holds the set of statements in the functions' body.
	Body *BlockStatement
}

func (fl *FunctionDefineLiteral) ExpressionNode() {}

// GetToken returns the token.
func (fl *FunctionDefineLiteral) GetToken() token.Token { return fl.Token }

// String returns this object as a string.
func (fl *FunctionDefineLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	fmt.Fprintf(&out, "%v(%v) %v", fl.GetToken().Literal, strings.Join(params, ", "), fl.Body)
	return out.String()

}

// CallExpression holds the invokation of a method-call.
type CallExpression struct {
	// Token stores the literal token
	Token token.Token

	// Function is the function to be invoked.
	Function asti.ExpressionI

	// Arguments are the arguments to be applied
	Arguments []asti.ExpressionI
}

func (ce *CallExpression) ExpressionNode() {}

// GetToken returns the token.
func (ce *CallExpression) GetToken() token.Token { return ce.Token }

// String returns this object as a string.
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	fmt.Fprintf(&out, "%v(%v)", ce.Function, strings.Join(args, ", "))
	return out.String()
}

// ObjectCallExpression is used when calling a method on an object.
type ObjectCallExpression struct {
	// Token is the literal token
	Token token.Token

	// Object is the object against which the call is invoked.
	Object asti.ExpressionI

	// Call is the method-name.
	Call asti.ExpressionI
}

func (oce *ObjectCallExpression) ExpressionNode() {}

// GetToken returns the token.
func (oce *ObjectCallExpression) GetToken() token.Token { return oce.Token }

// String returns this object as a string.
func (oce *ObjectCallExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%v.%v", oce.Object, oce.Call)
	return out.String()
}

// StringLiteral holds a string
type StringLiteral struct {
	// Token is the token
	Token token.Token

	// Value is the value of the string.
	Value string
}

func (sl *StringLiteral) ExpressionNode() {}

// GetToken returns the token.
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }

// String returns this object as a string.
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// RegexpLiteral holds a regular-expression.
type RegexpLiteral struct {
	// Token is the token
	Token token.Token

	// Value is the value of the regular expression.
	Value string

	// Flags contains any flags associated with the regexp.
	Flags string
}

func (rl *RegexpLiteral) ExpressionNode() {}

// GetToken returns the token.
func (rl *RegexpLiteral) GetToken() token.Token { return rl.Token }

// String returns this object as a string.
func (rl *RegexpLiteral) String() string {

	return (fmt.Sprintf("/%s/%s", rl.Value, rl.Flags))
}

// BacktickLiteral holds details of a command to be executed
type BacktickLiteral struct {
	// Token is the actual token
	Token token.Token

	// Value is the name of the command to execute.
	Value string
}

func (bl *BacktickLiteral) ExpressionNode() {}

// GetToken returns the token.
func (bl *BacktickLiteral) GetToken() token.Token { return bl.Token }

// String returns this object as a string.
func (bl *BacktickLiteral) String() string { return bl.Token.Literal }

// ArrayLiteral holds an inline array
type ArrayLiteral struct {
	// Token is the token
	Token token.Token

	// Elements holds the members of the array.
	Elements []asti.ExpressionI
}

func (al *ArrayLiteral) ExpressionNode() {}

// GetToken returns the token.
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }

// String returns this object as a string.
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	fmt.Fprintf(&out, "[%v]", strings.Join(elements, ", "))
	return out.String()
}

// IndexExpression holds an index-expression
type IndexExpression struct {
	// Token is the actual token
	Token token.Token

	// Left is the thing being indexed.
	Left asti.ExpressionI

	// Index is the value we're indexing
	Index asti.ExpressionI
}

func (ie *IndexExpression) ExpressionNode() {}

// GetToken returns the token.
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }

// String returns this object as a string.
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "(%v[%v])", ie.Left, ie.Index)
	return out.String()
}

// HashLiteral holds a hash definition
type HashLiteral struct {
	// Token holds the token
	Token token.Token // the '{' token

	// Pairs stores the name/value sets of the hash-content
	Pairs map[asti.ExpressionI]asti.ExpressionI
}

func (hl *HashLiteral) ExpressionNode() {}

// GetToken returns the token.
func (hl *HashLiteral) GetToken() token.Token { return hl.Token }

// String returns this object as a string.
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := make([]string, 0)
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	fmt.Fprintf(&out, "{%v}", strings.Join(pairs, ", "))
	return out.String()
}

// CaseExpression handles the case within a switch statement
type CaseExpression struct {
	// Token is the actual token
	Token token.Token

	// Default branch?
	Default bool

	// The thing we match
	Expr []asti.ExpressionI

	// The code to execute if there is a match
	Block *BlockStatement
}

func (ce *CaseExpression) ExpressionNode() {}

// GetToken returns the token.
func (ce *CaseExpression) GetToken() token.Token { return ce.Token }

// String returns this object as a string.
func (ce *CaseExpression) String() string {
	var out bytes.Buffer

	if ce.Default {
		out.WriteString("default ")
	} else {
		out.WriteString("case ")

		tmp := []string{}
		for _, exp := range ce.Expr {
			tmp = append(tmp, exp.String())
		}
		out.WriteString(strings.Join(tmp, ","))
	}
	out.WriteString(ce.Block.String())
	return out.String()
}

// SwitchExpression handles a switch statement
type SwitchExpression struct {
	// Token is the actual token
	Token token.Token

	// Value is the thing that is evaluated to determine
	// which block should be executed.
	Value asti.ExpressionI

	// The branches we handle
	Choices []*CaseExpression
}

func (se *SwitchExpression) ExpressionNode() {}

// GetToken returns the token.
func (se *SwitchExpression) GetToken() token.Token { return se.Token }

// String returns this object as a string.
func (se *SwitchExpression) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "\nswitch (%v)\n{\n", se.Value)
	for _, tmp := range se.Choices {
		if tmp != nil {
			out.WriteString(tmp.String())
		}
	}
	out.WriteString("}\n")

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

// GetToken returns the token.
func (fes *ForeachStatement) GetToken() token.Token { return fes.Token }

// String returns this object as a string.
func (fes *ForeachStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "foreach %v %v %v", fes.Ident, fes.Value, fes.Body)
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
	Operator tokentype.TokenType
	Value    asti.ExpressionI
}

func (as *AssignStatement) ExpressionNode() {}

// GetToken returns the token.
func (as *AssignStatement) GetToken() token.Token { return as.Token }

// String returns this object as a string.
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%v%v%v", as.Name, tokentype.Attrib[as.Operator].String, as.Value)
	return out.String()
}

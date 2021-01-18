package asti

// NodeI reresents a node.
type NodeI interface {
	// TokenLiteral returns the literal of the token.
	TokenLiteral() string

	// String returns this object as a string.
	String() string
}

// StatementI represents a single statement.
type StatementI interface {
	// NodeI is the node holding the actual statement
	NodeI

	StatementNode()
}

// Expression represents a single expression.
type ExpressionI interface {
	// NodeI is the node holding the expression.
	NodeI
	ExpressionNode()
}

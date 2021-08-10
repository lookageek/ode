package ast

import "lookageek.com/ode/token"

// AST structure holds the two basic types of code lines, statements & expressions
// statements do not evaluate and present back a value, expressions do

// Node is the basic node of the AST which can hold a token
type Node interface {
	TokenLiteral() string
}

// Statement node holds a statement
type Statement interface {
	Node
	statementNode()
}

// Expression node holds an expression
type Expression interface {
	Node
	expressionNode()
}

// Program node is the root node of the AST
// and any Ode program is just a series of statements
// hence we have a slice of statements here
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// Identifier implements the Expression interface because in other
// places the identifier DOES produce a value like `let x = valueProducingIdentifier;`
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// The LetStatement should hold three values
// the variable name, the token of the let statement
// and the expression on the RHS of the let statement
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

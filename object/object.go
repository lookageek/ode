package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

// Object holds the literal values for integer, boolean and null
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer holds the literal integer value
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Boolean holds the literal boolean value
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null holds nothing
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue holds the object representation of
// evaluated return value in the return statement
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Error object holds the error encountered during the evaluation
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

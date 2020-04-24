/*
Package object represents the object system for our language. It defined the various structs needed
in order to track types defined in the lanauge and evaluate them.

We've opted to represent everything as an object (like ruby) making things easier to work with
*/
package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

// Object is a wrapper interface around the object system for our language.
// It represents each type in the language.
type Object interface {
	// Type returns the native type of the object
	Type() ObjectType
	// Inspect is used for debugging
	Inspect() string
}

// ---------- int ----------

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// ---------- bool ----------

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// ------- null -------

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// --------- return ---------

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

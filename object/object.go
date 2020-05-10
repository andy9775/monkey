/*
Package object represents the object system for our language. It defined the various structs needed
in order to track types defined in the lanauge and evaluate them.

We've opted to represent everything as an object (like ruby) making things easier to work with
*/
package object

import (
	"bytes"
	"fmt"
	"strings"

	"hash/fnv"

	"github.com/andy9775/monkey/ast"
	"github.com/andy9775/monkey/code"
)

type ObjectType string

const (
	INTEGER_OBJ           ObjectType = "INTEGER"
	BOOLEAN_OBJ                      = "BOOLEAN"
	NULL_OBJ                         = "NULL"
	RETURN_VALUE_OBJ                 = "RETURN_VALUE"
	ERROR_OBJ                        = "ERROR"
	FUNCTION_OBJ                     = "FUNCTION"
	COMPILED_FUNCTION_OBJ            = "COMPILED_FUNCTION_OBJ"
	STRING_OBJ                       = "STRING"
	ARRAY_OBJ                        = "ARRAY"

	BUILTIN_OBJ = "BUILTIN"

	CLOSURE_OBJ = "CLOSURE"

	HASH_OBJ = "HASH"
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

// ---------- string ---------

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

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

// ---------- error ----------

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// ------------- func -------------

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	/*
		Each function has it's own environment which allows for closures
		By default we include the outer environment (which gives the function access to
		variables in the scope that it is defined)
	*/
	Env *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

//CompiledFunction contains a series of instructions which make up a function body
// it is used for the vm/compiler
type CompiledFunction struct {
	Instructions code.Instructions
	// NumLocals specifies the number of locally scope variables this function uses
	NumLocals int

	// NumParameters specifies how many arguments this function expects
	NumParameters int
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJ }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() ObjectType { return CLOSURE_OBJ }
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}

// --------- built in funcs ---------

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// ---------- array ----------

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// --------- hash ---------

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// --------- hash key ---------

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

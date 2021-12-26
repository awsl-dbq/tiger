package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/awsl-dbq/tiger/tiscript/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	QUOTE_OBJ        = "QUOTE"
	MACRO_OBJ        = "MACRO"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	ToBytesArray() []byte
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) ToBytesArray() []byte {
	return []byte{}
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}
func (i *Boolean) ToBytesArray() []byte {
	return []byte{}
}

type NULL struct {
}

func (i *NULL) Type() ObjectType { return NULL_OBJ }

func (i *NULL) Inspect() string {
	return "null"
}
func (i *NULL) ToBytesArray() []byte {
	return []byte{}
}

type ReturnValue struct {
	Value Object
}

func (i *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (i *ReturnValue) Inspect() string {
	return i.Value.Inspect()
}
func (i *ReturnValue) ToBytesArray() []byte {
	return []byte{}
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (i *Error) ToBytesArray() []byte {
	return []byte{}
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func FromBytesArray(input []byte) Object {
	return nil
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
func (i *Function) ToBytesArray() []byte {
	return []byte{}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (i *String) ToBytesArray() []byte {
	return []byte{}
}

/**
 build -in
**/
type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (i *Builtin) ToBytesArray() []byte {
	return []byte{}
}

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
func (i *Array) ToBytesArray() []byte {
	return []byte{}
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {
	return "QUOTE(" + q.Node.String() + ")"
}
func (i *Quote) ToBytesArray() []byte {
	return []byte{}
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (i *Macro) ToBytesArray() []byte {
	return []byte{}
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}

package object

import (
	"bytes"
	"encoding/json"
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

func FromBytesArray(input []byte) Object {
	t := input[0]
	bs := input[1:]
	var obj Object
	switch t {
	case byte(0x01):
		obj = &Integer{}
	case byte(0x02):
		obj = &Boolean{}
	case byte(0x03):
		obj = &NULL{}
	case byte(0x04):
		obj = &ReturnValue{}
	case byte(0x05):
		obj = &Error{}
	case byte(0x06):
		obj = &Function{}
	case byte(0x07):
		obj = &String{}
	case byte(0x08):
		obj = &Builtin{}
	case byte(0x09):
		obj = &Array{}
	case byte(0x10):
		obj = &Quote{}
	case byte(0x11):
		obj = &Macro{}
	}
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(obj.Inspect())
	return obj
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x01)
	ret = append(ret, bs...)
	return ret
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}
func (i *Boolean) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x02)
	ret = append(ret, bs...)
	return ret
}

type NULL struct {
}

func (i *NULL) Type() ObjectType { return NULL_OBJ }

func (i *NULL) Inspect() string {
	return "null"
}
func (i *NULL) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x03)
	ret = append(ret, bs...)
	return ret
}

type ReturnValue struct {
	Value Object
}

func (i *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (i *ReturnValue) Inspect() string {
	return i.Value.Inspect()
}
func (i *ReturnValue) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x04)
	ret = append(ret, bs...)
	return ret
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (i *Error) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x05)
	ret = append(ret, bs...)
	return ret
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
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
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x06)
	ret = append(ret, bs...)
	return ret
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (i *String) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x07)
	ret = append(ret, bs...)
	return ret
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
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x08)
	ret = append(ret, bs...)
	return ret
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
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x09)
	ret = append(ret, bs...)
	return ret
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {
	return "QUOTE(" + q.Node.String() + ")"
}
func (i *Quote) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x10)
	ret = append(ret, bs...)
	return ret
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (i *Macro) ToBytesArray() []byte {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1, 1+len(bs))
	ret[0] = byte(0x11)
	ret = append(ret, bs...)
	return ret
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

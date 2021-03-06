package value

import (
	types "github.com/awsl-dbq/tiger/tiscript/llvm/types"
	llvmValue "github.com/llir/llvm/ir/value"
)

type Value struct {
	Type  types.Type
	Value llvmValue.Value

	// Is true when Value points to an LLVM Allocated variable, and is false
	// when the value is a constant.
	// This is used to know if a "load" instruction is neccesary or not.
	// Pointers are not considered "variables" in this context.
	IsVariable bool

	// Is used when returning multiple types from a function
	// Type is set to MultiValue when this is case, and will also contain the
	// type information
	MultiValues []Value
}

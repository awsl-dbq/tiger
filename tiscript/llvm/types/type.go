package types

import (
	"fmt"

	"github.com/awsl-dbq/tiger/tiscript/llvm/strings"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	llvmtypes "github.com/llir/llvm/ir/types"
	llvmValue "github.com/llir/llvm/ir/value"
)

var I8 = &Int{Type: llvmtypes.I8, TypeName: "int8", TypeSize: 8 / 8}
var I16 = &Int{Type: llvmtypes.I16, TypeName: "int16", TypeSize: 18 / 8}
var I32 = &Int{Type: llvmtypes.I32, TypeName: "int32", TypeSize: 32 / 8}
var I64 = &Int{Type: llvmtypes.I64, TypeName: "int64", TypeSize: 64 / 8}

var Void = &VoidType{}
var String = &StringType{}
var Bool = &BoolType{}

type Type interface {
	LLVM() llvmtypes.Type
	Name() string

	// Size of type in bytes
	Size() int64

	AddMethod(string, *Method)
	GetMethod(string) (*Method, bool)

	Zero(*ir.Block, llvmValue.Value)
}

type backingType struct {
	methods map[string]*Method
}

func (b *backingType) AddMethod(name string, method *Method) {
	if b.methods == nil {
		b.methods = make(map[string]*Method)
	}
	b.methods[name] = method
}

func (b *backingType) GetMethod(name string) (*Method, bool) {
	m, ok := b.methods[name]
	return m, ok
}

func (backingType) Size() int64 {
	panic("Type does not have size set")
}

func (backingType) Zero(*ir.Block, llvmValue.Value) {
	// NOOP
}

type Struct struct {
	backingType

	Members       map[string]Type
	MemberIndexes map[string]int

	IsHeapAllocated bool

	SourceName string
	Type       llvmtypes.Type
}

func (s Struct) LLVM() llvmtypes.Type {
	return s.Type
}

func (s Struct) Name() string {
	return fmt.Sprintf("struct(%s)", s.SourceName)
}

func (s Struct) Zero(block *ir.Block, alloca llvmValue.Value) {
	for key, valType := range s.Members {
		ptr := block.NewGetElementPtr(alloca,
			constant.NewInt(llvmtypes.I32, 0),
			constant.NewInt(llvmtypes.I32, int64(s.MemberIndexes[key])),
		)
		valType.Zero(block, ptr)
	}
}

func (s Struct) Size() int64 {
	var sum int64
	for _, valType := range s.Members {
		sum += valType.Size()
	}
	return sum
}

type Method struct {
	backingType

	Function        *Function
	PointerReceiver bool
	MethodName      string
}

func (m Method) LLVM() llvmtypes.Type {
	return m.Function.LLVM()
}

func (m Method) Name() string {
	return m.MethodName
}

type Function struct {
	backingType

	LlvmFunction llvmValue.Named

	// The return type of the LLVM function (is always 1)
	LlvmReturnType Type
	// Return types of the Tre function
	ReturnTypes []Type

	IsVariadic    bool
	ArgumentTypes []Type
	IsExternal    bool

	// Is used when calling an interface method
	JumpFunction *ir.Func
}

func (f Function) LLVM() llvmtypes.Type {
	return f.LlvmFunction.Type()
}

func (f Function) Name() string {
	return "func"
}

type BoolType struct {
	backingType
}

func (BoolType) LLVM() llvmtypes.Type {
	return llvmtypes.I1
}

func (BoolType) Name() string {
	return "bool"
}

func (BoolType) Size() int64 {
	return 1
}

func (b BoolType) Zero(block *ir.Block, alloca llvmValue.Value) {
	block.NewStore(constant.NewInt(llvmtypes.I1, 0), alloca)
}

type VoidType struct {
	backingType
}

func (VoidType) LLVM() llvmtypes.Type {
	return llvmtypes.Void
}

func (VoidType) Name() string {
	return "void"
}

func (VoidType) Size() int64 {
	return 0
}

type Int struct {
	backingType

	Type     *llvmtypes.IntType
	TypeName string
	TypeSize int64
}

func (i Int) LLVM() llvmtypes.Type {
	return i.Type
}

func (i Int) Name() string {
	return i.TypeName
}

func (i Int) Size() int64 {
	return i.TypeSize
}

func (i Int) Zero(block *ir.Block, alloca llvmValue.Value) {
	block.NewStore(constant.NewInt(i.Type, 0), alloca)
}

type StringType struct {
	backingType
	Type llvmtypes.Type
}

// Populated by compiler.go
var ModuleStringType llvmtypes.Type
var EmptyStringConstant *ir.Global

func (StringType) LLVM() llvmtypes.Type {
	return ModuleStringType
}

func (StringType) Name() string {
	return "string"
}

func (StringType) Size() int64 {
	return 16
}

func (s StringType) Zero(block *ir.Block, alloca llvmValue.Value) {
	lenPtr := block.NewGetElementPtr(alloca, constant.NewInt(llvmtypes.I32, 0), constant.NewInt(llvmtypes.I32, 0))
	backingDataPtr := block.NewGetElementPtr(alloca, constant.NewInt(llvmtypes.I32, 0), constant.NewInt(llvmtypes.I32, 1))
	block.NewStore(constant.NewInt(llvmtypes.I64, 0), lenPtr)
	block.NewStore(strings.Toi8Ptr(block, EmptyStringConstant), backingDataPtr)
}

type Array struct {
	backingType
	Type     Type
	Len      uint64
	LlvmType llvmtypes.Type
}

func (a Array) LLVM() llvmtypes.Type {
	return a.LlvmType
}

func (a Array) Name() string {
	return "array"
}

func (a Array) Zero(block *ir.Block, alloca llvmValue.Value) {
	for i := uint64(0); i < a.Len; i++ {
		ptr := block.NewGetElementPtr(alloca, constant.NewInt(llvmtypes.I64, 0), constant.NewInt(llvmtypes.I64, int64(i)))
		a.Type.Zero(block, ptr)
	}
}

type Slice struct {
	backingType
	Type     Type // type of the items in the slice []int => int
	LlvmType llvmtypes.Type
}

func (s Slice) LLVM() llvmtypes.Type {
	return s.LlvmType
}

func (Slice) Name() string {
	return "slice"
}

func (Slice) Size() int64 {
	return 3*4 + 8 // 3 int32s and a pointer
}

// func (s Slice) SliceZero(block *ir.Block, mallocFunc llvmValue.Named, initCap int) *ir.InstAlloca {
// 	// The cap must always be larger than 0
// 	// Use 2 as the default value
// 	if initCap < 2 {
// 		initCap = 2
// 	}

// 	emptySlize := block.NewAlloca(s.LLVM())

// 	len := block.NewGetElementPtr(emptySlize, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
// 	len.SetName(name.Var("len"))
// 	cap := block.NewGetElementPtr(emptySlize, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
// 	cap.SetName(name.Var("cap"))
// 	offset := block.NewGetElementPtr(emptySlize, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 2))
// 	offset.SetName(name.Var("offset"))
// 	backingArray := block.NewGetElementPtr(emptySlize, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 3))
// 	backingArray.SetName(name.Var("backing"))

// 	block.NewStore(constant.NewInt(types.I32, 0), len)
// 	block.NewStore(constant.NewInt(types.I32, int64(initCap)), cap)
// 	block.NewStore(constant.NewInt(types.I32, 0), offset)

// 	mallocatedSpaceRaw := block.NewCall(mallocFunc, constant.NewInt(types.I64, int64(initCap)*s.Type.Size()))
// 	mallocatedSpaceRaw.SetName(name.Var("slicezero"))
// 	bitcasted := block.NewBitCast(mallocatedSpaceRaw, types.NewPointer(s.Type.LLVM()))
// 	block.NewStore(bitcasted, backingArray)

// 	return emptySlize
// }

type Pointer struct {
	backingType

	Type                  Type
	IsNonAllocDereference bool

	LlvmType llvmtypes.Type
}

func (p Pointer) LLVM() llvmtypes.Type {
	return llvmtypes.NewPointer(p.Type.LLVM())
}

func (p Pointer) Name() string {
	return fmt.Sprintf("pointer(%s)", p.Type.Name())
}

func (p Pointer) Size() int64 {
	return 8
}

// MultiValue is used when returning multiple values from a function
type MultiValue struct {
	backingType
	Types []Type
}

func (m MultiValue) Name() string {
	return "multivalue"
}

func (m MultiValue) LLVM() llvmtypes.Type {
	panic("MutliValue has no LLVM type")
}

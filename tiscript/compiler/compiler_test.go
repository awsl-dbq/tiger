package compiler

import (
	"testing"

	"github.com/awsl-dbq/tiger/tiscript/lexer"
	ir "github.com/awsl-dbq/tiger/tiscript/llvm"
	"github.com/awsl-dbq/tiger/tiscript/parser"
)

func TestStringIR(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		ir       string
	}{
		{"10", 10, "i32 10"},
		{"true", true, "i1 true"},
		{"false", false, "i1 false"},
		// {"if (true) { 10 }", 10},
		// {"if (false) { 10 }", nil},
		// {"if (1) { 10 }", 10},
		// {"if (1 < 2) { 10 }", 10},
		// {"if (1 > 2) { 10 }", nil}, {"if (1 > 2) { 10 } else { 20 }", 20},
		// {"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testCompile(tt.input)
		switch integer := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(integer), tt.ir)
		case nil:
			testNullObject(t, evaluated)
		case bool:
			testBoolObjct(t, evaluated, integer, tt.ir)
		}
	}
}
func testCompile(input string) ir.IRObject {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Compile(program)
}

func testIntegerObject(t *testing.T, obj ir.IRObject, expected int64, expectedIr string) bool {
	result, ok := obj.(*ir.IntegerObject)
	if obj.IR() != expectedIr {
		t.Errorf("IRcode is %v, want %s", obj, expectedIr)
		return false
	}
	if !ok {
		t.Errorf("object is %T (%+v), want Integer", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Value is %d, want %d", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj ir.IRObject) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testBoolObjct(t *testing.T, obj ir.IRObject, expected bool, expectedIr string) bool {
	result, ok := obj.(*ir.BooleanObject)
	if obj.IR() != expectedIr {
		t.Errorf("IRcode is %v, want %s", obj.IR(), expectedIr)
		return false
	}
	if !ok {
		t.Errorf("object is %T (%+v), want Boolean", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Value is %v, want %v", result.Value, expected)
		return false
	}
	return true
}

func TestIR(t *testing.T) {
	MakeTargetTriple()
	input := "true"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	Compile(program)
	t.Errorf("\n%v", GetIR())
}

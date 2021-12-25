package ast

import (
	"testing"

	"github.com/awsl-dbq/tiger/tiscript/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "someVar"},
					Value: "someVar",
				},
			},
		},
	}
	if program.String() != "let myVar = someVar;" {
		t.Fatalf("expect 'let myVar = someVar;'but,got  '%v'", program.String())
	}
}

// func TestNodeType(t *testing.T) {
// 	program := &Program{
// 		Statements: []Statement{
// 			&NodeTypeLiteral{
// 				Token:    token.Token{Type: token.NODETYPE, Literal: "NodeType"},
// 				NodeName: "Page",
// 			},
// 		},
// 	}
// 	if program.String() != "NodeType  Page{}" {
// 		t.Fatalf("expect 'NodeType Page{};'but,got  '%v'", program.String())
// 	}
// }

func TestEdgeType(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&EdgeTypeLiteral{
				Token:       token.Token{Type: token.EDGETYPE, Literal: "EdgeType"},
				EdgeName:    "friend",
				SourceRefer: "st",
				SourceType:  "Person",
				TargetRefer: "tt",
				TargetType:  "Person",
				IsReverse:   true,
			},
		},
	}
	expt := `EdgeType  friend(st: Person,tt: Person)@reverse{
}`
	if program.String() != expt {
		t.Fatalf("expect \n%v\n but,got  '\n%v\n'", expt, program.String())
	}
}

package ast

import (
	"bytes"

	"github.com/awsl-dbq/tiger/tiscript/token"
)

type MakeType string //function or sth else
type ReturnObject struct {
	Token   token.Token
	Literal string
}
type MakeLiteral struct {
	Token  token.Token
	Type   MakeType
	Name   *Identifier
	Params []*FunctionParam
	Body   *BlockStatement
	Return ReturnObject
}

func (node *MakeLiteral) ExpressionNode() {}
func (node *MakeLiteral) TokenLiteral() string {
	return node.Token.Literal
}
func (node *MakeLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(node.Token.Literal + " ")
	if node.Type == "function" {
		out.WriteString("function ")
	}
	out.WriteString(node.Name.Value)
	out.WriteString("(")
	for i, v := range node.Params {
		out.WriteString(v.String())
		if i != len(node.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	out.WriteString(node.Return.Literal)
	out.WriteString("{\n")
	out.WriteString(node.Body.String() + "\n")
	out.WriteString("}")
	return out.String()
}

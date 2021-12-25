package ast

import (
	"bytes"

	"github.com/awsl-dbq/tiger/tiscript/token"
)

/**

EdgeType  referTo(s:Person,t:Person) @filter(has(s.Name)) @reverse {
    id,
    lable: s.Name,
}

**/
type EdgeTypeLiteral struct {
	Token       token.Token
	EdgeName    string
	SourceType  string
	SourceRefer string
	TargetType  string
	TargetRefer string
	IsReverse   bool
	Properties  []*EdgeProperty
	Filter      *EdgeFilter
}
type EdgeProperty struct {
	Name string
	Type string
	As   *AsExpression
}
type AsExpression struct {
	Literal string
	Tokens  []token.Token
}
type EdgeFilter struct {
	Literal    string
	FuncName   string
	FuncParams []*Identifier
}

func (node *EdgeTypeLiteral) expressionNode() {}
func (node *EdgeTypeLiteral) TokenLiteral() string {
	return node.Token.Literal
}
func (node *EdgeTypeLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("EdgeType  ")
	out.WriteString(node.EdgeName)
	out.WriteString("(")
	out.WriteString(node.SourceRefer + ": " + node.SourceType)
	out.WriteString(",")
	out.WriteString(node.TargetRefer + ": " + node.TargetType)
	out.WriteString(") ")
	if node.IsReverse {
		out.WriteString("@reverse ")
	}
	if node.Filter != nil {
		out.WriteString(node.Filter.String())
	}
	out.WriteString("{\n")
	// append property
	for _, property := range node.Properties {
		out.WriteString(property.Name + ": " + property.Type)
		if property.As != nil {
			out.WriteString(" as " + property.As.String())
		}
		out.WriteString(",\n")
	}
	out.WriteString("}")
	return out.String()
}
func (node *EdgeTypeLiteral) statementNode() {}

func (node *EdgeFilter) String() string {
	var out bytes.Buffer
	out.WriteString("@filter(")
	out.WriteString(node.FuncName + "(")
	for i, funcParam := range node.FuncParams {
		out.WriteString(funcParam.String())
		if i != len(node.FuncParams)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("))")
	return out.String()
}
func (node *AsExpression) String() string {
	return node.Literal
}

package evaluator

import (
	"github.com/awsl-dbq/tiger/tiscript/ast"
	"github.com/awsl-dbq/tiger/tiscript/object"
)

func evalMakeLiteral(node *ast.MakeLiteral, env *object.Environment) object.Object {
	params := []*ast.Identifier{}
	for _, p := range node.Params {
		params = append(params, &ast.Identifier{
			Value: p.Value,
		})
	}
	val := &object.Function{Parameters: params, Env: env, Body: node.Body}
	env.Set(node.Name.Value, val)
	return val
}

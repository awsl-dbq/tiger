package parser

import (
	"github.com/awsl-dbq/tiger/tiscript/ast"
	"github.com/awsl-dbq/tiger/tiscript/token"
)

func (p *Parser) parsePackageLiteral() ast.Expression {
	// package a.b.c.d ;
	ex := &ast.PackageLiteral{
		Token: p.curToken,
	}
	p.nextToken()
	if !p.peekTokenIs(token.SEMICOLON) {
		return nil
	}
	ex.Value = p.curToken.Literal
	return ex
}

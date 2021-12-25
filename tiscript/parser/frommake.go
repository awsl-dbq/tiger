package parser

import (
	"github.com/awsl-dbq/tiger/tiscript/ast"
	"github.com/awsl-dbq/tiger/tiscript/token"
)



func (p *Parser) parseIdentifierArrays() []*ast.Identifier {
	nts := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) { // return )
		p.nextToken()
		return nts
	}
	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	p.nextToken()
	if p.curTokenIs(token.RBRACKET) {
		return nts
	}
	nt := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	nts = append(nts, nt)
	// check , []
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if p.peekTokenIs(token.RBRACKET) {
			break
		}
		if p.peekTokenIs(token.COMMA) {
			continue
		}
		p.nextToken()

		nt := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		nts = append(nts, nt)
	}
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return nts
}
func (p *Parser) parseFromGraphEdgeTypes() []*ast.Identifier {
	if p.peekTokenIs(token.RPAREN) {
		return nil
	}
	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
	}
	return p.parseIdentifierArrays()
}
func (p *Parser) parseFromGraphNodeTypes() []*ast.Identifier {
	return p.parseIdentifierArrays()
}

func (p *Parser) parseMakeLiteral() ast.Expression {
	if !p.curTokenIs(token.MAKE) {
		p.addError("need toke.Make")
		return nil
	}
	ml := &ast.MakeLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.FUNCTION) {
		return nil
	}
	ml.Type = "function"
	if !p.expectPeek(token.IDENT) {
		p.addError("need ident")
		return nil
	}
	ml.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal}
	ml.Params = p.parseFunctionWithTypeParams()
	if !p.peekTokenIs(token.LBRACE) {
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		ml.Return.Literal = p.curToken.Literal
		ml.Return.Token = p.curToken
		p.nextToken()
	} else {
		ml.Return.Literal = "void"
		p.nextToken()
	}
	ml.Body = p.parseBlockStatement()
	return ml
}

func (p *Parser) parseFunctionWithTypeParams() []*ast.FunctionParam {
	params := []*ast.FunctionParam{}
	if p.peekTokenIs(token.RPAREN) {
		return params
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	parm := &ast.FunctionParam{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	p.nextToken()
	if !p.curTokenIs(token.IDENT) {
		p.addError("Type should be ident")
		return nil
	}
	parm.Type = p.curToken.Literal
	params = append(params, parm)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		parm := &ast.FunctionParam{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			p.addError("Type should be ident")
			return nil
		}
		parm.Type = p.curToken.Literal
		params = append(params, parm)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

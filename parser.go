package main

import (
	"fmt"

	"github.com/rabierre/compiler/ast"
	"github.com/rabierre/compiler/token"
)

var debug bool

type Parser struct {
	val string
	tok token.Type
	pos int

	scope    *ast.Scope
	topScope *ast.Scope
	scanner  *Scanner

	comments *ast.CommentList

	// If we handle source codes in files
	// This should go in file struct
	//
	decls      []ast.Decl
	UnResolved []*ast.Ident
}

func (p *Parser) Init(src []byte) {
	p.scanner = &Scanner{}
	p.scanner.Init()
	p.scanner.src = src
	p.comments = &ast.CommentList{}
	p.OpenScope() // Top scope
	p.topScope = p.scope

	p.next()
}

func (p *Parser) Parse() {
	p.OpenScope()

	// TODO parse comment
	// case token.COMMENT:
	// 		p.parseComment()

	for !p.scanner.fullScaned {
		// TODO use p.peek() after move parsecomment phase to scanner

		switch p.tok {
		default:
			p.parseDecl()
		}
	}

	p.CloseScope()

	unResolved := p.UnResolved
	p.UnResolved = []*ast.Ident{}
	old := p.scope
	p.scope = p.topScope
	for _, id := range unResolved {
		p.resolve(id)
	}
	p.scope = old

	if debug {
		for _, un := range p.UnResolved {
			print(un.Name, " ")
		}
		println()
	}

	if len(p.UnResolved) > 0 {
		panic("Unresolved ident exist")
	}
}

func (p *Parser) parseDecl() {
	trace("parseDecl")

	switch p.tok {
	// By spec for now, no global variable, no imports are available.
	// Function is top scope
	//
	case token.FUNC:
		p.parseFunc()
	default:
	}
}

func (p *Parser) parseFunc() ast.Decl {
	trace("parseFunc")

	p.next() // consune func token
	ident := p.parseIdent()

	p.expect(token.LPAREN)
	params := p.parseParamList()
	p.expect(token.RPAREN)

	// TODO parse func type
	var _typ token.Type
	if p.tok == token.INT || p.tok == token.DOUBLE {
		_typ = p.tok
		p.next() // consume type token
	} else {
		_typ = token.VOID
	}

	p.OpenScope()
	for _, param := range params.List {
		if decl := param.(*ast.VarDeclStmt); decl != nil {
			obj := ast.NewObject(decl, ast.VAR)
			p.scope.Insert(obj, decl.Name.Name)
		}
	}

	body := p.parseBody()
	decl := &ast.FuncDecl{Name: ident, Body: body, Params: params, Type: _typ}

	// TODO move this to specific function like parse function decl only
	p.decls = append(p.decls, decl)

	// TODO generalize this
	// p.declare(decl, nil, p.pkgScope, ast.Fun, ident)
	obj := ast.NewObject(decl, ast.FUNC)
	p.topScope.Insert(obj, ident.Name)

	return decl
}

func (p *Parser) parseIdent() *ast.Ident {
	trace("parseIdent")

	if p.tok != token.IDENT {
		panic("Expect IDENT, GOT: " + p.tok.String())
	}

	id := &ast.Ident{Name: p.val, Pos: p.pos}
	p.next()

	return id
}

func (p *Parser) parseParamList() *ast.StmtList {
	trace("parseParamList")

	list := []ast.Stmt{}
	for p.tok == token.INT || p.tok == token.DOUBLE {
		list = append(list, p.parseParam())

		if p.tok == token.RPAREN {
			break
		}
		if p.tok == token.COMMA {
			p.next() // consume ,
		}
	}

	return &ast.StmtList{List: list}
}

func (p *Parser) parseParam() ast.Stmt {
	trace("parseParam")

	param := &ast.VarDeclStmt{Pos: p.pos, Type: p.tok}
	p.next() // consume type
	param.Name = &ast.Ident{Pos: p.pos, Name: p.val}
	p.next() // consume variable
	return param
}

func (p *Parser) parseBody() *ast.CompoundStmt {
	trace("parseBody")

	lbrace := p.expect(token.LBRACE)

	list := p.parseStmtList()
	p.CloseScope()

	rbrace := p.expect(token.RBRACE)

	return &ast.CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseCompoundStmt() *ast.CompoundStmt {
	trace("parseCompoundStmt")

	lbrace := p.expect(token.LBRACE)
	p.OpenScope()

	list := p.parseStmtList()

	p.CloseScope()
	rbrace := p.expect(token.RBRACE)

	return &ast.CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseStmtList() []ast.Stmt {
	trace("parseStmtList")

	list := []ast.Stmt{}
	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseStmt())
		if p.tok == token.COMMA {
			p.next() // consume ,
		}
	}
	return list
}

func (p *Parser) parseStmt() ast.Stmt {
	trace("parseStmt")

	switch p.tok {
	case token.INT, token.DOUBLE:
		return p.parseVarDecl()
	case token.IDENT:
		return p.parseExprStmt()
	case token.FOR:
		return p.parseForStmt()
	case token.IF:
		return p.parseIfStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.LBRACE:
		return p.parseCompoundStmt()
	case token.RBRACE:
		return &ast.EmptyStmt{ /*position for semicolon if need*/ }
	default:
		// unkown type but progress
		p.next()
	}

	return &ast.BadStmt{From: p.pos}
}

// parse variable declaration
// int a = 1
// double b = 1.0
// int c
//
func (p *Parser) parseVarDecl() ast.Stmt {
	trace("parseVarDecl")

	decl := &ast.VarDeclStmt{Pos: p.pos, Type: p.tok}
	p.next() // consume type

	ident := p.parseIdent()

	var value ast.Expr
	if p.tok == token.ASSIGN {
		p.next() //consume =
		value = p.parseExpr(true)
	}

	decl.Name = ident
	decl.RValue = value

	// TODO generalize this
	obj := ast.NewObject(decl, ast.VAR)
	p.scope.Insert(obj, ident.Name)

	return decl
}

// Parse Expr in Statement
// a = 10
// funcCall()
//
func (p *Parser) parseExprStmt() ast.Stmt {
	trace("parseExprStmt")

	x := p.parseExpr(true)
	if p.tok == token.ASSIGN {
		p.next() // consume =
		y := p.parseExpr(true)
		x = &ast.AssignExpr{Pos: x.(*ast.Ident).Pos, LValue: x, RValue: y}
	}
	return &ast.ExprStmt{Val: x}
}

func (p *Parser) parseForStmt() ast.Stmt {
	trace("parseForStmt")

	pos := p.pos
	p.next() //consume for

	p.expect(token.LPAREN)

	_init := p.parseStmt()
	p.expect(token.SEMI_COLON)

	_cond := p.parseExpr(true)
	p.expect(token.SEMI_COLON)

	_post := p.parseExpr(true)

	p.expect(token.RPAREN)

	// TODO extract to function
	// if decl := _init.(*VarDeclStmt); decl != nil {
	// 	obj := &Object{decl: decl}
	// 	p.scope.Insert(obj, decl.Name.Name)
	// }

	body := p.parseBody()

	return &ast.ForStmt{Pos: pos, Cond: _cond, Init: _init, Post: _post, Body: body}
}

func (p *Parser) parseIfStmt() ast.Stmt {
	trace("parseIfStmt")

	pos := p.pos
	p.next() // consume if

	p.expect(token.LPAREN)
	cond := p.parseExpr(true) // TODO: If cond is nil, error
	p.expect(token.RPAREN)

	body := p.parseCompoundStmt()

	var elseBody ast.Stmt
	if p.tok == token.ELSE {
		p.next() // consume else
		elseBody = p.parseCompoundStmt()
	}
	return &ast.IfStmt{Pos: pos, Cond: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseReturnStmt() ast.Stmt {
	trace("parseReturnStmt")

	pos := p.pos
	p.next() // consume return

	var expr ast.Expr
	if p.tok != token.EOF && p.tok != token.RBRACE {
		expr = p.parseExpr(true)
	}

	return &ast.ReturnStmt{Pos: pos, Value: expr}
}

func (p *Parser) parseExprList() *ast.ExprList {
	trace("parseExprList")

	var exprs []ast.Expr
	for {
		// 1. parse expr
		exprs = append(exprs, p.parseExpr(true))
		// 2. if next token is COMMA then continue
		// 2-1. else return
		if p.tok != token.COMMA {
			break
		}
		p.next() // consume ,
	}
	return &ast.ExprList{List: exprs}
}

func (p *Parser) parseExpr(lookup bool) ast.Expr {
	trace("parseExpr")

	return p.parseBinaryExpr(token.LowestPriority+1, lookup)
}

// Term
func (p *Parser) parseBinaryExpr(prio int, lookup bool) ast.Expr {
	trace("parseBinaryExpr")

	x := p.parseUnaryExpr(lookup)
	for {
		// check if current token is operator
		if p.tok.Priority() < prio {
			return x
		}

		op := ast.Operator{Type: p.tok}
		p.next() // consume operator

		y := p.parseBinaryExpr(p.tok.Priority()+1, lookup)
		x = &ast.BinaryExpr{Pos: p.pos, Op: op, LValue: x, RValue: y}
	}
}

// Factor ::= "(" Expr ")"
//         | AddSub Factor
//         | number
//         | string
func (p *Parser) parseUnaryExpr(lookup bool) ast.Expr {
	trace("parseUnaryExpr")

	switch p.tok {
	case token.PLUS, token.MINUS:
		op := ast.Operator{Type: p.tok}
		p.next() // consume operator

		x := p.parseUnaryExpr(lookup)

		// TODO fix position: use operator's position
		return &ast.UnaryExpr{Pos: p.pos, Op: op, RValue: x}
	}

	return p.parsePrimaryExpr(lookup)
}

// Parse function call or variable reference
// identifier "(" ExprList ? ")"
// identifier
//
func (p *Parser) parsePrimaryExpr(lookup bool) ast.Expr {
	trace("parsePrimaryExpr")

	x := p.parseOperand(lookup)

	switch p.tok {
	case token.LPAREN:
		if !lookup {
			p.resolve(x)
		}

		return p.parseCallExpr(x)
	case token.INC, token.DEC:
		if lookup {
			p.resolve(x)
		}

		op := ast.Operator{Type: p.tok}
		// TODO shortExpr is similar with unaryExpr()
		x = &ast.ShortExpr{Pos: p.pos, Op: op, RValue: x}
		p.next()
	}

	return x
}

func (p *Parser) parseOperand(lookup bool) ast.Expr {
	trace("parseOperand")

	switch p.tok {
	case token.IDENT:
		x := p.parseIdent()
		if lookup {
			p.resolve(x)
		}
		return x
	case token.INT_LIT, token.DOUBLE_LIT, token.TRUE, token.FALSE:
		lit := &ast.BasicLit{Pos: p.pos, Value: p.val, Type: p.tok}
		p.next()
		return lit
	default:
		p.next()
	}

	return &ast.BadExpr{From: p.pos, To: p.pos}
}

func (p *Parser) parseCallExpr(x ast.Expr) ast.Expr {
	trace("parseCallExpr")

	lparen := p.expect(token.LPAREN)
	list := []ast.Expr{}
	for p.tok != token.EOF && p.tok != token.RPAREN {
		list = append(list, p.parseRHS())

		if p.tok == token.COMMA {
			p.next() // consume ,
		}
	}
	params := &ast.ExprList{list}
	rparen := p.expect(token.RPAREN)

	return &ast.CallExpr{Name: x, LParenPos: lparen, RParenPos: rparen, Params: params}
}

func (p *Parser) parseRHS() ast.Expr {
	trace("parseRHS")

	return p.parseExpr(true)
}

func (p *Parser) next() {
	for tok, _ := p.scanner.peek(); tok.Kind == token.COMMENT; tok, _ = p.scanner.peek() {
		p.scanner.nextLine()
	}
	tok, pos := p.scanner.next()
	p.tok = tok.Kind
	p.val = tok.Val
	p.pos = pos
}

func (p *Parser) expect(expected token.Type) int {
	pos := p.pos
	if p.tok != expected {
		panic(fmt.Sprintf("Expected: %s Found: %s, %s", expected.String(), p.val, p.tok.String()))
	}
	p.next()
	return pos
}

func (p *Parser) parseComment() {
	trace("parseComment")

	token, pos := p.scanner.nextLine()

	comment := &ast.Comment{pos, token.Val}
	p.comments.Insert(comment)
}

func (p *Parser) resolve(expr ast.Expr) {
	trace("resolve")

	id := expr.(*ast.Ident)
	if id == nil {
		return
	}

	for s := p.scope; s != nil; s = s.Outer {
		_, exist := s.Objects[id.Name]
		if exist {
			return
		}
	}

	p.UnResolved = append(p.UnResolved, id)
}

func (p *Parser) OpenScope() {
	p.scope = &ast.Scope{p.scope, map[string]*ast.Object{}}
}

func (p *Parser) CloseScope() {
	p.scope = p.scope.Outer
}

func trace(name string) {
	if debug {
		println(name)
	}
}

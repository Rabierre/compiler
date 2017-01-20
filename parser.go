package main

import (
	"fmt"

	"github.com/rabierre/compiler/token"
)

const debug = true

type Parser struct {
	val string
	tok token.Type
	pos int

	scope    *Scope
	topScope *Scope
	scanner  *Scanner

	comments *CommentList

	// If we handle source codes in files
	// This should go in file struct
	//
	decls      []Decl
	UnResolved []*Ident
}

func (p *Parser) Init(src []byte) {
	p.scanner = &Scanner{}
	p.scanner.Init()
	p.scanner.src = src
	p.comments = &CommentList{}
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
	p.UnResolved = []*Ident{}
	old := p.scope
	p.scope = p.topScope
	for _, id := range unResolved {
		p.resolve(id)
	}
	p.scope = old

	for _, un := range p.UnResolved {
		print(un.Name, " ")
	}
	println()
	if len(p.UnResolved) != 0 {
		panic("Unresolved ident exist")
	}
}

func (p *Parser) parseDecl() {
	println("parseDecl")

	switch p.tok {
	// By spec for now, no global variable, no imports are available.
	// Function is top scope
	//
	case token.FUNC:
		p.parseFunc()
	default:
	}
}

func (p *Parser) parseFunc() Decl {
	println("parseFunc")
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
		if decl := param.(*VarDeclStmt); decl != nil {
			obj := &Object{decl: decl.Name}
			p.scope.Insert(obj, decl.Name.Name)
		}
	}

	body := p.parseBody()
	funcDecl := &FuncDecl{Name: ident, Body: body, Params: params, Type: _typ}

	// TODO move this to specific function like parse function decl only
	p.decls = append(p.decls, funcDecl)

	// TODO generalize this
	// p.declare(decl, nil, p.pkgScope, ast.Fun, ident)
	obj := &Object{decl: funcDecl, kind: FUNC}
	p.topScope.Insert(obj, ident.Name)

	return funcDecl
}

func (p *Parser) parseIdent() *Ident {
	println("parseIdent")
	if p.tok != token.IDENT {
		panic("Expect IDENT, GOT: " + p.tok.String())
	}

	id := &Ident{Name: p.val, Pos: p.pos}
	p.next()

	return id
}

func (p *Parser) parseParamList() *StmtList {
	println("parseParamList")

	list := []Stmt{}
	for p.tok == token.INT || p.tok == token.DOUBLE {
		list = append(list, p.parseParam())

		if p.tok == token.RPAREN {
			break
		}
		if p.tok == token.COMMA {
			p.next() // consume ,
		}
	}

	return &StmtList{List: list}
}

func (p *Parser) parseParam() Stmt {
	println("parseParam")
	param := &VarDeclStmt{Pos: p.pos, Type: p.tok}
	p.next() // consume type
	param.Name = &Ident{Pos: p.pos, Name: p.val}
	p.next() // consume variable
	return param
}

func (p *Parser) parseBody() *CompoundStmt {
	println("parseBody")
	lbrace := p.expect(token.LBRACE)
	// p.OpenScope()

	list := p.parseStmtList()
	p.CloseScope()

	rbrace := p.expect(token.RBRACE)

	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseCompoundStmt() *CompoundStmt {
	println("parseCompoundStmt")
	lbrace := p.expect(token.LBRACE)
	p.OpenScope()

	list := p.parseStmtList()

	p.CloseScope()
	rbrace := p.expect(token.RBRACE)

	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseStmtList() []Stmt {
	println("parseStmtList")
	list := []Stmt{}
	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseStmt())
		if p.tok == token.COMMA {
			p.next() // consume ,
		}
	}
	return list
}

func (p *Parser) parseStmt() Stmt {
	println("parseStmt")

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
		return &EmptyStmt{ /*position for semicolon if need*/ }
	default:
		// unkown type but progress
		p.next()
	}

	return &BadStmt{From: p.pos}
}

// parse variable declaration
// int a = 1
// double b = 1.0
// int c
//
func (p *Parser) parseVarDecl() Stmt {
	println("parseVarDecl")
	decl := &VarDeclStmt{Pos: p.pos, Type: p.tok}
	p.next() // consume type

	ident := p.parseIdent()

	var value Expr
	if p.tok == token.ASSIGN {
		p.next() //consume =
		value = p.parseExpr(true)
	}

	decl.Name = ident
	decl.RValue = value

	// TODO generalize this
	obj := &Object{decl: decl, kind: VAR}
	p.scope.Insert(obj, ident.Name)

	return decl
}

// Parse Expr in Statement
// a = 10
// funcCall()
//
func (p *Parser) parseExprStmt() Stmt {
	println("parseExprStmt")
	x := p.parseExpr(true)
	if p.tok == token.ASSIGN {
		p.next() // consume =
		y := p.parseExpr(true)
		x = &AssignExpr{Pos: x.(*Ident).Pos, LValue: x, RValue: y}
	}
	return &ExprStmt{expr: x}
}

func (p *Parser) parseForStmt() Stmt {
	println("parseForStmt")
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

	return &ForStmt{Pos: pos, Cond: _cond, Init: _init, Post: _post, Body: body}
}

func (p *Parser) parseIfStmt() Stmt {
	println("parseIfStmt")
	pos := p.pos
	p.next() // consume if

	p.expect(token.LPAREN)
	cond := p.parseExpr(true) // TODO: If cond is nil, error
	p.expect(token.RPAREN)

	body := p.parseCompoundStmt()

	var elseBody Stmt
	if p.tok == token.ELSE {
		p.next() // consume else
		elseBody = p.parseCompoundStmt()
	}
	return &IfStmt{Pos: pos, Cond: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseReturnStmt() Stmt {
	println("parseReturnStmt")
	pos := p.pos
	p.next() // consume return

	var expr Expr
	if p.tok != token.EOF && p.tok != token.RBRACE {
		expr = p.parseExpr(true)
	}

	return &ReturnStmt{Pos: pos, Value: expr}
}

func (p *Parser) parseExprList() *ExprList {
	println("parseExprList")
	var exprs []Expr
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
	return &ExprList{List: exprs}
}

func (p *Parser) parseExpr(lookup bool) Expr {
	println("parseExpr")
	return p.parseBinaryExpr(token.LowestPriority+1, lookup)
}

// Term
func (p *Parser) parseBinaryExpr(prio int, lookup bool) Expr {
	println("parseBinaryExpr")

	x := p.parseUnaryExpr(lookup)
	for {
		// check if current token is operator
		if p.tok.Priority() < prio {
			return x
		}

		op := Operator{Type: p.tok}
		p.next() // consume operator

		y := p.parseBinaryExpr(p.tok.Priority()+1, lookup)
		x = &BinaryExpr{Pos: p.pos, Op: op, LValue: x, RValue: y}
	}
}

// Factor ::= "(" Expr ")"
//         | AddSub Factor
//         | number
//         | string
func (p *Parser) parseUnaryExpr(lookup bool) Expr {
	println("parseUnaryExpr")

	switch p.tok {
	case token.PLUS, token.MINUS:
		op := Operator{Type: p.tok}
		p.next() // consume operator

		x := p.parseUnaryExpr(lookup)

		// TODO fix position: use operator's position
		return &UnaryExpr{Pos: p.pos, Op: op, RValue: x}
	}

	return p.parsePrimaryExpr(lookup)
}

// Parse function call or variable reference
// identifier "(" ExprList ? ")"
// identifier
//
func (p *Parser) parsePrimaryExpr(lookup bool) Expr {
	println("parsePrimaryExpr")
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

		op := Operator{Type: p.tok}
		// TODO shortExpr is similar with unaryExpr()
		x = &ShortExpr{Pos: p.pos, Op: op, RValue: x}
		p.next()
	}

	return x
}

func (p *Parser) parseOperand(lookup bool) Expr {
	println("parseOperand")

	switch p.tok {
	case token.IDENT:
		x := p.parseIdent()
		if lookup {
			p.resolve(x)
		}
		return x
	case token.INT_LIT, token.DOUBLE_LIT, token.TRUE, token.FALSE:
		lit := &BasicLit{Pos: p.pos, Value: p.val, Type: p.tok}
		p.next()
		return lit
	default:
		p.next()
	}

	return &BadExpr{From: p.pos, To: p.pos}
}

func (p *Parser) parseCallExpr(x Expr) Expr {
	println("parseCallExpr")
	lparen := p.expect(token.LPAREN)
	list := []Expr{}
	for p.tok != token.EOF && p.tok != token.RPAREN {
		list = append(list, p.parseRHS())

		if p.tok == token.COMMA {
			p.next() // consume ,
		}
	}
	params := &ExprList{list}
	rparen := p.expect(token.RPAREN)

	return &CallExpr{Name: x, LParenPos: lparen, RParenPos: rparen, Params: params}
}

func (p *Parser) parseRHS() Expr {
	println("parseRHS")

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
	token, pos := p.scanner.nextLine()

	comment := &Comment{pos, token.Val}
	p.comments.Insert(comment)

	if debug {
		println("Trace: ")
		for _, cmt := range p.comments.comments {
			fmt.Println(cmt)
		}
	}
}

func (p *Parser) resolve(expr Expr) {
	println("resolve")
	id := expr.(*Ident)
	if id == nil {
		return
	}

	for s := p.scope; s != nil; s = s.outer {
		_, exist := s.Objects[id.Name]
		if exist {
			return
		}
	}

	p.UnResolved = append(p.UnResolved, id)
}

func (p *Parser) OpenScope() {
	p.scope = &Scope{p.scope, map[string]*Object{}}
}

func (p *Parser) CloseScope() {
	p.scope = p.scope.outer
}

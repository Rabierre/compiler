package main

import (
	"fmt"

	"github.com/rabierre/compiler/token"
)

const debug = true

type Parser struct {
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
}

func (p *Parser) Parse() {
	p.OpenScope()

	for !p.scanner.fullScaned {
		// TODO use p.peek() after move parsecomment phase to scanner
		tok, _ := p.scanner.peek()
		switch tok.Kind {
		case token.COMMENT:
			// TODO move to scanner
			p.parseComment()
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
	tok, _ := p.peek()
	switch tok.Kind {
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
	p.next()
	ident := p.parseIdent()

	p.expect(token.LPAREN)
	params := p.parseParamList()
	p.expect(token.RPAREN)

	// TODO parse func type
	var _typ token.Token
	if tok, _ := p.peek(); tok.Kind == token.INT || tok.Kind == token.DOUBLE {
		p.next()
		_typ = tok
	} else {
		_typ = token.Token{Kind: token.VOID}
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
	tok, pos := p.peek()
	if tok.Kind == token.IDENT {
		// TODO else error
		p.next()
	}

	return &Ident{Name: tok.Val, Pos: pos}
}

func (p *Parser) parseParamList() *StmtList {
	println("parseParamList")

	list := []Stmt{}
	for tok, _ := p.peek(); tok.Kind == token.INT || tok.Kind == token.DOUBLE; {
		list = append(list, p.parseParam())

		if tok, _ := p.peek(); tok.Kind != token.COMMENT && tok.Kind == token.RPAREN {
			break
		}
		p.next()
	}

	return &StmtList{List: list}
}

func (p *Parser) parseParam() Stmt {
	println("parseParam")
	tok, pos := p.next()
	param := &VarDeclStmt{Pos: pos, Type: tok}
	tok, pos = p.next()
	param.Name = &Ident{Pos: pos, Name: tok.Val}
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
	list := []Stmt{}
	for {
		tok, _ := p.peek()
		if tok.Kind == token.RBRACE || tok.Kind == token.EOF {
			break
		}
		list = append(list, p.parseStmt())
	}
	return list
}

func (p *Parser) parseStmt() Stmt {
	println("parseStmt")
	tok, pos := p.peek()

	switch tok.Kind {
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
		p.next()
	}
	// No statement? error
	return &BadStmt{From: pos}
}

// parse variable declaration
// int a = 1
// double b = 1.0
// int c
//
func (p *Parser) parseVarDecl() Stmt {
	typ, pos := p.next()

	ident := p.parseIdent()
	var value Expr
	if tok, _ := p.peek(); tok.Kind == token.ASSIGN {
		p.next()
		value = p.parseExpr(true)
	}

	decl := &VarDeclStmt{Pos: pos, Type: typ, Name: ident, RValue: value}

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
	x := p.parseExpr(true)
	tok, _ := p.peek()
	if tok.Kind == token.ASSIGN {
		p.next()
		y := p.parseExpr(true)
		x = &AssignExpr{Pos: x.(*Ident).Pos, LValue: x, RValue: y}
	}
	return &ExprStmt{expr: x}
}

func (p *Parser) parseForStmt() Stmt {
	_, pos := p.next()

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
	_, pos := p.next()

	p.expect(token.LPAREN)
	cond := p.parseExpr(true) // TODO: If cond is nil, error
	p.expect(token.RPAREN)

	body := p.parseCompoundStmt()

	var elseBody Stmt
	if tok, _ := p.peek(); tok.Kind == token.ELSE {
		p.next()
		elseBody = p.parseCompoundStmt()
	}
	return &IfStmt{Pos: pos, Cond: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseReturnStmt() Stmt {
	_, pos := p.next()

	var expr Expr
	if tok, _ := p.peek(); tok.Kind != token.EOF && tok.Kind != token.RBRACE {
		expr = p.parseExpr(true)
	}

	return &ReturnStmt{Pos: pos, Value: expr}
}

func (p *Parser) parseExprList() *ExprList {
	var exprs []Expr
	for {
		// 1. parse expr
		exprs = append(exprs, p.parseExpr(true))
		// 2. if next token is COMMA then continue
		// 2-1. else return
		if tok, _ := p.peek(); tok.Kind != token.COMMA {
			break
		}
		p.next()
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
		tok, _ := p.peek()
		// 1. if tok has high priority or equal than previous op
		// parse continously
		// 1-2. else return x
		if tok.Priority() < prio {
			return x
		}
		tok, pos := p.next()
		// 2. parse y as binay expr
		y := p.parseBinaryExpr(tok.Priority()+1, lookup)
		op := Operator{Val: tok}
		x = &BinaryExpr{Pos: pos, Op: op, LValue: x, RValue: y}
	}
}

// Factor ::= "(" Expr ")"
//         | AddSub Factor
//         | number
//         | string
func (p *Parser) parseUnaryExpr(lookup bool) Expr {
	println("parseUnaryExpr")

	tok, _ := p.peek()
	switch tok.Kind {
	case token.PLUS, token.MINUS:
		tok, pos := p.next()
		x := p.parseUnaryExpr(lookup)
		op := Operator{Val: tok}
		return &UnaryExpr{Pos: pos, Op: op, RValue: x}
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

	tok, _ := p.peek()
	switch tok.Kind {
	case token.LPAREN:
		if !lookup {
			p.resolve(x)
		}

		return p.parseCallExpr(x)
	case token.INC, token.DEC:
		if lookup {
			p.resolve(x)
		}

		tok, pos := p.next()
		op := Operator{Val: tok}
		// TODO shortExpr is similar with unaryExpr()
		return &ShortExpr{Pos: pos, Op: op, RValue: x}
	}
	return x
}

func (p *Parser) parseOperand(lookup bool) Expr {
	println("parseOperand")

	tok, pos := p.peek()
	switch tok.Kind {
	case token.IDENT:
		x := p.parseIdent()
		if lookup {
			p.resolve(x)
		}

		return x
	case token.INT_LIT, token.DOUBLE_LIT, token.TRUE, token.FALSE:
		p.next()
		return &BasicLit{Pos: pos, Value: tok.Val, Type: tok.Kind}
	}

	_, to := p.peek()
	return &BadExpr{From: pos, To: to}
}

func (p *Parser) parseCallExpr(x Expr) Expr {
	println("parseCallExpr")
	lparen := p.expect(token.LPAREN)
	list := []Expr{}
	for tok, _ := p.peek(); tok.Kind != token.EOF && tok.Kind != token.RPAREN; tok, _ = p.peek() {
		list = append(list, p.parseRHS())

		if tok, _ := p.peek(); tok.Kind == token.COMMA {
			p.next()
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

func (p *Parser) next() (token.Token, int) {
	for tok, _ := p.scanner.peek(); tok.Kind == token.COMMENT; tok, _ = p.scanner.peek() {
		p.scanner.nextLine()
	}
	return p.scanner.next()
}

func (p *Parser) peek() (token.Token, int) {
	for tok, _ := p.scanner.peek(); tok.Kind == token.COMMENT; tok, _ = p.scanner.peek() {
		p.scanner.nextLine()
	}
	return p.scanner.peek()
}

func (p *Parser) expect(expected token.Type) int {
	tok, pos := p.next()
	if tok.Kind != expected {
		panic("Expected: " + expected.String() + ", found: " + tok.Val + " " + tok.Kind.String())
	}
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

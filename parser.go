package main

import (
	"fmt"
)

const debug = true

type Parser struct {
	topScope *Scope
	scanner  *Scanner

	comments *CommentList

	// If we handle source codes in files
	// This should go in file struct
	//
	decls []Decl
}

func (p *Parser) Init(src []byte) {
	p.scanner = &Scanner{}
	p.scanner.Init()
	p.scanner.src = src
	p.comments = &CommentList{}
	p.OpenScope() // Top scope
}

func (p *Parser) Parse() {
	p.OpenScope()

	for !p.scanner.fullScaned {
		token, _ := p.scanner.peek()
		println("parse: ", token.val, token.kind.String())
		switch token.kind {
		case CommentType:
			// TODO move to scanner
			p.parseComment()
		default:
			p.parseDecl()
		}
	}

	p.CloseScope()
}

func (p *Parser) parseDecl() {
	tok, _ := p.scanner.next()
	println("parse decl: ", tok.val, tok.kind.String())
	switch tok.kind {
	// By spec for now, no global variable, no imports are available.
	// Function is top scope
	//
	case FuncType:
		p.parseFunc()
	}
}

func (p *Parser) parseFunc() {
	// 1. get function name
	ident := p.parseIdent()
	// 2. get parameters
	p.expect(LParenType)
	// TODO get param as Expr
	params := p.parseParamList()
	p.expect(RParenType)
	// 3. parse body
	body := p.parseBody() // parse compound statement
	// 4. make funcDecl
	funcDecl := &FuncDecl{Name: ident, Body: body, Params: params}
	// 5. add to decl
	p.decls = append(p.decls, funcDecl)
}

func (p *Parser) parseIdent() Ident {
	tok, pos := p.peek()
	if tok.kind == IdentType {
		p.next()
	}
	// TODO else error
	return Ident{Name: tok, Pos: pos}
}

func (p *Parser) parseParamList() *ArgList {
	list := []Arg{}

	for tok, _ := p.peek(); tok.kind == IntType || tok.kind == DoubleType; {
		list = append(list, p.parseParam())

		if tok, _ := p.peek(); tok.kind != CommaType && tok.kind == RParenType {
			break
		}
		p.next()
	}
	return &ArgList{List: list}
}

func (p *Parser) parseParam() Arg {
	tok, pos := p.next()
	arg := Arg{Pos: pos, Type: tok}
	tok, pos = p.next()
	arg.Name = Ident{Pos: pos, Name: tok}
	return arg
}

func (p *Parser) parseBody() *CompoundStmt {
	lbrace := p.expect(LBraceType)
	// TODO open scope
	list := p.parseStmtList()
	// TODO close scope
	rbrace := p.expect(RBraceType)

	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseCompoundStmt() *CompoundStmt {
	lbrace := p.expect(LBraceType)
	// TODO open scope
	list := p.parseStmtList()
	// TODO close scope
	rbrace := p.expect(RBraceType)
	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseStmtList() []Stmt {
	list := []Stmt{}
	for {
		token, _ := p.peek()
		println("parseStmtList: ", token.val, token.kind.String())
		if token.kind == RBraceType || token.kind == EOFType {
			break
		}
		list = append(list, p.parseStmt())
	}
	return list
}

func (p *Parser) parseStmt() Stmt {
	token, pos := p.peek()
	fmt.Printf("parseStmt: %s, %s\n", token.val, token.kind.String())

	switch token.kind {
	case IntType, DoubleType:
		return p.parseVarDecl()
	case ForType:
		return p.parseForStmt()
	case IfType:
		return p.parseIfStmt()
	case ReturnType:
		// "return" Expr ?
		return p.parseReturnStmt()
	case LBraceType:
		return p.parseCompoundStmt()
	case RBraceType:
		return &EmptyStmt{ /*position for semicolon if need*/ }
	default:
		p.next()
	}
	// No statement? error
	return &BadStmt{From: pos}
}

func (p *Parser) parseVarDecl() Stmt {
	typ, pos := p.next()
	fmt.Printf("parseVarDecl: %s, %s\n", typ.val, typ.kind.String())

	// int a = 1
	// double b = 1.0
	// int c
	ident := p.parseIdent()
	var value Expr
	tok, pos := p.peek()
	if tok.kind == AssignType {
		p.next()
		value = p.parseExpr()
	}

	// TODO resolve with scope

	return &VarDeclStmt{Pos: pos, Type: typ, Name: ident, RValue: value}
}

func (p *Parser) parseForStmt() Stmt {
	_, pos := p.next()

	// 1. get initial status
	p.expect(LParenType)
	// init := &EmptyStmt{}
	init := p.parseStmt()
	p.expect(SemiColType)
	// 2. get condition
	cond := p.parseExpr()
	p.expect(SemiColType)
	// 3. get post stmt
	post := p.parseExpr()
	p.expect(RParenType)

	// 4. parse body
	body := p.parseBody()

	// 5. make forDecl
	return &ForStmt{Pos: pos, Cond: cond, Init: init, Post: post, Body: body}
}

func (p *Parser) parseIfStmt() Stmt {
	_, pos := p.next()

	p.expect(LParenType)
	cond := p.parseExpr() // TODO: If cond is nil, error
	p.expect(RParenType)

	body := p.parseBody()

	var elseBody Stmt
	if tok, _ := p.peek(); tok.kind == ElseType {
		p.next()
		elseBody = p.parseBody()
	}
	return &IfStmt{Pos: pos, Cond: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseReturnStmt() Stmt {
	_, pos := p.next()

	var expr Expr
	if tok, _ := p.peek(); tok.kind != EOFType && tok.kind != RBraceType {
		expr = p.parseExpr()
	}

	return &ReturnStmt{Pos: pos, Value: expr}
}

func (p *Parser) parseExpr() Expr {
	return p.parseBinaryExpr(LowestPriority + 1)
}

func (p *Parser) parseExprList() *ExprList {
	var exprs []Expr
	for {
		// 1. parse expr
		exprs = append(exprs, p.parseExpr())
		// 2. if next token is COMMA then continue
		// 2-1. else return
		if tok, _ := p.peek(); tok.kind != CommaType {
			break
		}
		p.next()
	}
	return &ExprList{List: exprs}
}

// parse Term
func (p *Parser) parseBinaryExpr(prio int) Expr {
	// indetifier op expr
	x := p.parseUnaryExpr()
	for {
		tok, _ := p.peek()
		// 1. if tok has high priority or equal than previous op
		// parse continously
		// 1-2. else return x
		if tok.Priority() < prio {
			return x
		}
		op, pos := p.next()
		// 2. parse y as binay expr
		y := p.parseBinaryExpr(op.Priority() + 1)

		x = &BinaryExpr{Pos: pos, Op: op, LValue: x, RValue: y}
	}
}

// parse Factor
func (p *Parser) parseUnaryExpr() Expr {
	// Factor ::= "(" Expr ")"
	//         | AddSub Factor
	//         | number
	//         | string
	tok, _ := p.peek()
	switch tok.kind {
	case PlusType, MinusType:
		op, pos := p.next()
		x := p.parseUnaryExpr()
		return &UnaryExpr{Pos: pos, Op: op, RValue: x}
	}

	return p.parsePrimaryExpr()
}

func (p *Parser) parsePrimaryExpr() Expr {
	// identifier "(" ExprList ? ")"
	// identifier
	x := p.parseOperand()
	tok, _ := p.peek()
	switch tok.kind {
	case LParenType:
		lparen := p.expect(LParenType)
		params := p.parseExprList()
		rparen := p.expect(RParenType)
		return &CallExpr{Name: x, LParenPos: lparen, RParenPos: rparen, Params: params}
	}
	return x
}

func (p *Parser) parseOperand() Expr {
	tok, pos := p.peek()
	switch tok.kind {
	case IdentType:
		x := p.parseIdent()
		// TODO check x is declared in this scope
		return &x
	case IntLit, DoubleLit, TrueLit, FalseLit:
		p.next()
		return &BasicLit{Pos: pos, Value: tok.val, Type: tok.kind}
	}

	// TODO resolve ident is what type of identifier
	// function, variable,,

	_, to := p.peek()
	return &BadExpr{From: pos, To: to}
}

func (p *Parser) next() (Token, int) {
	var tok Token
	for tok, _ = p.scanner.peek(); tok.kind == CommentType; tok, _ = p.scanner.peek() {
		p.scanner.nextLine()
	}
	return p.scanner.next()
}

func (p *Parser) peek() (Token, int) {
	var tok Token
	for tok, _ = p.scanner.peek(); tok.kind == CommentType; tok, _ = p.scanner.peek() {
		p.scanner.nextLine()
	}
	return p.scanner.peek()
}

func (p *Parser) expect(expected TokenType) int {
	tok, pos := p.next()
	if tok.kind != expected {
		panic("Expected: " + expected.String() + ", found: " + tok.val + " " + tok.kind.String())
	}
	return pos
}

func (p *Parser) OpenScope() {
	p.topScope = &Scope{p.topScope, []*Object{}}
}

func (p *Parser) CloseScope() {
	p.topScope = p.topScope.outer
}

func (p *Parser) parseComment() {
	token, pos := p.scanner.nextLine()

	comment := &Comment{pos, token.val}
	p.comments.Insert(comment)

	if debug {
		println("Trace: ")
		for _, cmt := range p.comments.comments {
			fmt.Println(cmt)
		}
	}
}

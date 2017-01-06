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
		case COMMENT_SLASH:
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
	case FUNC:
		p.parseFunc()
	}
}

func (p *Parser) parseFunc() Decl {
	ident := p.parseIdent()

	p.expect(LPAREN)
	params := p.parseParamList()
	p.expect(RPAREN)

	var _typ Token
	if tok, _ := p.peek(); tok.kind == INT || tok.kind == DOUBLE {
		p.next()
		_typ = tok
	}

	body := p.parseBody()
	funcDecl := &FuncDecl{Name: ident, Body: body, Params: params, Type: _typ}
	p.decls = append(p.decls, funcDecl)
	return funcDecl
}

func (p *Parser) parseIdent() Ident {
	tok, pos := p.peek()
	if tok.kind == IDENT {
		p.next()
	}
	// TODO else error
	return Ident{Name: tok, Pos: pos}
}

func (p *Parser) parseParamList() *ArgList {
	list := []Arg{}

	for tok, _ := p.peek(); tok.kind == INT || tok.kind == DOUBLE; {
		list = append(list, p.parseParam())

		if tok, _ := p.peek(); tok.kind != COMMENT_SLASH && tok.kind == RPAREN {
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
	lbrace := p.expect(LBRACE)
	// TODO open scope
	list := p.parseStmtList()
	// TODO close scope
	rbrace := p.expect(RBRACE)

	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseCompoundStmt() *CompoundStmt {
	lbrace := p.expect(LBRACE)
	// TODO open scope
	list := p.parseStmtList()
	// TODO close scope
	rbrace := p.expect(RBRACE)
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
		if token.kind == RBRACE || token.kind == EOF {
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
	case INT, DOUBLE:
		return p.parseVarDecl()
	case FOR:
		return p.parseForStmt()
	case IF:
		return p.parseIfStmt()
	case RETURN:
		// "return" Expr ?
		return p.parseReturnStmt()
	case LBRACE:
		return p.parseCompoundStmt()
	case RBRACE:
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
	if tok.kind == ASSIGN {
		p.next()
		value = p.parseExpr()
	}

	// TODO resolve with scope

	return &VarDeclStmt{Pos: pos, Type: typ, Name: ident, RValue: value}
}

func (p *Parser) parseForStmt() Stmt {
	_, pos := p.next()

	// 1. get initial status
	p.expect(LPAREN)
	// init := &EmptyStmt{}
	init := p.parseStmt()
	p.expect(SEMI_COLON)
	// 2. get condition
	cond := p.parseExpr()
	p.expect(SEMI_COLON)
	// 3. get post stmt
	post := p.parseExpr()
	p.expect(RPAREN)

	// 4. parse body
	body := p.parseBody()

	// 5. make forDecl
	return &ForStmt{Pos: pos, Cond: cond, Init: init, Post: post, Body: body}
}

func (p *Parser) parseIfStmt() Stmt {
	_, pos := p.next()

	p.expect(LPAREN)
	cond := p.parseExpr() // TODO: If cond is nil, error
	p.expect(RPAREN)

	body := p.parseBody()

	var elseBody Stmt
	if tok, _ := p.peek(); tok.kind == ELSE {
		p.next()
		elseBody = p.parseBody()
	}
	return &IfStmt{Pos: pos, Cond: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseReturnStmt() Stmt {
	_, pos := p.next()

	var expr Expr
	if tok, _ := p.peek(); tok.kind != EOF && tok.kind != RBRACE {
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
		if tok, _ := p.peek(); tok.kind != COMMA {
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
	case PLUS, MINUS:
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
	case LPAREN:
		lparen := p.expect(LPAREN)
		params := p.parseExprList()
		rparen := p.expect(RPAREN)
		return &CallExpr{Name: x, LParenPos: lparen, RParenPos: rparen, Params: params}
	}
	return x
}

func (p *Parser) parseOperand() Expr {
	tok, pos := p.peek()
	switch tok.kind {
	case IDENT:
		x := p.parseIdent()
		// TODO check x is declared in this scope
		return &x
	case INT_LIT, DOUBLE_LIT, TRUE, FALSE:
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
	for tok, _ = p.scanner.peek(); tok.kind == COMMENT_SLASH; tok, _ = p.scanner.peek() {
		p.scanner.nextLine()
	}
	return p.scanner.next()
}

func (p *Parser) peek() (Token, int) {
	var tok Token
	for tok, _ = p.scanner.peek(); tok.kind == COMMENT_SLASH; tok, _ = p.scanner.peek() {
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

package main

import (
	"fmt"

	"github.com/rabierre/compiler/token"
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
}

func (p *Parser) parseDecl() {
	tok, _ := p.next()
	switch tok.Kind {
	// By spec for now, no global variable, no imports are available.
	// Function is top scope
	//
	case token.FUNC:
		p.parseFunc()
	}
}

func (p *Parser) parseFunc() Decl {
	ident := p.parseIdent()

	p.expect(token.LPAREN)
	params := p.parseParamList()
	p.expect(token.RPAREN)

	var _typ token.Token
	if tok, _ := p.peek(); tok.Kind == token.INT || tok.Kind == token.DOUBLE {
		p.next()
		_typ = tok
	} else {
		// TODO move to Type
		_typ = Token{kind: VOID}
	}

	body := p.parseBody()
	funcDecl := &FuncDecl{Name: ident, Body: body, Params: params, Type: _typ}

	// TODO move this to specific function like parse function decl only
	p.decls = append(p.decls, funcDecl)
	return funcDecl
}

func (p *Parser) parseIdent() Ident {
	tok, pos := p.peek()
	if tok.Kind == token.IDENT {
		p.next()
	}
	// TODO else error
	return Ident{Pos: pos, Name: tok.val}
}

func (p *Parser) parseParamList() *ArgList {
	list := []Arg{}

	for tok, _ := p.peek(); tok.Kind == token.INT || tok.Kind == token.DOUBLE; {
		list = append(list, p.parseParam())

		if tok, _ := p.peek(); tok.Kind != token.COMMENT && tok.Kind == token.RPAREN {
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
	arg.Name = Ident{Pos: pos, Name: tok.val}
	return arg
}

func (p *Parser) parseBody() *CompoundStmt {
	lbrace := p.expect(token.LBRACE)
	// TODO open scope
	list := p.parseStmtList()
	// TODO close scope
	rbrace := p.expect(token.RBRACE)

	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) parseCompoundStmt() *CompoundStmt {
	// TODO open scope
	lbrace := p.expect(token.LBRACE)
	list := p.parseStmtList()
	// TODO close scope
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
		value = p.parseExpr()
	}

	// TODO resolve with scope

	return &VarDeclStmt{Pos: pos, Type: typ, Name: ident, RValue: value}
}

func (p *Parser) parseExprStmt() Stmt {
	// a = 10
	x := p.parseExpr()
	tok, _ := p.peek()
	if tok.kind == ASSIGN {
		p.next()
		y := p.parseExpr()
		x = &AssignExpr{Pos: x.(*Ident).Pos, LValue: x, RValue: y}
	}
	return &ExprStmt{expr: x}
}

func (p *Parser) parseForStmt() Stmt {
	_, pos := p.next()

	// 1. get initial status
	p.expect(token.LPAREN)
	// init := &EmptyStmt{}
	init := p.parseStmt()
	p.expect(token.SEMI_COLON)
	// 2. get condition
	cond := p.parseExpr()
	p.expect(token.SEMI_COLON)
	// 3. get post stmt
	post := p.parseExpr()
	p.expect(token.RPAREN)

	// 4. parse body
	body := p.parseBody()

	// 5. make forDecl
	return &ForStmt{Pos: pos, Cond: cond, Init: init, Post: post, Body: body}
}

func (p *Parser) parseIfStmt() Stmt {
	_, pos := p.next()

	p.expect(token.LPAREN)
	cond := p.parseExpr() // TODO: If cond is nil, error
	p.expect(token.RPAREN)

	body := p.parseBody()

	var elseBody Stmt
	if tok, _ := p.peek(); tok.Kind == token.ELSE {
		p.next()
		elseBody = p.parseBody()
	}
	return &IfStmt{Pos: pos, Cond: cond, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseReturnStmt() Stmt {
	_, pos := p.next()

	var expr Expr
	if tok, _ := p.peek(); tok.Kind != token.EOF && tok.Kind != token.RBRACE {
		expr = p.parseExpr()
	}

	return &ReturnStmt{Pos: pos, Value: expr}
}

func (p *Parser) parseExpr() Expr {
	return p.parseBinaryExpr(token.LowestPriority + 1)
}

func (p *Parser) parseExprList() *ExprList {
	var exprs []Expr
	for {
		// 1. parse expr
		exprs = append(exprs, p.parseExpr())
		// 2. if next token is COMMA then continue
		// 2-1. else return
		if tok, _ := p.peek(); tok.Kind != token.COMMA {
			break
		}
		p.next()
	}
	return &ExprList{List: exprs}
}

// Term
func (p *Parser) parseBinaryExpr(prio int) Expr {
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

// Factor
func (p *Parser) parseUnaryExpr() Expr {
	// Factor ::= "(" Expr ")"
	//         | AddSub Factor
	//         | number
	//         | string
	tok, _ := p.peek()
	switch tok.Kind {
	case token.PLUS, token.MINUS:
		op, pos := p.next()
		x := p.parseUnaryExpr()
		return &UnaryExpr{Pos: pos, Op: op, RValue: x}
	}

	return p.parsePrimaryExpr()
}

// Parse function call or variable reference
// identifier "(" ExprList ? ")"
// identifier
//
func (p *Parser) parsePrimaryExpr() Expr {
	x := p.parseOperand()
	tok, _ := p.peek()
	switch tok.Kind {
	case token.LPAREN:
		lparen := p.expect(token.LPAREN)
		params := p.parseExprList()
		rparen := p.expect(token.RPAREN)
		return &CallExpr{Name: x, LParenPos: lparen, RParenPos: rparen, Params: params}
	}
	return x
}

func (p *Parser) parseOperand() Expr {
	tok, pos := p.peek()
	switch tok.Kind {
	case token.IDENT:
		x := p.parseIdent()
		// TODO Resolve identity in scope
		return &x
	case token.INT_LIT, token.DOUBLE_LIT, token.TRUE, token.FALSE:
		p.next()
		return &BasicLit{Pos: pos, Value: tok.Val, Type: tok.Kind}
	}

	_, to := p.peek()
	return &BadExpr{From: pos, To: to}
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

func (p *Parser) OpenScope() {
	p.topScope = &Scope{p.topScope, []*Object{}}
}

func (p *Parser) CloseScope() {
	p.topScope = p.topScope.outer
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

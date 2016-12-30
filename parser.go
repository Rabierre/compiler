package main

import (
	"fmt"
)

const debug = true

type Parser struct {
	topScope *Scope // currently don't support multi level scope. eg can't use function in function
	scanner  *Scanner

	comments *CommentList

	// If we handle source codes in files
	// This should go in file struct
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
	tok, pos := p.next()
	if tok.kind != IdentType {
		panic("Not function identifier: " + tok.val + " " + tok.kind.String())
	}

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
		// expression
	case ForType:
		return p.parseForStmt()
	case IfType:
		// IfStmt
	case ReturnType:
		// "return" Expr ?
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

func (p *Parser) parseForStmt() Stmt {
	_, pos := p.next()
	// 1. get initial status
	p.expect(LParenType)
	init := &EmptyStmt{}
	p.expect(SemiColType)
	// 2. get condition
	cond := &EmptyStmt{}
	p.expect(SemiColType)
	// 3. get post stmt
	post := &EmptyStmt{}
	p.expect(RParenType)
	// 4. parse body
	body := p.parseBody() // parse compound statement
	// 5. make forDecl
	return &ForStmt{Pos: pos, Cond: cond, Init: init, Post: post, Body: body}
}

func (p *Parser) next() (Token, int) {
	tok, _ := p.scanner.peek()
	// TODO skip all comment
	if tok.kind == CommentType {
		p.scanner.nextLine()
	}
	return p.scanner.next()
}

func (p *Parser) peek() (Token, int) {
	tok, _ := p.scanner.peek()
	// TODO skip all comment
	if tok.kind == CommentType {
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

// TODO move to scanner
func (p *Parser) parseComment() {
	println("parsecomment")
	token, pos := p.scanner.nextLine()

	// TODO Use integer position in source code not token index
	comment := &Comment{pos, token.val}
	p.comments.Insert(comment)

	if debug {
		println("Trace: ")
		for _, cmt := range p.comments.comments {
			fmt.Println(cmt)
		}
	}
}

// // What to do with tokens? where should we save them?
// func ForStatement(forToken Token) Block {
// 	block := Block{[]Filed{}, ForBlock}
// 	// "for"
// 	block.tokens = append(block.tokens, forToken)

// 	token := scanner.next()
// 	// "(" Expr ";" OptExpr ";" OptExpr ")"
// 	if token.kind == LParenType {
// 		// TODO initial statement
// 		Expression(block)
// 		scanner.next() // drop ";"
// 		// TODO condition statement
// 		Expression(block)
// 		scanner.next() // drop ";"
// 		// TODO increase statement
// 		Expression(block)
// 		scanner.next() // drop ")"
// 	}
// 	// CompoundStmt
// 	CompoundStatement(block)

// 	// // TODO add block to somewhere
// 	return block
// }

// func Expression(block Block) {
// 	token := scanner.next()
// 	if token.kind == IdentType {
// 		// identifier "=" Expr
// 	} else {
// 		// Rvalue
// 	}
// }

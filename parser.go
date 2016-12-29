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

func (p *Parser) next() (Token, int) {
	tok, _ := p.scanner.peek()
	if tok.kind == CommentType {
		// Consume comment list
		p.scanner.next()
		p.parseComment()
	}
	return p.scanner.next()
}

func (p *Parser) Parse() {
	p.OpenScope()

	for {
		token, _ := p.next()

		// println("Parsed token: ", token.kind.String())
		switch token.kind {
		case CommentType:
			p.parseComment()
		// TODO seperate comment consume and decl parse
		case FuncType:
			p.parseFunc()
		}
		if p.scanner.fullScaned {
			break
		}
	}

	p.CloseScope()
}

func (p *Parser) OpenScope() {
	p.topScope = &Scope{p.topScope, []*Object{}}
}

func (p *Parser) CloseScope() {
	p.topScope = p.topScope.outer
}

func (p *Parser) parseComment() {
	// '//' already consumed
	token := p.scanner.nextLine() // TODO get line not token
	println("Parse Comment: ", token.val)

	// TODO Use integer position in source code not token index
	comment := &Comment{p.scanner.srcIndex, token.val}
	p.comments.Insert(comment)

	if debug {
		println("Trace: ")
		for _, cmt := range p.comments.comments {
			fmt.Println(cmt)
		}
	}
}

type Ident struct {
	Pos  int
	Name string
}

type Field struct {
}

type ArgList struct {
	fields []*Field
}

type Stmt interface {
	// IfStmt
	// ForStmt
	// Expr
	// CompoundStmt
	// "return" Expr ?
	stmtNode()
}

type CompoundStmt struct {
	LBracePos int
	RBracePos int
	List      []Stmt
}

type EmptyStmt struct {
}

type BadStmt struct {
	From int
}

func (*CompoundStmt) stmtNode() {}
func (*EmptyStmt) stmtNode()    {}
func (*BadStmt) stmtNode()      {}

type FuncDecl struct {
	Name   Ident
	Body   *CompoundStmt // body or nil
	Params *ArgList      // list of parameters
}

func (*FuncDecl) declNode() {}

func (p *Parser) parseFunc() {
	// 1. get function name
	ident := p.parseIdent()
	// 2. get parameters
	// TODO
	// 3. parse body
	body := p.parseBody() // parse compound statement
	// 4. make funcDecl
	funcDecl := &FuncDecl{Name: ident, Body: body}
	// 5. add to decl
	p.decls = append(p.decls, funcDecl)
}

func (p *Parser) parseIdent() Ident {
	tok, pos := p.next()
	if tok.kind != IdentType {
		panic("Not function identifier: " + tok.val + " " + tok.kind.String())
	}

	for i := 0; i < len(tok.val); i++ {
		if rune(tok.val[i]) == '(' {
			tok.val = tok.val[:i]
		}
	}

	return Ident{Name: tok.val, Pos: pos}
}

func (p *Parser) parseBody() *CompoundStmt {
	lbrace := p.expect(LBraceType)
	// TODO open scope
	println("parseBody")
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
	println("parseBody")
	list := p.parseStmtList()
	// TODO close scope
	rbrace := p.expect(RBraceType)
	return &CompoundStmt{
		LBracePos: lbrace,
		RBracePos: rbrace,
		List:      list,
	}
}

func (p *Parser) expect(expected TokenType) int {
	tok, pos := p.next()
	if tok.kind != expected {
		panic("Expected: " + expected.String() + ", found: " + tok.val + " " + tok.kind.String())
	}
	return pos
}

func (p *Parser) parseStmtList() []Stmt {
	list := []Stmt{}
	for {
		token, _ := p.scanner.peek()
		if token.kind == RBraceType || token.kind == EOFType {
			break
		}
		list = append(list, p.parseStmt())
	}
	return list
}

func (p *Parser) parseStmt() Stmt {
	token, pos := p.scanner.peek()
	println("parseStmt: ", token.val)

	switch token.kind {
	case IntType, DoubleType:
		// expression
	case ForType:
		// ForStmt
	case IfType:
		// IfStmt
	case ReturnType:
		// "return" Expr ?
	case LBraceType:
		return p.parseCompoundStmt()
	case RBraceType:
		return &EmptyStmt{ /*position for semicolon if need*/ }
	}
	// No statement? error
	return &BadStmt{From: pos}
}

// TODO parse decl(function, const ..) here
func (p *Parser) parseDecl() {
}

// ------------------------------------------------------------------
// TODO go somewhere
type CommentList struct {
	comments []*Comment
}

func (c *CommentList) Insert(comment *Comment) {
	c.comments = append(c.comments, comment)
}

type Comment struct {
	pos  int // TODO position of comment slash' in source code
	text string
}

type Decl interface {
	declNode()
}

type Scope struct {
	outer   *Scope
	Objects []*Object // better contain name of it for convenient search decl in this scope
}

func (s *Scope) Insert(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

type Object struct {
	// Kind Type
	decl interface{} // statement(function, for, if..), expression
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

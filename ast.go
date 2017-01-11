package main

import (
	"github.com/rabierre/compiler/token"
)

type Node interface {
}

//--------------------------------------------------------------------------------------
// Expression
//
type Expr interface {
	Node
	exprNode()
}

type ExprList struct {
	List []Expr
}

type BasicLit struct {
	Pos   int
	Value string
	Type  token.Type
}

// Term
type BinaryExpr struct {
	Pos    int
	LValue Expr
	RValue Expr
	Op     token.Token
}

// Factor
type UnaryExpr struct {
	Pos    int
	Op     token.Token
	RValue Expr
}

type Ident struct {
	Pos  int
	Name token.Token
}

type CallExpr struct {
	Name      Expr
	Params    *ExprList
	LParenPos int
	RParenPos int
}

type Arg struct {
	Pos  int
	Type token.Token
	Name Ident
}

type ArgList struct {
	List []Arg
}

type BadExpr struct {
	From int
	To   int
}

func (*BasicLit) exprNode()   {}
func (*Ident) exprNode()      {}
func (*BinaryExpr) exprNode() {}
func (*UnaryExpr) exprNode()  {}
func (*CallExpr) exprNode()   {}
func (*Arg) exprNode()        {}
func (*BadExpr) exprNode()    {}

//--------------------------------------------------------------------------------------
// Declaration
//
type Decl interface {
	Node
	declNode()
}

type FuncDecl struct {
	// TODO Pos
	Name   Ident
	Type   token.Token
	Params *ArgList
	Body   *CompoundStmt
}

func (*FuncDecl) declNode()    {}
func (*VarDeclStmt) declNode() {}

//--------------------------------------------------------------------------------------
// Statement
//
type Stmt interface {
	// IfStmt, ForStmt, Expr, CompoundStmt, ReturnStmt
	Node
	stmtNode()
}

type CompoundStmt struct {
	LBracePos int
	RBracePos int
	List      []Stmt
}

type IfStmt struct {
	Pos      int
	Cond     Expr // Assign is not available
	Body     *CompoundStmt
	ElseBody Stmt
}

type ForStmt struct {
	Pos  int
	Init Stmt
	Cond Expr
	Post Expr
	Body *CompoundStmt
}

type VarDeclStmt struct {
	Pos    int
	Type   token.Token
	Name   Ident
	RValue Expr
}

type ReturnStmt struct {
	Pos   int
	Value Expr
}

type EmptyStmt struct {
}

type BadStmt struct {
	From int
}

func (*CompoundStmt) stmtNode() {}
func (*ForStmt) stmtNode()      {}
func (*IfStmt) stmtNode()       {}
func (*VarDeclStmt) stmtNode()  {}
func (*ReturnStmt) stmtNode()   {}
func (*EmptyStmt) stmtNode()    {}
func (*BadStmt) stmtNode()      {}

//--------------------------------------------------------------------------------------
// Comment
//
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

//--------------------------------------------------------------------------------------
// Scope
//
type Scope struct {
	outer   *Scope
	Objects []*Object // better contain name of it for convenient when resolving
}

func (s *Scope) Insert(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

type ObjectType int

const (
	FUNC = iota
	VAR
)

type Object struct {
	kind ObjectType
	decl interface{}
}

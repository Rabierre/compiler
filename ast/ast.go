package ast

import (
	"github.com/rabierre/compiler/token"
)

type Node interface {
}

type Operator struct {
	Node
	Type token.Type
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
	Op     Operator
}

// Factor
type UnaryExpr struct {
	Pos    int
	RValue Expr
	Op     Operator
}

type ShortExpr struct {
	Pos    int
	Op     Operator
	RValue Expr
}

type Ident struct {
	Pos  int
	Name string
}

type CallExpr struct {
	Name      Expr
	Params    *ExprList
	LParenPos int
	RParenPos int
}

type AssignExpr struct {
	Pos    int
	LValue Expr
	RValue Expr
}

type BadExpr struct {
	From int
	To   int
}

func (*BasicLit) exprNode()   {}
func (*Ident) exprNode()      {}
func (*BinaryExpr) exprNode() {}
func (*UnaryExpr) exprNode()  {}
func (*ShortExpr) exprNode()  {}
func (*CallExpr) exprNode()   {}
func (*AssignExpr) exprNode() {}
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
	Name   *Ident
	Type   token.Type
	Params *StmtList
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

type StmtList struct {
	List []Stmt
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
	Type   token.Type
	Name   *Ident
	RValue Expr
}

type ReturnStmt struct {
	Pos   int
	Value Expr
}

type ExprStmt struct {
	Val Expr
}

type EmptyStmt struct {
}

type BadStmt struct {
	From int
}

func (*StmtList) stmtNode()     {}
func (*CompoundStmt) stmtNode() {}
func (*ForStmt) stmtNode()      {}
func (*IfStmt) stmtNode()       {}
func (*VarDeclStmt) stmtNode()  {}
func (*ReturnStmt) stmtNode()   {}
func (*ExprStmt) stmtNode()     {}
func (*EmptyStmt) stmtNode()    {}
func (*BadStmt) stmtNode()      {}

//--------------------------------------------------------------------------------------
// Comment
//
type CommentList struct {
	List []*Comment
}

func (c *CommentList) Insert(comment *Comment) {
	c.List = append(c.List, comment)
}

type Comment struct {
	Pos  int
	Text string
}

//--------------------------------------------------------------------------------------
// Scope
//
type Scope struct {
	Outer   *Scope
	Objects map[string]*Object
}

func (s *Scope) Insert(obj *Object, name string) {
	_, exist := s.Objects[name]
	if exist {
		panic("Already exist: " + name)
	}
	s.Objects[name] = obj
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

func NewObject(decl interface{}, kind ObjectType) *Object {
	return &Object{decl: decl, kind: kind}
}

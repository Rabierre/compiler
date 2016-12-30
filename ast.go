package main

type Ident struct {
	Pos  int
	Name Token
}

type Node interface {
}

//--------------------------------------------------------------------------------------
// Expression
type Expr interface {
	Node
	exprNode()
}

type BinaryExpr struct {
	Pos      int
	LValue   Expr
	RValue   Expr
	Operator Token
}

type Arg struct {
	Pos  int
	Type Token // int, double
	Name Ident
}

func (*BinaryExpr) exprNode() {}
func (*Arg) exprNode()        {}

type ArgList struct {
	List []Arg
}

//--------------------------------------------------------------------------------------
// Declaration
//
type Decl interface {
	Node
	declNode()
}

type FuncDecl struct {
	Name   Ident
	Body   *CompoundStmt // body or nil
	Params *ArgList      // list of parameters
}

func (*FuncDecl) declNode() {}

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
	Pos  int
	Cond Expr // Assign is not available
	Body *CompoundStmt
}

type ForStmt struct {
	Pos  int
	Init Stmt
	Cond Stmt // TODO expr
	Post Stmt
	Body *CompoundStmt
}

type EmptyStmt struct {
}

type BadStmt struct {
	From int
}

func (*CompoundStmt) stmtNode() {}
func (*ForStmt) stmtNode()      {}
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
	Objects []*Object // better contain name of it for convenient search decl in this scope
}

func (s *Scope) Insert(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

type Object struct {
	// Kind Type
	decl interface{} // statement(function, for, if..), expression
}

package main

type Node interface {
}

//--------------------------------------------------------------------------------------
// Expression
type Expr interface {
	Node
	exprNode()
}

type BasicLit struct {
	Pos   int
	Value string
	Type  TokenType
}

// Term
type BinaryExpr struct {
	Pos    int
	LValue Expr
	RValue Expr
	Op     Token
}

// Factor
type UnaryExpr struct {
	Pos    int
	Op     Token
	RValue Expr
}

type Ident struct {
	Pos  int
	Name Token
}

type CallExpr struct {
	Name      Expr
	Params    *ExprList
	LParenPos int
	RParenPos int
}

type ExprList struct {
	List []Expr
}

type Arg struct {
	Pos  int
	Type Token // int, double
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
	Cond Expr
	Post Expr
	Body *CompoundStmt
}

// TODO also value
type VarDeclStmt struct {
	Pos    int
	Type   Token
	Name   Ident
	RValue Expr
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

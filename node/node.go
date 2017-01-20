package node

import (
	"fmt"

	"github.com/rabierre/compiler/ast"
)

type Node struct {
	left  *Node
	right *Node
	list  []*Node
	Rlist *Node // multi return not allowed

	Type *Type // Type not Node
	Name *Name

	Func *Func

	Sym *Sym
	E   interface{} // Value

	Op Op
}

// object's declaration in file, package
type Sym struct {
	Name string
	Def  *Node // definition of sym

	// imported package
}

type Type struct {
}

type Func struct {
	NBody  []*Node
	NParam *Node
}

type Name struct {
	Val string

	Op Op
	// TODO
	// symbol
	// link: declared position
}

type Op int

const (
	OIDENT = iota
	ODCLFUNC
	ODCLVAR
	ONAME
	OTFUNC

	OLIT

	OASSGN

	// expr
	OCALL

	OEMPTY
)

type Value struct {
	// Pint, Pdouble
	E interface{} // primitive types in c
}

type Pint struct {
	// Val
	// overflow
}

func funcDecl(decl *ast.FuncDecl) *Node {
	// Header
	f := newNode(ODCLFUNC)
	f.Name = name(decl.Name)
	// TODO collect all varExpr and sort

	t := newNode(OTFUNC) // Type of function
	t.list = params(decl.Params)
	t.Rlist = newNode(OLIT)
	t.Rlist.E = Value{E: Pint{}}
	// t.Rlist.Type = TIDEAL // TODO

	// TODOf.Func.Nname.Name.Param.Ntype = t // TODO: check if nname already has an ntyp
	f.Func.NParam = t

	// Body
	var b []*Node
	if decl.Body != nil {
		b = stmts(decl.Body.List)
	}
	if len(b) == 0 { // TODO change to compare nil
		b = []*Node{newNode(OEMPTY)}
	}
	f.Func.NBody = b

	return f
}

func param(decl ast.Stmt) *Node {
	p := newNode(ODCLVAR)
	p.left = nameNode(name(decl.(*ast.VarDeclStmt).Name))

	n := newNode(OLIT) // TODO custom type later
	n.Type = new(Type) // TODO filed of Type
	// TODO set n.Type's type as TIDEAL
	p.right = n

	return p
}

func params(stmts *ast.StmtList) []*Node {
	var list []*Node
    if stmts == nil {
        return list
    }

	for _, s := range stmts.List {
		list = append(list, param(s))
	}
	return list
}

func nameNode(name *Name) *Node {
	n := newNode(ONAME)
	n.Sym = &Sym{Def: n}
	return n
}

func name(name *ast.Ident) *Name {
	return &Name{Val: name.Name}
}

func stmt(stmt ast.Stmt) *Node {
	switch stmt := stmt.(type) {
	case *ast.ForStmt:
		println(stmt)
	case *ast.CompoundStmt:
	case *ast.IfStmt:
	case *ast.VarDeclStmt:
	case *ast.ReturnStmt:
	case *ast.ExprStmt:
		// TODO return expr(stmt.Val)
	case *ast.EmptyStmt:
		// TODO insert sentinel
	case *ast.BadStmt:
	}
	panic(fmt.Sprintf("unhandled stmt: %+v", stmt))
}

func stmts(stmts []ast.Stmt) (l []*Node) {
	for _, p := range stmts {
		l = append(l, stmt(p))
	}
	return
}

func decls(decls []ast.Decl) []*Node {
	var list []*Node
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			list = append(list, funcDecl(decl))
		}
	}

	return list
}

func newNode(op Op) *Node {
	n := new(Node)
	n.Op = op
	if op == ODCLFUNC {
		n.Func = new(Func)
	}
	return n
}

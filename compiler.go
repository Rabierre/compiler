package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rabierre/compiler/token"
)

type Compiler struct {
	buf bytes.Buffer

	input  string
	output string
}

func (c *Compiler) Init(input, output string) {
	c.input = input
	c.output = output
}

func (c *Compiler) Compile(src []byte) {
	parser := Parser{}
	parser.Init(src)
	parser.Parse()

	// function is top scope
	for _, decl := range parser.decls {
		fn := decl.(*FuncDecl)
		c.emitType(fn.Type)
		c.buf.WriteByte(' ')
		c.buf.WriteString(fn.Name.Name)

		c.emitParams(fn.Params)
		c.emitBody(fn.Body)
	}

	ioutil.WriteFile(path.Join(c.output, "mid.c"), c.buf.Bytes(), 0777)
}

func (c *Compiler) emitBody( /*Don't handle ast directly*/ stnt Stmt) {
	c.emitCompoundStmt(stnt)
}

func (c *Compiler) emitCompoundStmt( /*Don't handle ast directly*/ stmt Stmt) {
	cs := stmt.(*CompoundStmt)
	c.buf.WriteString("{\n")
	for _, s := range cs.List {
		c.emitStmt(s)
	}
	c.buf.WriteString("}\n")
}

func (c *Compiler) emitIfStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*IfStmt)
	c.buf.WriteString("if (")
	c.emitExpr(s.Cond)
	c.buf.WriteString(")\n")
	c.emitBody(s.Body)
	c.buf.WriteString("else\n")
	c.emitStmt(s.ElseBody)
}

func (c *Compiler) emitForStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*ForStmt)
	c.buf.WriteString("for (")
	c.emitStmt(s.Init) // TODO emitShortVarDecl
	c.buf.WriteRune(';')
	c.emitExpr(s.Cond)
	c.buf.WriteRune(';')
	c.emitExpr(s.Post)
	c.buf.WriteRune(')')
	c.emitBody(s.Body)
}

func (c *Compiler) emitReturnStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*ReturnStmt)
	c.buf.WriteString("return ")
	c.emitExpr(s.Value)
	c.buf.WriteString(";\n")
}

func (c *Compiler) emitVarDeclStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*VarDeclStmt)
	c.buf.WriteString(s.Type.Val)
	c.buf.WriteRune(' ')
	c.buf.WriteString(s.Name.Name)
	if s.RValue != nil {
		c.buf.WriteRune('=')
		c.emitExpr(s.RValue)
	}
	c.buf.WriteRune(';')
	c.buf.WriteRune('\n')
}

func (c *Compiler) emitExprStmt( /*Don't handle ast directly*/ stmt Stmt) {
	c.emitExpr(stmt.(*ExprStmt).expr)
}

func (c *Compiler) emitStmt( /*Don't handle ast directly*/ stmt Stmt) {
	switch typ := stmt.(type) {
	case (*CompoundStmt):
		c.emitCompoundStmt(stmt)
	case (*IfStmt):
		c.emitIfStmt(stmt)
	case (*ForStmt):
		c.emitForStmt(stmt)
	case (*ReturnStmt):
		c.emitReturnStmt(stmt)
	case (*VarDeclStmt):
		c.emitVarDeclStmt(stmt)
	case (*ExprStmt):
		c.emitExprStmt(stmt)
	default:
		println("Type: ", typ)
	}
}

func (c *Compiler) emitLiteracy( /*Don't handle ast directly*/ expr Expr) {
	ex := expr.(*BasicLit)
	// TODO is type need?
	// buf.WriteString(ex.Type.String())
	c.buf.WriteString(ex.Value)
}

func (c *Compiler) emitBinaryExpr( /*Don't handle ast directly*/ expr Expr) {
	bin := expr.(*BinaryExpr)
	c.emitExpr(bin.LValue)
	c.buf.WriteString(bin.Op.Val)
	c.emitExpr(bin.RValue)
}

func (c *Compiler) emitIdent( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*Ident)
	c.buf.WriteString(e.Name)
}

func (c *Compiler) emitCallExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*CallExpr)
	c.emitExpr(e.Name)
	c.buf.WriteRune('(')
	c.emitExprList(e.Params.List)
	c.buf.WriteRune(')')
}

func (c *Compiler) emitAssignExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*AssignExpr)
	c.emitExpr(e.LValue)
	c.buf.WriteRune('=')
	c.emitExpr(e.RValue)
	c.buf.WriteRune(';')
	c.buf.WriteRune('\n')
}

func (c *Compiler) emitExprList(list []Expr) {
	for _, e := range list {
		c.emitExpr(e)
	}
}

// TODO maybe function chain
func (c *Compiler) emitExpr( /*Don't handle ast directly*/ expr Expr) {
	switch typ := expr.(type) {
	case (*BasicLit):
		c.emitLiteracy(expr)
	case (*BinaryExpr):
		c.emitBinaryExpr(expr)
	case (*Ident):
		c.emitIdent(expr)
	case (*CallExpr):
		c.emitCallExpr(expr)
	case (*AssignExpr):
		c.emitAssignExpr(expr)
	default:
		fmt.Println("Type: ", typ)
	}
}

func (c *Compiler) emitType( /*Don't handle ast directly*/ typ token.Token) {
	c.buf.WriteString(typ.Kind.String())
}

func (c *Compiler) emitParams( /*Don't handle ast directly*/ params *StmtList) {
	c.buf.WriteString("(")
	for i, p := range params.List {
		d := p.(*VarDeclStmt)
		c.buf.WriteString(fmt.Sprintf("%s %s", d.Type.Kind.String(), d.Name.Name))
		if i < len(params.List)-1 {
			c.buf.WriteString(", ")
		}
	}
	c.buf.WriteString(")\n")
}

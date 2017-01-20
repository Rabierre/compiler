package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rabierre/compiler/ast"
	"github.com/rabierre/compiler/token"
)

type Compiler struct {
	buf bytes.Buffer

	input  string
	output string

	tlevel int
}

func (c *Compiler) Init(input, output string) {
	c.input = input
	c.output = output
}

func (c *Compiler) Compile(src []byte) {
	parser := Parser{}
	parser.Init(src)
	parser.Parse()

	for _, decl := range parser.decls {
		fn := decl.(*ast.FuncDecl)
		c.emitType(fn.Type)
		c.buf.WriteByte(' ')
		c.buf.WriteString(fn.Name.Name)

		c.emitParamTypes(fn.Params)
		c.buf.WriteByte(';')
		c.buf.WriteByte('\n')
	}

	ioutil.WriteFile(path.Join(c.output, "mid.h"), c.buf.Bytes(), 0777)
	c.buf.Truncate(c.buf.Len())

	// function is top scope
	for _, decl := range parser.decls {
		fn := decl.(*ast.FuncDecl)
		c.emitType(fn.Type)
		c.buf.WriteByte(' ')
		c.buf.WriteString(fn.Name.Name)

		c.emitParams(fn.Params)
		c.buf.WriteByte('\n')
		c.emitBody(fn.Body)
	}

	ioutil.WriteFile(path.Join(c.output, "mid.c"), c.buf.Bytes(), 0777)
}

func (c *Compiler) emitBody( /*Don't handle ast directly*/ stnt ast.Stmt) {
	c.emitCompoundStmt(stnt)

}

func (c *Compiler) emitCompoundStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	c.write("{\n")
	c.tlevel++
	cs := stmt.(*ast.CompoundStmt)
	for _, s := range cs.List {
		c.emitStmt(s)
	}
	c.tlevel--
	c.write("}\n")
}

func (c *Compiler) emitIfStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	s := stmt.(*ast.IfStmt)
	c.write("if (")
	c.emitExpr(s.Cond)
	c.buf.WriteString(")\n")
	c.emitBody(s.Body)
	c.write("else\n")
	c.emitBody(s.ElseBody)
}

func (c *Compiler) emitForStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	s := stmt.(*ast.ForStmt)
	c.write("for (")
	c.emitShortDeclStmt(s.Init) // TODO emitShortVarDecl
	c.buf.WriteRune(';')
	c.emitExpr(s.Cond)
	c.buf.WriteRune(';')
	c.emitExpr(s.Post)
	c.buf.WriteString(")\n")
	c.emitBody(s.Body)
}

func (c *Compiler) emitReturnStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	s := stmt.(*ast.ReturnStmt)
	c.write("return ")
	c.emitExpr(s.Value)
	c.buf.WriteString(";\n")
}

// for for stmt
func (c *Compiler) emitShortDeclStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	s := stmt.(*ast.VarDeclStmt)
	c.buf.WriteString(s.Type.String())
	c.buf.WriteRune(' ')
	c.buf.WriteString(s.Name.Name)
	if s.RValue != nil {
		c.buf.WriteRune('=')
		c.emitExpr(s.RValue)
	}
}

func (c *Compiler) emitVarDeclStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	s := stmt.(*ast.VarDeclStmt)
	c.write("")
	c.buf.WriteString(s.Type.String())
	c.buf.WriteRune(' ')
	c.buf.WriteString(s.Name.Name)
	if s.RValue != nil {
		c.buf.WriteRune('=')
		c.emitExpr(s.RValue)
	}
	c.buf.WriteRune(';')
	c.buf.WriteRune('\n')
}

func (c *Compiler) emitExprStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	c.write("")
	c.emitExpr(stmt.(*ast.ExprStmt).Val)
	c.buf.WriteByte(';')
	c.buf.WriteByte('\n')
}

func (c *Compiler) emitStmt( /*Don't handle ast directly*/ stmt ast.Stmt) {
	switch typ := stmt.(type) {
	case (*ast.CompoundStmt):
		c.emitCompoundStmt(stmt)
	case (*ast.IfStmt):
		c.emitIfStmt(stmt)
	case (*ast.ForStmt):
		c.emitForStmt(stmt)
	case (*ast.ReturnStmt):
		c.emitReturnStmt(stmt)
	case (*ast.VarDeclStmt):
		c.emitVarDeclStmt(stmt)
	case (*ast.ExprStmt):
		c.emitExprStmt(stmt)
	default:
		println("Type: ", typ)
	}
}

func (c *Compiler) emitLiteracy( /*Don't handle ast directly*/ expr ast.Expr) {
	ex := expr.(*ast.BasicLit)
	// TODO is type need?
	// buf.WriteString(ex.Type.String())
	c.buf.WriteString(ex.Value)
}

func (c *Compiler) emitBinaryExpr( /*Don't handle ast directly*/ expr ast.Expr) {
	e := expr.(*ast.BinaryExpr)
	c.emitExpr(e.LValue)
	c.buf.WriteString(e.Op.Type.String())
	c.emitExpr(e.RValue)
}

func (c *Compiler) emitIdent( /*Don't handle ast directly*/ expr ast.Expr) {
	e := expr.(*ast.Ident)
	c.buf.WriteString(e.Name)
}

func (c *Compiler) emitCallExpr( /*Don't handle ast directly*/ expr ast.Expr) {
	e := expr.(*ast.CallExpr)
	c.emitExpr(e.Name)
	c.buf.WriteRune('(')

	list := e.Params.List
	for i := 0; i < len(list); i++ {
		c.emitExpr(list[i])
		if i < len(list)-1 {
			c.buf.WriteString(", ")
		}
	}

	c.buf.WriteRune(')')
}

func (c *Compiler) emitAssignExpr( /*Don't handle ast directly*/ expr ast.Expr) {
	e := expr.(*ast.AssignExpr)
	c.emitExpr(e.LValue)
	c.buf.WriteRune('=')
	c.emitExpr(e.RValue)
}

func (c *Compiler) emitShortExpr( /*Don't handle ast directly*/ expr ast.Expr) {
	e := expr.(*ast.ShortExpr)
	c.emitExpr(e.RValue)
	c.buf.WriteString(e.Op.Type.String())
}

// TODO maybe function chain
func (c *Compiler) emitExpr( /*Don't handle ast directly*/ expr ast.Expr) {
	switch typ := expr.(type) {
	case (*ast.BasicLit):
		c.emitLiteracy(expr)
	case (*ast.BinaryExpr):
		c.emitBinaryExpr(expr)
	case (*ast.Ident):
		c.emitIdent(expr)
	case (*ast.CallExpr):
		c.emitCallExpr(expr)
	case (*ast.AssignExpr):
		c.emitAssignExpr(expr)
	case (*ast.ShortExpr):
		c.emitShortExpr(expr)
	default:
		fmt.Println("Type: ", typ)
	}
}

func (c *Compiler) emitType( /*Don't handle ast directly*/ typ token.Type) {
	c.buf.WriteString(typ.String())
}

func (c *Compiler) emitParamTypes( /*Don't handle ast directly*/ params *ast.StmtList) {
	c.buf.WriteString("(")
	for i, p := range params.List {
		d := p.(*ast.VarDeclStmt)
		c.buf.WriteString(d.Type.String())
		if i < len(params.List)-1 {
			c.buf.WriteString(", ")
		}
	}
	c.buf.WriteString(")")
}

func (c *Compiler) emitParams( /*Don't handle ast directly*/ params *ast.StmtList) {
	c.buf.WriteString("(")
	for i, p := range params.List {
		d := p.(*ast.VarDeclStmt)
		c.buf.WriteString(fmt.Sprintf("%s %s", d.Type.String(), d.Name.Name))
		if i < len(params.List)-1 {
			c.buf.WriteString(", ")
		}
	}
	c.buf.WriteString(")")
}

func (c *Compiler) write(s string) {
	for i := 0; i < c.tlevel; i++ {
		c.buf.WriteByte('\t')
	}
	c.buf.WriteString(s)
}

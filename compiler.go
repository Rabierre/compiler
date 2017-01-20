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
		fn := decl.(*FuncDecl)
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
		fn := decl.(*FuncDecl)
		c.emitType(fn.Type)
		c.buf.WriteByte(' ')
		c.buf.WriteString(fn.Name.Name)

		c.emitParams(fn.Params)
		c.buf.WriteByte('\n')
		c.emitBody(fn.Body)
	}

	ioutil.WriteFile(path.Join(c.output, "mid.c"), c.buf.Bytes(), 0777)
}

func (c *Compiler) emitBody( /*Don't handle ast directly*/ stnt Stmt) {
	c.emitCompoundStmt(stnt)

}

func (c *Compiler) emitCompoundStmt( /*Don't handle ast directly*/ stmt Stmt) {
	c.write("{\n")
	c.tlevel++
	cs := stmt.(*CompoundStmt)
	for _, s := range cs.List {
		c.emitStmt(s)
	}
	c.tlevel--
	c.write("}\n")
}

func (c *Compiler) emitIfStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*IfStmt)
	c.write("if (")
	c.emitExpr(s.Cond)
	c.buf.WriteString(")\n")
	c.emitBody(s.Body)
	c.write("else\n")
	c.emitBody(s.ElseBody)
}

func (c *Compiler) emitForStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*ForStmt)
	c.write("for (")
	c.emitShortDeclStmt(s.Init) // TODO emitShortVarDecl
	c.buf.WriteRune(';')
	c.emitExpr(s.Cond)
	c.buf.WriteRune(';')
	c.emitExpr(s.Post)
	c.buf.WriteString(")\n")
	c.emitBody(s.Body)
}

func (c *Compiler) emitReturnStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*ReturnStmt)
	c.write("return ")
	c.emitExpr(s.Value)
	c.buf.WriteString(";\n")
}

// for for stmt
func (c *Compiler) emitShortDeclStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*VarDeclStmt)
	c.buf.WriteString(s.Type.Val)
	c.buf.WriteRune(' ')
	c.buf.WriteString(s.Name.Name)
	if s.RValue != nil {
		c.buf.WriteRune('=')
		c.emitExpr(s.RValue)
	}
}

func (c *Compiler) emitVarDeclStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*VarDeclStmt)
	c.write("")
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
	c.write("")
	c.emitExpr(stmt.(*ExprStmt).expr)
	c.buf.WriteByte(';')
	c.buf.WriteByte('\n')
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
	e := expr.(*BinaryExpr)
	c.emitExpr(e.LValue)
	c.buf.WriteString(e.Op.Val.Kind.String())
	c.emitExpr(e.RValue)
}

func (c *Compiler) emitIdent( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*Ident)
	c.buf.WriteString(e.Name)
}

func (c *Compiler) emitCallExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*CallExpr)
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

func (c *Compiler) emitAssignExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*AssignExpr)
	c.emitExpr(e.LValue)
	c.buf.WriteRune('=')
	c.emitExpr(e.RValue)
}

func (c *Compiler) emitShortExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*ShortExpr)
	c.emitExpr(e.RValue)
	c.buf.WriteString(e.Op.Val.Kind.String())
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
	case (*ShortExpr):
		c.emitShortExpr(expr)
	default:
		fmt.Println("Type: ", typ)
	}
}

func (c *Compiler) emitType( /*Don't handle ast directly*/ typ token.Token) {
	c.buf.WriteString(typ.Kind.String())
}

func (c *Compiler) emitParamTypes( /*Don't handle ast directly*/ params *StmtList) {
	c.buf.WriteString("(")
	for i, p := range params.List {
		d := p.(*VarDeclStmt)
		c.buf.WriteString(d.Type.Kind.String())
		if i < len(params.List)-1 {
			c.buf.WriteString(", ")
		}
	}
	c.buf.WriteString(")")
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
	c.buf.WriteString(")")
}

func (c *Compiler) write(s string) {
	for i := 0; i < c.tlevel; i++ {
		c.buf.WriteByte('\t')
	}
	c.buf.WriteString(s)
}

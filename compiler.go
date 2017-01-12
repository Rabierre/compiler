package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/rabierre/compiler/token"
)

var buf bytes.Buffer

func Compile(src []byte) {
	parser := Parser{}
	parser.Init(src)
	parser.Parse()

	// function is top scope
	for _, decl := range parser.decls {
		fn := decl.(*FuncDecl)
		emitType(fn.Type)
		buf.WriteByte(' ')
		buf.WriteString(fn.Name.Name)

		emitParams(fn.Params)
		emitBody(fn.Body)
	}

	ioutil.WriteFile("mid.a", buf.Bytes(), 0777)
}

func emitBody( /*Don't handle ast directly*/ stnt Stmt) {
	emitCompoundStmt(stnt)
}

func emitCompoundStmt( /*Don't handle ast directly*/ stmt Stmt) {
	buf.WriteString("{\n")
	list := stmt.(*CompoundStmt).List
	for i := 0; i < len(list); i++ {
		emitStmt(list[i])
	}
	buf.WriteString("}\n")
}

func emitIfStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*IfStmt)
	buf.WriteString("if (")
	emitExpr(s.Cond)
	buf.WriteString(")\n")
	emitBody(s.Body)
	buf.WriteString("else\n")
	emitStmt(s.ElseBody)
}

func emitForStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*ForStmt)
	buf.WriteString("for (")
	emitStmt(s.Init) // TODO emitShortVarDecl
	buf.WriteRune(';')
	emitExpr(s.Cond)
	buf.WriteRune(';')
	emitExpr(s.Post)
	buf.WriteRune(')')
	emitBody(s.Body)
}

func emitReturnStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*ReturnStmt)
	buf.WriteString("return ")
	emitExpr(s.Value)
	buf.WriteString(";\n")
}

func emitVarDeclStmt( /*Don't handle ast directly*/ stmt Stmt) {
	s := stmt.(*VarDeclStmt)
	buf.WriteString(s.Type.Val)
	buf.WriteRune(' ')
	buf.WriteString(s.Name.Name)
	if s.RValue != nil {
		buf.WriteRune('=')
		emitExpr(s.RValue)
	}
	buf.WriteRune(';')
	buf.WriteRune('\n')
}

func emitExprStmt( /*Don't handle ast directly*/ stmt Stmt) {
	emitExpr(stmt.(*ExprStmt).expr)
}

func emitStmt( /*Don't handle ast directly*/ stmt Stmt) {
	switch typ := stmt.(type) {
	case (*CompoundStmt):
		emitCompoundStmt(stmt)
	case (*IfStmt):
		emitIfStmt(stmt)
	case (*ForStmt):
		emitForStmt(stmt)
	case (*ReturnStmt):
		emitReturnStmt(stmt)
	case (*VarDeclStmt):
		emitVarDeclStmt(stmt)
	case (*ExprStmt):
		emitExprStmt(stmt)
	default:
		println("Type: ", typ)
	}
}

func emitLiteracy( /*Don't handle ast directly*/ expr Expr) {
	ex := expr.(*BasicLit)
	// TODO is type need?
	// buf.WriteString(ex.Type.String())
	buf.WriteString(ex.Value)
}

func emitBinaryExpr( /*Don't handle ast directly*/ expr Expr) {
	bin := expr.(*BinaryExpr)
	emitExpr(bin.LValue)
	buf.WriteString(bin.Op.Val)
	emitExpr(bin.RValue)
}

func emitIdent( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*Ident)
	buf.WriteString(e.Name)
}

func emitCallExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*CallExpr)
	emitExpr(e.Name)
	buf.WriteRune('(')
	emitExprList(e.Params.List)
	buf.WriteRune(')')
}

func emitAssignExpr( /*Don't handle ast directly*/ expr Expr) {
	e := expr.(*AssignExpr)
	emitExpr(e.LValue)
	buf.WriteRune('=')
	emitExpr(e.RValue)
	buf.WriteRune(';')
	buf.WriteRune('\n')
}

func emitExprList(list []Expr) {
	for _, e := range list {
		emitExpr(e)
	}
}

// TODO maybe function chain
func emitExpr( /*Don't handle ast directly*/ expr Expr) {
	switch typ := expr.(type) {
	case (*BasicLit):
		emitLiteracy(expr)
	case (*BinaryExpr):
		emitBinaryExpr(expr)
	case (*Ident):
		emitIdent(expr)
	case (*CallExpr):
		emitCallExpr(expr)
	case (*AssignExpr):
		emitAssignExpr(expr)
	default:
		fmt.Println("Type: ", typ)
	}
}

func emitType( /*Don't handle ast directly*/ typ token.Token) {
	buf.WriteString(typ.Kind.String())
}

func emitParams( /*Don't handle ast directly*/ params *StmtList) {
	buf.WriteString("(")
	for i, p := range params.List {
		d := p.(*VarDeclStmt)
		buf.WriteString(fmt.Sprintf("%s %s", d.Type.Kind.String(), d.Name.Name))
		if i < len(params.List)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")\n")
}

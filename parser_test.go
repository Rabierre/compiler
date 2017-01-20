package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/rabierre/compiler/ast"
	"github.com/rabierre/compiler/token"
	"github.com/stretchr/testify/assert"
)

func initParser(src string) *Parser {
	parser := &Parser{}
	parser.Init([]byte(src))
	return parser
}

func TestParseComment(t *testing.T) {
	parser := initParser("// this is comment")
	parser.parseComment()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, 1, len(parser.comments.List))
}

func TestParseFunction(t *testing.T) {
	src := `func func1() {}
		// Comment 1
		func func2() {
			// Comment 2
		}
		func func3() {
			for(int i = 0; i < 10; i++) {
				// Comment 3
			}
		}
		func func4(int a, double b) int {
			return a
		}
		func func5(int a, double b) int {
			func4(a, b)
		}
	`
	parser := initParser(src)
	parser.Parse()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, 5, len(parser.decls))
	assert.Equal(t, "func1", parser.decls[0].(*ast.FuncDecl).Name.Name)
	assert.Equal(t, 0, len(parser.decls[0].(*ast.FuncDecl).Body.List))
	assert.Equal(t, "func2", parser.decls[1].(*ast.FuncDecl).Name.Name)
	assert.Equal(t, 0, len(parser.decls[1].(*ast.FuncDecl).Body.List))
}

func TestparseForStmt(t *testing.T) {
	src := `for (int i = 0;i < 10; i++) {
			// Comment

		}
	`
	parser := initParser(src)
	stmt := parser.parseForStmt().(*ast.ForStmt)
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.Cond)
	assert.NotNil(t, stmt.Post)
	assert.NotNil(t, stmt.Post.(*ast.ShortExpr))
}

func TestParseIfStmt(t *testing.T) {
	src := `if (1 == 2) {
			// Comment

		}
		else {
			// comment
		}
	`
	parser := initParser(src)
	stmt := parser.parseIfStmt().(*ast.IfStmt)
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.Cond)
	assert.NotNil(t, stmt.ElseBody)

	cond := stmt.Cond.(*ast.BinaryExpr)
	assert.True(t, DeepEqual(&ast.BasicLit{Pos: 4, Value: "1", Type: token.INT_LIT}, cond.LValue))
	assert.True(t, DeepEqual(ast.Operator{Type: token.EQ}, cond.Op))
	assert.True(t, DeepEqual(&ast.BasicLit{Pos: 9, Value: "2", Type: token.INT_LIT}, cond.RValue))

	src = `if (1 == 2) {
			// Comment
			// Comment
			// Comment
		}
	`
	parser = initParser(src)
	stmt = parser.parseIfStmt().(*ast.IfStmt)
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.Cond)
	assert.Nil(t, stmt.ElseBody)
}

func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func TestParseVarDecl(t *testing.T) {
	src := `int a = 10`
	parser := initParser(src)
	stmt := parser.parseVarDecl().(*ast.VarDeclStmt)
	assert.NotNil(t, stmt)
	assert.Equal(t, "a", stmt.Name.Name)
	assert.Equal(t, "10", stmt.RValue.(*ast.BasicLit).Value)

	src = `int a = funcCall(b,c)`
	parser = initParser(src)
	stmt = parser.parseVarDecl().(*ast.VarDeclStmt)
	assert.NotNil(t, stmt)
	assert.Equal(t, "a", stmt.Name.Name)
	assert.Equal(t, "funcCall", stmt.RValue.(*ast.CallExpr).Name.(*ast.Ident).Name)
}

func TestParseReturnStmt(t *testing.T) {
	src := `return`
	parser := initParser(src)
	stmt := parser.parseReturnStmt().(*ast.ReturnStmt)
	assert.NotNil(t, stmt)
	assert.Nil(t, stmt.Value)

	src = `return funcCall(a,1,true)`
	parser = initParser(src)
	stmt = parser.parseReturnStmt().(*ast.ReturnStmt)
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.Value)
	params := stmt.Value.(*ast.CallExpr).Params.List

	assert.Equal(t, 3, len(params))
	set := [][]string{
		{"Ident", reflect.TypeOf(params[0]).String()},
		{"BasicLit", reflect.TypeOf(params[1]).String()},
		{"BasicLit", reflect.TypeOf(params[2]).String()},
	}
	for _, s := range set {
		assert.True(t, strings.Contains(s[1], s[0]))
	}
}

func TestParseExprStmt(t *testing.T) {
	src := `
		a = 10
	`
	parser := initParser(src)
	stmt := parser.parseExprStmt()
	e := stmt.(*ast.ExprStmt).Val.(*ast.AssignExpr)
	assert.Equal(t, "a", e.LValue.(*ast.Ident).Name)
	assert.Equal(t, "10", e.RValue.(*ast.BasicLit).Value)
}

func TestNext(t *testing.T) {
	src := `
		// Comment 1
		// Comment 2
		int a
		// Comment 3
		// Comment 4
		int b
	`
	parser := initParser(src)
	expects := []token.Type{token.INT, token.IDENT, token.INT, token.IDENT}
	for _, exp := range expects {
		assert.Equal(t, exp, parser.tok)
		parser.next()
	}
}

func TestResolve(t *testing.T) {
	src := `
		func1()
	`
	parser := initParser(src)
	parser.parseExpr(true)
	assert.Equal(t, 1, len(parser.UnResolved))

	// TODO parse exprstmt need
	src = `
		func func1() {
			func1()
			func2()
			a = func1()
		}
	`
	parser = initParser(src)
	assert.Panics(t, func() {
		parser.Parse()
	})
	assert.Equal(t, 2, len(parser.UnResolved))
}

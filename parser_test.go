package main

import (
	"reflect"
	"testing"

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
	parser.Parse()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, len(parser.comments.comments), 1)
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
	`
	parser := initParser(src)
	parser.Parse()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, 4, len(parser.decls))
	assert.Equal(t, "func1", parser.decls[0].(*FuncDecl).Name.Name)
	assert.Equal(t, 0, len(parser.decls[0].(*FuncDecl).Body.List))
	assert.Equal(t, "func2", parser.decls[1].(*FuncDecl).Name.Name)
	assert.Equal(t, 0, len(parser.decls[1].(*FuncDecl).Body.List))
}

func TestparseForStmt(t *testing.T) {
	src := `for (int i = 0;i < 10; i++) {
			// Comment

		}
	`
	parser := initParser(src)
	stmt := parser.parseForStmt()
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.(*ForStmt).Cond)
	assert.NotNil(t, stmt.(*ForStmt).Post)
	assert.NotNil(t, stmt.(*ForStmt).Post.(*ShortExpr))
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
	stmt := parser.parseIfStmt()
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.(*IfStmt).Cond)
	assert.NotNil(t, stmt.(*IfStmt).ElseBody)

	cond := stmt.(*IfStmt).Cond.(*BinaryExpr)
	assert.True(t, DeepEqual(&BasicLit{Pos: 4, Value: "1", Type: token.INT_LIT}, cond.LValue))
	assert.True(t, DeepEqual(token.Token{"==", token.EQ}, cond.Op))
	assert.True(t, DeepEqual(&BasicLit{Pos: 9, Value: "2", Type: token.INT_LIT}, cond.RValue))

	src = `if (1 == 2) {
			// Comment
			// Comment
			// Comment
		}
	`
	parser = initParser(src)
	stmt = parser.parseIfStmt()
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.(*IfStmt).Cond)
	assert.Nil(t, stmt.(*IfStmt).ElseBody)
}

func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func TestParseVarDecl(t *testing.T) {
	src := `int a = 10`
	parser := initParser(src)
	stmt := parser.parseVarDecl()
	assert.NotNil(t, stmt)
	varDecl := stmt.(*VarDeclStmt)
	assert.Equal(t, "a", varDecl.Name.Name)
	assert.Equal(t, "10", varDecl.RValue.(*BasicLit).Value)

	src = `int a = funcCall(b,c)`
	parser = initParser(src)
	stmt = parser.parseVarDecl()
	assert.NotNil(t, stmt)
	varDecl = stmt.(*VarDeclStmt)
	assert.Equal(t, "a", varDecl.Name.Name)
	assert.Equal(t, "funcCall", varDecl.RValue.(*CallExpr).Name.(*Ident).Name)
}

func TestParseReturnStmt(t *testing.T) {
	src := `return`
	parser := initParser(src)
	stmt := parser.parseReturnStmt()
	assert.NotNil(t, stmt)
	assert.Nil(t, stmt.(*ReturnStmt).Value)

	src = `return funcCall(a,1,true)`
	parser = initParser(src)
	stmt = parser.parseReturnStmt()
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.(*ReturnStmt).Value)
	params := stmt.(*ReturnStmt).Value.(*CallExpr).Params.List

	assert.Equal(t, 3, len(params))
	set := [][]string{
		{"*main.Ident", reflect.TypeOf(params[0]).String()},
		{"*main.BasicLit", reflect.TypeOf(params[1]).String()},
		{"*main.BasicLit", reflect.TypeOf(params[2]).String()},
	}
	for _, s := range set {
		assert.Equal(t, s[0], s[1])
	}
}

func TestParseExprStmt(t *testing.T) {
	src := `
		a = 10
	`
	parser := initParser(src)
	stmt := parser.parseExprStmt()
	e := stmt.(*ExprStmt).expr.(*AssignExpr)
	assert.Equal(t, "a", e.LValue.(*Ident).Name)
	assert.Equal(t, "10", e.RValue.(*BasicLit).Value)
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
		tok, _ := parser.next()
		assert.Equal(t, exp, tok.Kind)
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

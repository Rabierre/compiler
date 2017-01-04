package main

import (
	"reflect"
	"testing"

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
		func func4(int a, double b) {
		}
	`
	parser := initParser(src)
	parser.Parse()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, 4, len(parser.decls))
	assert.Equal(t, "func1", parser.decls[0].(*FuncDecl).Name.Name.val)
	assert.Equal(t, 0, len(parser.decls[0].(*FuncDecl).Body.List))
	assert.Equal(t, "func2", parser.decls[1].(*FuncDecl).Name.Name.val)
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
	assert.True(t, DeepEqual(&BasicLit{Pos: 4, Value: "1", Type: IntLit}, cond.LValue))
	assert.True(t, DeepEqual(Token{val: "==", kind: EqType}, cond.Op))
	assert.True(t, DeepEqual(&BasicLit{Pos: 9, Value: "2", Type: IntLit}, cond.RValue))

	src = `if (1 == 2) {
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
	assert.Equal(t, "a", varDecl.Name.Name.val)
	assert.Equal(t, "10", varDecl.RValue.(*BasicLit).Value)

	src = `int a = funcCall(b,c)`
	parser = initParser(src)
	stmt = parser.parseVarDecl()
	assert.NotNil(t, stmt)
	varDecl = stmt.(*VarDeclStmt)
	assert.Equal(t, "a", varDecl.Name.Name.val)
	assert.Equal(t, "funcCall", varDecl.RValue.(*CallExpr).Name.(*Ident).Name.val)
}

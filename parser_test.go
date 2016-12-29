package main

import (
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
        // Comment Here
        func func2() {

        }
    `
	parser := initParser(src)
	parser.Parse()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, 2, len(parser.decls))
	assert.Equal(t, "func1", parser.decls[0].(*FuncDecl).Name.Name)
	assert.Equal(t, 0, len(parser.decls[0].(*FuncDecl).Body.List))
	assert.Equal(t, "func2", parser.decls[1].(*FuncDecl).Name.Name)
	assert.Equal(t, 0, len(parser.decls[1].(*FuncDecl).Body.List))
	// assert.Equal(t, expected, actual, ...)
}

func TestParseFor(t *testing.T) {
	src := `for (;;) {

        }
    `
	parser := initParser(src)
	stmt := parser.parseStmt()
	assert.NotNil(t, stmt)
	assert.NotNil(t, stmt.(*ForStmt).Cond)
}

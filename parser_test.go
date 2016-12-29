package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseComment(t *testing.T) {
	parser := Parser{}
	parser.Init([]byte("// this is comment"))
	parser.Parse()

	assert.NotNil(t, parser.topScope)
	assert.Equal(t, len(parser.comments.comments), 1)
}

func TestParseFunction(t *testing.T) {
	parser := Parser{}
	parser.Init([]byte(`func func1() {}
        // Comment Here
        func func2() {
        }
        `))
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
	// initScanner("for (;;) {}")
}

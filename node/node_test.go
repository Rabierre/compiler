package node

import (
	"testing"

	"github.com/rabierre/compiler/ast"
	"github.com/stretchr/testify/assert"
)

func TestDecls(t *testing.T) {
	name := &ast.Ident{}
	l := []ast.Decl{&ast.FuncDecl{Name: name}}

	res := decls(l)
	assert.NotNil(t, res[0])
	assert.NotNil(t, res[0].Func.NParam)
	assert.NotNil(t, res[0].Func.NBody)
}

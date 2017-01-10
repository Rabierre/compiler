package main

import (
	"bytes"
	"fmt"
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
		buf.WriteString(fn.Name.Name.val)

		emitParams(fn.Params)
		emitBody(fn.Body)
	}

	fmt.Print(buf.String())
}

func emitBody( /*Don't handle ast directly*/ body Stmt) {
	buf.WriteString("{\n")

	list := body.(*CompoundStmt).List
	for i := 0; i < len(list); i++ {
		emitStmt(list[i])
		buf.WriteByte('\n')
	}

	buf.WriteString("}\n")
}

func emitStmt( /*Don't handle ast directly*/ stmt Stmt) {
	switch typ := stmt.(type) {
	case (*CompoundStmt):
		println("*CompoundStmt")
	case (*IfStmt):
		println("*IfStmt")
	case (*ForStmt):
		println("*ForStmt")
	case (*ReturnStmt):
		println("*ReturnStmt")
	default:
		println("Type: ", typ)
	}
}

func emitType( /*Don't handle ast directly*/ typ Token) {
	buf.WriteString(typ.kind.String())
}

func emitParams( /*Don't handle ast directly*/ params *ArgList) {
	buf.WriteString("(")
	for i := 0; i < len(params.List); i++ {
		buf.WriteString(params.List[i].Name.Name.val)
		if i < len(params.List)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")\n")
}

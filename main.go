package main

import (
	"bytes"
	"io"
)

type TokenKind struct {
	val string
}

const (
	If       = "if"
	For      = "for"
	LParen   = "("
	RParen   = ")"
	Int      = "int"
	Double   = "double"
	Plus     = "+"
	Minus    = "-"
	Multi    = "*"
	Divide   = "/"
	Assign   = "="
	Less     = "<"
	Great    = ">"
	LessEq   = "<="
	GreateEq = ">="
	Equal    = "=="
	NotEq    = "!="
	Space    = " "
)

type Token struct {
	kind TokenKind
}

func IsSpace(ch string) bool {
	if ch == Space {
		return true
	}
	return false
}

func NextToken() Token {
	// TODO skip space
	ch, err := NextChar()
	for IsSpace(ch) {
		ch, err = NextChar()
		if err != nil && err != io.EOF {
			panic(err)
		}
	}

	text := ""

	switch Kind(ch) {
	case LETTER:
		for ch != Space && err != io.EOF {
			text += ch
			ch, err = NextChar()
		}
	case DIGIT:
	default:
		// TODO store
	}

	// get token from input and
	// check if token can be made
	// else get next char

	return Token{}
}

type CharKind int

const (
	LETTER CharKind = iota
	DIGIT
	DOUBLE_QUOTE
	OTHER
)

func Kind(ch string) CharKind {
	switch ch {
	case "a", "b", "c":
		return LETTER
	case "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9":
		return DIGIT
	case ".":
		return DOUBLE_QUOTE
	default:
		return OTHER
	}
}

var fin *bytes.Buffer

func NextChar() (string, error) {
	ch, _, err := fin.ReadRune()
	return string(ch), err
}

func main() {
	fin = &bytes.Buffer{}
	fin.Write([]byte("if else for + - * /"))

	for fin.Len() > 0 {
		ch, _ := NextChar()
		println(ch)
	}
}

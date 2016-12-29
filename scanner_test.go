package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initScanner(src string) *Scanner {
	scanner := Scanner{}
	scanner.Init()
	scanner.src = []byte(src)
	return &scanner
}

func TestNextCh(t *testing.T) {
	scanner := initScanner("if else for + - * /")
	for {
		ch, err := scanner.nextCh()
		if err != nil && err == io.EOF {
			break
		}
		println(ch)
	}
}

func TestScan(t *testing.T) {
	scanner := initScanner("func main() {}")
	tokType := []TokenType{FuncType, IdentType, LParenType, RParenType, LBraceType, RBraceType, EOFType}
	for i := 0; !scanner.fullScaned; i++ {
		tok, _ := scanner.next()
		assert.Equal(t, tokType[i], tok.kind)
	}

	scanner = initScanner(`for (;;) {
		}
	`)
	tokType = []TokenType{ForType, LParenType, SemiColType, SemiColType, RParenType, LBraceType, RBraceType, EOFType}
	for i := 0; !scanner.fullScaned; i++ {
		tok, _ := scanner.next()
		assert.Equal(t, tokType[i], tok.kind)
	}
}

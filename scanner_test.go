package main

import (
	"io"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func initScanner(src string) {
	scanner.src = []byte(src)
	scanner.tokens = []Token{}
	scanner.srcIndex = -1
	scanner.tokenIndex = -1
}

func TestNextChar(t *testing.T) {
	initScanner("if else for + - * /")
	for {
		ch, err := NextChar()
		if err != nil && err == io.EOF {
			break
		}
		println(ch)
	}
}

func TestTokenize(t *testing.T) {
	initScanner("if else for + - * /")
	tokens := Tokenize()
	fmt.Println(tokens)
	assert.Equal(t, len(tokens), 7)

	initScanner("func main() {}")
	tokens = Tokenize()
	fmt.Println(tokens)
	assert.Equal(t, len(tokens), 3)

	initScanner("for (;;) {}")
	tokens = Tokenize()
	fmt.Println(tokens)
	assert.Equal(t, len(tokens), 3)
}

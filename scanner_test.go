package main

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/rabierre/compiler/token"
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
	expect := strings.Split("if else for + - * /", "")
	var result []string
	for {
		ch, err := scanner.nextCh()
		if err == io.EOF {
			break
		}
		result = append(result, ch)
	}
	assert.True(t, reflect.DeepEqual(result, expect))
}

type Suite struct {
	src    string
	tokens []token.Type
}

func SuiteCase() []*Suite {
	return []*Suite{
		&Suite{"func main(int a, double b) {}", []token.Type{token.FUNC, token.IDENT, token.LPAREN, token.INT, token.IDENT, token.COMMA, token.DOUBLE, token.IDENT, token.RPAREN, token.LBRACE, token.RBRACE, token.EOF}},
		&Suite{`for (;;) {
	 		}
	 	`, []token.Type{token.FOR, token.LPAREN, token.SEMI_COLON, token.SEMI_COLON, token.RPAREN, token.LBRACE, token.RBRACE, token.EOF}},
		&Suite{`if (1 == 2) {
				// comment
			} else {
				// comment
			}
		`, []token.Type{token.IF, token.LPAREN, token.INT_LIT, token.EQ, token.INT_LIT, token.RPAREN, token.LBRACE, token.COMMENT, token.IDENT, token.RBRACE, token.ELSE, token.LBRACE, token.COMMENT, token.IDENT, token.RBRACE, token.EOF}},
		&Suite{`func func3() {
					for(int i = 0; i < 10; i++) {
					// Comment
				}
			}
		`, []token.Type{token.FUNC, token.IDENT, token.LPAREN, token.RPAREN, token.LBRACE, token.FOR, token.LPAREN, token.INT, token.IDENT, token.ASSIGN, token.INT_LIT, token.SEMI_COLON, token.IDENT, token.LESS, token.INT_LIT, token.SEMI_COLON, token.IDENT, token.INC, token.RPAREN, token.LBRACE, token.COMMENT, token.IDENT, token.RBRACE, token.RBRACE, token.EOF}},
	}
}

func TestScan(t *testing.T) {
	for _, suite := range SuiteCase() {
		scanner := initScanner(suite.src)
		for i := 0; !scanner.fullScaned; i++ {
			tok, _ := scanner.next()
			assert.Equal(t, suite.tokens[i], tok.Kind)
		}
	}
}

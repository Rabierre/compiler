package main

import (
	"io"
	"reflect"
	"strings"
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
	tokens []TokenType
}

func SuiteCase() []*Suite {
	return []*Suite{
		&Suite{"func main(int a, double b) {}", []TokenType{FUNC, IDENT, LPAREN, INT, IDENT, COMMA, DOUBLE, IDENT, RPAREN, LBRACE, RBRACE, EOF}},
		&Suite{`for (;;) {
	 		}
	 	`, []TokenType{FOR, LPAREN, SEMI_COLON, SEMI_COLON, RPAREN, LBRACE, RBRACE, EOF}},
		&Suite{`if (1 == 2) {
				// comment
			} else {
				// comment
			}
		`, []TokenType{IF, LPAREN, INT_LIT, EQ, INT_LIT, RPAREN, LBRACE, COMMENT_SLASH, IDENT, RBRACE, ELSE, LBRACE, COMMENT_SLASH, IDENT, RBRACE, EOF}},
		&Suite{`func func3() {
					for(int i = 0; i < 10; i++) {
					// Comment
				}
			}
		`, []TokenType{FUNC, IDENT, LPAREN, RPAREN, LBRACE, FOR, LPAREN, INT, IDENT, ASSIGN, INT_LIT, SEMI_COLON, IDENT, LESS, INT_LIT, SEMI_COLON, IDENT, RPAREN, LBRACE, COMMENT_SLASH, IDENT, RBRACE, RBRACE, EOF}},
	}
}

func TestScan(t *testing.T) {
	for _, suite := range SuiteCase() {
		scanner := initScanner(suite.src)
		for i := 0; !scanner.fullScaned; i++ {
			tok, _ := scanner.next()
			assert.Equal(t, suite.tokens[i], tok.kind)
		}
	}
}

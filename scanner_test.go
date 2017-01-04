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
		&Suite{"func main(int a, double b) {}", []TokenType{FuncType, IdentType, LParenType, IntType, IdentType, CommaType, DoubleType, IdentType, RParenType, LBraceType, RBraceType, EOFType}},
		&Suite{`for (;;) {
	 		}
	 	`, []TokenType{ForType, LParenType, SemiColType, SemiColType, RParenType, LBraceType, RBraceType, EOFType}},
		&Suite{`if (1 == 2) {
				// comment
			} else {
				// comment
			}
		`, []TokenType{IfType, LParenType, IntLit, EqType, IntLit, RParenType, LBraceType, CommentType, IdentType, RBraceType, ElseType, LBraceType, CommentType, IdentType, RBraceType, EOFType}},
		&Suite{`func func3() {
					for(int i = 0; i < 10; i++) {
					// Comment
				}
			}
		`, []TokenType{FuncType, IdentType, LParenType, RParenType, LBraceType, ForType, LParenType, IntType, IdentType, AssignType, IntLit, SemiColType, IdentType, LessType, IntLit, SemiColType, IdentType, RParenType, LBraceType, CommentType, IdentType, RBraceType, RBraceType, EOFType}},
	}
}

func TestScan(t *testing.T) {
	TestSuite := SuiteCase()
	for _, suite := range TestSuite {
		scanner := initScanner(suite.src)
		for i := 0; !scanner.fullScaned; i++ {
			tok, _ := scanner.next()
			assert.Equal(t, suite.tokens[i], tok.kind)
		}
	}
}

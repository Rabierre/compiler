package main

import (
	"strconv"
)

type TokenType int

type Token struct {
	val  string
	kind TokenType
}

const (
	BEGIN TokenType = iota

	IF
	ELSE
	FOR
	FUNC
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	INT
	DOUBLE
	INT_LIT
	DOUBLE_LIT
	PLUS
	MINUS
	MULTI
	DIVIDE
	ASSIGN
	LESS
	GRT
	LEQ
	GEQ
	EQ
	NEQ
	IDENT
	COMMENT_SLASH
	RETURN
	SEMI_COLON
	COMMA
	TRUE
	FALSE
	SPACE

	END

	EOF
)

var Tokens map[string]TokenType

func init() {
	Tokens = make(map[string]TokenType)
	for i := BEGIN + 1; i < END; i++ {
		Tokens[Keywords[i]] = i
	}
}

func KeywordType(token string) TokenType {
	typ, exist := Tokens[token]
	if exist {
		return typ
	}
	return IDENT
}

func (t TokenType) String() string {
	s := ""
	if 0 <= t && t < TokenType(len(Keywords)) {
		s = Keywords[t]
	}
	if s == "" {
		s = "Token: " + strconv.Itoa(int(t))
	}

	return s
}

type CharType int

const (
	LETTER_LIT CharType = iota
	DIGIT_LIT
	DOT_LIT
	LBRACE_LIT
	RBRACE_LIT
	LPAREN_LIT
	RPAREN_LIT
	SEMICOLON_LIT
	COMMA_LIT
	OTHER_LIT
)

func Kind(ch string) CharType {
	if ch == "" {
		return OTHER_LIT
	}

	c := []rune(ch)[0]
	switch {
	case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		return LETTER_LIT
	case '0' <= c && c <= '9':
		return DIGIT_LIT
	case c == '.':
		return DOT_LIT
	case c == '{':
		return LBRACE_LIT
	case c == '}':
		return RBRACE_LIT
	case c == '(':
		return LPAREN_LIT
	case c == ')':
		return RPAREN_LIT
	case c == ';':
		return SEMICOLON_LIT
	case c == ',':
		return COMMA_LIT
	default:
		return OTHER_LIT
	}
}

const (
	LowestPriority  = 0 // non-operators
	UnaryPriority   = 6
	HighestPriority = 7
)

func (t Token) Priority() int {
	switch t.kind {
	// case LOR:
	// 	return 1
	// case LAND:
	// 	return 2
	case EQ, NEQ, LESS, LEQ, GRT, GEQ:
		return 3
	case PLUS, MINUS:
		return 4
	case MULTI, DIVIDE:
		return 5
	}
	return LowestPriority
}

package token

import (
	"strconv"
)

type Type int

type Token struct {
	Val  string
	Kind Type
}

const (
	BEGIN Type = iota

	IF
	ELSE
	FOR
	FUNC
	INT
	DOUBLE
	RETURN
	TRUE
	FALSE

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	COMMENT
	SEMI_COLON
	COMMA

	PLUS
	MINUS
	INC
	DEC
	MULTI
	DIVIDE
	ASSIGN
	LESS
	GRT
	LEQ
	GEQ
	EQ
	NEQ
	SPACE
	VOID

	END

	INT_LIT
	DOUBLE_LIT
	IDENT
	EOF
)

var Tokens map[string]Type

func init() {
	Tokens = make(map[string]Type)
	for i := BEGIN + 1; i < END; i++ {
		Tokens[Keywords[i]] = i
	}
}

func KeywordType(token string) Type {
	typ, exist := Tokens[token]
	if exist {
		return typ
	}
	return IDENT
}

func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(Keywords)) {
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
	switch t.Kind {
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

var Keywords = [...]string{
	IF:     "if",
	ELSE:   "else",
	FOR:    "for",
	FUNC:   "func",
	INT:    "int",
	DOUBLE: "double",
	RETURN: "return",
	TRUE:   "true",
	FALSE:  "false",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",

	COMMENT:    "//",
	SEMI_COLON: ";",
	COMMA:      ",",

	PLUS:   "+",
	MINUS:  "-",
	INC:    "++",
	DEC:    "--",
	MULTI:  "*",
	DIVIDE: "/",
	ASSIGN: "=",
	LESS:   "<",
	GRT:    ">",
	LEQ:    "<=",
	GEQ:    ">=",
	EQ:     "==",
	NEQ:    "!=",

	SPACE: " ",
	VOID:  "void",
}

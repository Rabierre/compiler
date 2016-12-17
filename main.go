package main

import (
	"bytes"
	"io"

	"fmt"
)

func Tokenizer() []Token {
	tokens := []Token{}
	for {
		token, err := NextToken()
		tokens = append(tokens, token)
		if err != nil && err == io.EOF {
			break
		}

	}
	return tokens
}

func NextToken() (Token, error) {
	ch, err := NextChar()
	if err != nil && err == io.EOF {
		return Token{ch, EOFType}, err
	}

	for IsSpace(ch) {
		ch, err = NextChar()
		if err != nil && err != io.EOF {
			panic(err)
		}
	}

	text := ""
	isNum := false

	switch Kind(ch) {
	case LETTER:
		for ch != Space && err != io.EOF {
			text += ch
			ch, err = NextChar()
		}
	case DIGIT:
		isNum = true
		for ch != Space && err != io.EOF {
			text += ch
			if Kind(ch) == LETTER {
				panic("Invalid variable name: " + text)
			}
			ch, err = NextChar()
		}
	default: // Operator
		for ch != Space && err != io.EOF {
			text += ch
			if Kind(PeepChar()) != OTHER {
				break
			}
			ch, err = NextChar()
		}
	}

	return Tokenize(text, isNum), err
}

func IsSpace(ch string) bool {
	if ch == Space {
		return true
	}
	return false
}

func Tokenize(token string, num bool) Token {
	if num {
		return Token{token, NumberType}
	}

	kind := KeywordType(token)
	return Token{token, kind}
}

func KeywordType(token string) (kind TokenType) {
	switch token {
	case If:
		kind = IfType
	case Else:
		kind = ElseType
	case For:
		kind = ForType
	case LParen:
		kind = LParenType
	case RParen:
		kind = RParenType
	case Int:
		kind = IntType
	case Double:
		kind = DoubleType
	case Plus:
		kind = PlusType
	case Minus:
		kind = MinusType
	case Multi:
		kind = MultiType
	case Divide:
		kind = DivideType
	case Assign:
		kind = AssignType
	case Less:
		kind = LessType
	case Great:
		kind = GreatType
	case LessEq:
		kind = LessEqType
	case GreateEq:
		kind = GreateEqType
	case Equal:
		kind = EqualType
	case NotEq:
		kind = NotEqType
	case Space:
		kind = SpaceType
	default:
		kind = IdentType
	}
	return kind
}

var fin *bytes.Buffer

func NextChar() (string, error) {
	ch, _, err := fin.ReadRune()
	return string(ch), err
}

func PeepChar() string {
	ch, _, _ := fin.ReadRune()
	fin.UnreadRune()
	return string(ch)
}

func Kind(ch string) CharType {
	switch ch {
	case "a", "b", "c", "d", "e",
		"f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o",
		"p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z":
		return LETTER
	case "0", "1", "2", "3", "4",
		"5", "6", "7", "8", "9":
		return DIGIT
	case ".":
		return DOUBLE_QUOTE
	default:
		return OTHER
	}
}

func main() {
	fin = &bytes.Buffer{}
	fin.Write([]byte("if else for + - * /"))
	for fin.Len() > 0 {
		ch, _ := NextChar()
		println(ch)
	}

	fin.Write([]byte("if else for + - * / els 1234 12.34 <= >="))
	tokens := Tokenizer()
	fmt.Println(tokens)
}

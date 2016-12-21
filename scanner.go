package main

import "io"

type Scanner struct {
	src        []byte
	srcIndex   int
	tokens     []Token
	tokenIndex int
}

var scanner *Scanner

func init() {
	scanner = &Scanner{tokens: []Token{}, tokenIndex: -1, srcIndex: -1}
}

func Tokenize() []Token {
	for {
		token, err := NextToken()
		scanner.push(token)
		if err != nil && err == io.EOF {
			break
		}
	}
	return scanner.tokens
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

	return ToToken(text, isNum), err
}

func NextChar() (string, error) {
	scanner.srcIndex += 1
	if scanner.srcIndex >= len(scanner.src) {
		return "", io.EOF
	}
	return string(scanner.src[scanner.srcIndex]), nil
}

func ToToken(token string, num bool) Token {
	if num {
		return Token{token, NumberType}
	}

	kind := KeywordType(token)
	return Token{token, kind}
}

func IsSpace(ch string) bool {
	if ch == Space {
		return true
	}
	return false
}

func (s *Scanner) next() Token {
	s.tokenIndex += 1
	return s.tokens[s.tokenIndex]
}

func (s *Scanner) push(token Token) {
	s.tokens = append(s.tokens, token)
}

func PeepChar() string {
	return string(scanner.src[scanner.srcIndex])
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

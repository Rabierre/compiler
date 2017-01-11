package main

import (
	"io"
	"strings"

	"github.com/rabierre/compiler/token"
)

type Scanner struct {
	src        []byte
	srcIndex   int
	tokens     []token.Token
	tokenIndex int
	fullScaned bool
}

func (s *Scanner) Init() {
	s.srcIndex = -1
}

func (s *Scanner) peek() (token.Token, int) {
	index := s.srcIndex
	tok, pos := s.next()
	s.srcIndex = index
	return tok, pos
}

// For comment
func (s *Scanner) nextLine() (token.Token, int) {
	text := ""
	pos := s.srcIndex
	if pos < 0 {
		pos = 0
	}

	ch, err := s.nextCh()
	for ch != "\n" && err != io.EOF {
		text += ch
		ch, err = s.nextCh()
	}
	return token.Token{text, token.COMMENT}, pos
}

func (s *Scanner) next() (token.Token, int) {
	ch, err := s.skipWhiteSpace()

	pos := s.srcIndex

	if err != nil && err == io.EOF {
		s.fullScaned = true
		return token.Token{"", token.EOF}, pos
	}

	text := ""
	isNum := false
	switch token.Kind(ch) {
	case token.LETTER_LIT:
		for ch != "\n" && ch != "\t" && ch != " " && err != io.EOF {
			if ch == token.Keywords[token.LPAREN] ||
				ch == token.Keywords[token.RPAREN] ||
				ch == token.Keywords[token.COMMA] {
				s.undoCh()
				break
			}
			text += ch
			ch, err = s.nextCh()
		}
	case token.DIGIT_LIT:
		isNum = true
		for ch != " " && err != io.EOF {
			text += ch
			if token.Kind(ch) == token.LETTER_LIT {
				panic("Invalid variable name: " + text)
			}

			ch, err = s.nextCh()
			if kind := token.Kind(ch); kind != token.DIGIT_LIT && kind != token.DOT_LIT {
				s.undoCh()
				break
			}
		}
	case token.COMMA_LIT:
		text += ch
	default: // Operator
		for ch != " " && ch != "\n" && err != io.EOF {
			text += ch

			ch, err = s.nextCh()
			if token.Kind(ch) != token.OTHER_LIT {
				s.undoCh()
				break
			}
		}
	}

	if err != nil && err == io.EOF {
		// if err and not EOF, increase err count
		s.fullScaned = true
	}

	return ToToken(text, isNum), pos
}

func (s *Scanner) nextCh() (string, error) {
	s.srcIndex += 1
	if s.srcIndex >= len(s.src) {
		return "", io.EOF
	}
	return string(s.src[s.srcIndex]), nil
}

func (s *Scanner) PeepCh() (string, error) {
	if s.srcIndex+1 >= len(s.src) {
		return "", io.EOF
	}
	return string(s.src[s.srcIndex+1]), nil
}

func (s *Scanner) undoCh() {
	if 0 <= s.srcIndex-1 {
		s.srcIndex -= 1
	}
}

func (s *Scanner) skipWhiteSpace() (string, error) {
	ch, err := s.nextCh()
	for ch == " " || ch == "\n" || ch == "\t" || ch == "\r" {
		ch, err = s.nextCh()
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
	return ch, err
}

func ToToken(keyword string, num bool) token.Token {
	if num {
		if strings.Contains(keyword, ".") {
			return token.Token{keyword, token.DOUBLE_LIT}
		} else {
			return token.Token{keyword, token.INT_LIT}
		}
	}

	kind := token.KeywordType(keyword)
	return token.Token{keyword, kind}
}

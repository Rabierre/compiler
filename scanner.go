package main

import (
	"io"
	"strings"
)

type Scanner struct {
	src        []byte
	srcIndex   int
	tokens     []Token
	tokenIndex int
	fullScaned bool
}

func (s *Scanner) Init() {
	s.srcIndex = -1
}

func (s *Scanner) peek() (Token, int) {
	index := s.srcIndex
	tok, pos := s.next()
	s.srcIndex = index
	return tok, pos
}

// For comment
func (s *Scanner) nextLine() (Token, int) {
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
	return Token{text, COMMENT_SLASH}, pos
}

func (s *Scanner) next() (Token, int) {
	ch, err := s.skipWhiteSpace()

	pos := s.srcIndex

	if err != nil && err == io.EOF {
		s.fullScaned = true
		return Token{"", EOF}, pos
	}

	text := ""
	isNum := false
	switch Kind(ch) {
	case LETTER_LIT:
		for ch != "\n" && ch != "\t" && ch != " " && err != io.EOF {
			if ch == Keywords[LPAREN] || ch == Keywords[RPAREN] || ch == Keywords[COMMA] {
				s.undoCh()
				break
			}
			text += ch
			ch, err = s.nextCh()
		}
	case DIGIT_LIT:
		isNum = true
		for ch != " " && err != io.EOF {
			text += ch
			if Kind(ch) == LETTER_LIT {
				panic("Invalid variable name: " + text)
			}

			ch, err = s.nextCh()
			if kind := Kind(ch); kind != DIGIT_LIT && kind != DOT_LIT {
				s.undoCh()
				break
			}
		}
	case COMMA_LIT:
		text += ch
	default: // Operator
		for ch != " " && ch != "\n" && err != io.EOF {
			text += ch

			ch, err = s.nextCh()
			if Kind(ch) != OTHER_LIT {
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

func ToToken(token string, num bool) Token {
	if num {
		if strings.Contains(token, ".") {
			return Token{token, DOUBLE_LIT}
		} else {
			return Token{token, INT_LIT}
		}
	}

	kind := KeywordType(token)
	return Token{token, kind}
}

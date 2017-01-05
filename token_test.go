package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeywordType(t *testing.T) {
	assert.Equal(t, LPAREN, KeywordType("("))
	assert.Equal(t, INT, KeywordType("int"))
	assert.Equal(t, IDENT, KeywordType("hello"))
	assert.Equal(t, COMMENT_SLASH, KeywordType("//"))
}

func TestPriority(t *testing.T) {
	assert.Equal(t, Token{kind: PLUS}.Priority(), Token{kind: MINUS}.Priority())
	assert.True(t, (Token{kind: MINUS}.Priority() < Token{kind: MULTI}.Priority()))
}

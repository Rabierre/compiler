package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeywordType(t *testing.T) {
	assert.Equal(t, LPAREN, KeywordType("("))
	assert.Equal(t, INT, KeywordType("int"))
	assert.Equal(t, IDENT, KeywordType("hello"))
	assert.Equal(t, COMMENT, KeywordType("//"))
}

func TestPriority(t *testing.T) {
	assert.Equal(t, Token{Kind: PLUS}.Priority(), Token{Kind: MINUS}.Priority())
	assert.True(t, (Token{Kind: MINUS}.Priority() < Token{Kind: MULTI}.Priority()))
}

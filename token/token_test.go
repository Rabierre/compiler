package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeywordType(t *testing.T) {
	assert.Equal(t, FUNC, KeywordType("func"))
	assert.Equal(t, LPAREN, KeywordType("("))
	assert.Equal(t, INT, KeywordType("int"))
	assert.Equal(t, IDENT, KeywordType("hello"))
	assert.Equal(t, COMMENT, KeywordType("//"))
}

func TestPriority(t *testing.T) {
	assert.True(t, PLUS.Priority() == MINUS.Priority())
	assert.True(t, MINUS.Priority() < MULTI.Priority())
}

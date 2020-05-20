package punct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPureWord(t *testing.T) {
	w, b, e := PureWord("mama.")
	assert.Equal(t, "mama", w)
	assert.Equal(t, "", b)
	assert.Equal(t, ".", e)
	w, b, e = PureWord("mama")
	assert.Equal(t, "mama", w)
	assert.Equal(t, "", b)
	assert.Equal(t, "", e)
}

func TestPureWordStartEnd(t *testing.T) {
	w, b, e := PureWord("(,mama.")
	assert.Equal(t, "mama", w)
	assert.Equal(t, "(,", b)
	assert.Equal(t, ".", e)
	w, b, e = PureWord("(-mama")
	assert.Equal(t, "mama", w)
	assert.Equal(t, "(-", b)
	assert.Equal(t, "", e)
}

func TestPureWord_NoChange(t *testing.T) {
	w, b, e := PureWord("„mama“")
	assert.Equal(t, "„mama“", w)
	assert.Equal(t, "", b)
	assert.Equal(t, "", e)
}

func TestTrimQuotes(t *testing.T) {
	s1, s2, s3 := TrimWord("<mama>", IsParentheses)
	test3(t, s1, s2, s3, "mama", "<", ">")
	s1, s2, s3 = TrimWord("<mama>", IsPunct)
	test3(t, s1, s2, s3, "<mama>", "", "")
	s1, s2, s3 = TrimWord(",..<mama>,-+", IsAllSep)
	test3(t, s1, s2, s3, "mama", ",..<", ">,-+")
}

func test3(t *testing.T, s1, s2, s3, e1, e2, e3 string) {
	assert.Equal(t, e1, s1)
	assert.Equal(t, e2, s2)
	assert.Equal(t, e3, s3)
}

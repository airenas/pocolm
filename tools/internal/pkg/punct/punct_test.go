package punct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPureWord(t *testing.T) {
	w, e := PureWord("mama.")
	assert.Equal(t, "mama", w)
	assert.Equal(t, ".", e)
	w, e = PureWord("mama")
	assert.Equal(t, "mama", w)
	assert.Equal(t, "", e)
	w, e = PureWord("„mama“.,")
	assert.Equal(t, "mama", w)
	assert.Equal(t, ".,", e)
}

func TestPureWord_NoChange(t *testing.T) {
	w, e := PureWord("<mama>")
	assert.Equal(t, "<mama>", w)
	assert.Equal(t, "", e)
}

package punct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPureWord(t *testing.T) {
	assert.Equal(t, "mama", PureWord("mama."))
	assert.Equal(t, "mama", PureWord("„mama“."))
}

func TestPureWord_NoChange(t *testing.T) {
	assert.Equal(t, "<mama>", PureWord("<mama>"))
}

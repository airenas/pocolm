package lema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	assert.True(t, IsNumber("0.0"))
	assert.True(t, IsNumber("0"))
	assert.True(t, IsNumber("0."))
	assert.True(t, IsNumber("10.0"))
	assert.True(t, IsNumber("-500.0"))
	assert.True(t, IsNumber("+4540.0"))
	assert.True(t, IsNumber("+4540,0,0"))
}

func TestNotNumber(t *testing.T) {
	assert.False(t, IsNumber(""))
	assert.False(t, IsNumber(","))
	assert.False(t, IsNumber("a0"))
	assert.False(t, IsNumber("2 10"))
	assert.False(t, IsNumber("x500.0"))
}

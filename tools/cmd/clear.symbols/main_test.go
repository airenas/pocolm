package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	assert.Equal(t, "mama", changeLine("mama"))
}

func TestChange(t *testing.T) {
	assert.Equal(t, "mama", changeLine("mama\t"))
	assert.Equal(t, "mama mama", changeLine("mama\tmama"))
	assert.Equal(t, "mama - mama", changeLine("mama – mama"))
}

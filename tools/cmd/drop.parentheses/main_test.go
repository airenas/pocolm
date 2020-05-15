package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	assert.Equal(t, "mama ir olia", changeLine("mama ir olia"))
}

func TestDrops(t *testing.T) {
	assert.Equal(t, "mama <SKL>", changeLine("mama (olia)"))
	assert.Equal(t, "mama <SKL> xx", changeLine("mama (tatata hhha) xx"))
	assert.Equal(t, "<SKL> olia", changeLine("(mama) olia"))
	assert.Equal(t, "mama <SKL> dddd <SKL> ccc", changeLine("mama () dddd (aaa aaa) ccc"))
}

func TestDropsVarious(t *testing.T) {
	assert.Equal(t, "mama <SKL>", changeLine("mama [olia]"))
	assert.Equal(t, "mama <SKL>", changeLine("mama {olia}"))
	assert.Equal(t, "mama <SKL>", changeLine("mama <olia>"))
}

func TestNoDrop(t *testing.T) {
	assert.Equal(t, "mama (olia]", changeLine("mama (olia]"))
}

func TestDropsMultiple(t *testing.T) {
	assert.Equal(t, "mama <SKL> tt", changeLine("mama (olia (xxx) ir) tt"))
	assert.Equal(t, "mama ( <SKL> tt", changeLine("mama ((olia (xxx) ir) tt"))
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	assert.Equal(t, "mama ir olia", changeLine("mama ir olia"))
}

func TestDrops(t *testing.T) {
	assert.Equal(t, "mama <SKLIAUSTUOSE>", changeLine("mama (olia)"))
	assert.Equal(t, "mama <SKLIAUSTUOSE> xx", changeLine("mama (tatata hhha) xx"))
	assert.Equal(t, "<SKLIAUSTUOSE> olia", changeLine("(mama) olia"))
	assert.Equal(t, "mama <SKLIAUSTUOSE> dddd <SKLIAUSTUOSE> ccc", changeLine("mama () dddd (aaa aaa) ccc"))
}

func TestDropsVarious(t *testing.T) {
	assert.Equal(t, "mama <SKLIAUSTUOSE>", changeLine("mama [olia]"))
	assert.Equal(t, "mama <SKLIAUSTUOSE>", changeLine("mama {olia}"))
	assert.Equal(t, "mama <SKLIAUSTUOSE>", changeLine("mama <olia>"))
}

func TestNoDrop(t *testing.T) {
	assert.Equal(t, "mama (olia]", changeLine("mama (olia]"))
}

func TestDropsMultiple(t *testing.T) {
	assert.Equal(t, "mama <SKLIAUSTUOSE> tt", changeLine("mama (olia (xxx) ir) tt"))
	assert.Equal(t, "mama ( <SKLIAUSTUOSE> tt", changeLine("mama ((olia (xxx) ir) tt"))
}

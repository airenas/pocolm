package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	initRegexp()
	assert.Equal(t, "mama", changeLine("mama"))
}

func TestURL(t *testing.T) {
	initRegexp()
	assert.Equal(t, "mama <URL>", changeLine("mama http://www.delfi.lt"))
	assert.Equal(t, "mama <URL> ir", changeLine("mama http://delfi.lt?olia?tatata=tatat  ir"))
	assert.Equal(t, "mama <URL>", changeLine("mama http://www.delfi.lt"))
}

func TestEmail(t *testing.T) {
	initRegexp()
	assert.Equal(t, "mama <EMAIL>", changeLine("mama a@a.lt"))
	assert.Equal(t, "mama <EMAIL> ir", changeLine("mama ai@delfi.lt  ir"))
	assert.Equal(t, "<EMAIL> olia", changeLine("aaa@adelfi.lt olia"))
}

func TestUnderscore(t *testing.T) {
	initRegexp()
	assert.Equal(t, "mama <PILDOMA> <PILDOMA>", changeLine("mama __ _______"))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama________"))
}

func TestDot(t *testing.T) {
	initRegexp()
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama ...."))
	assert.Equal(t, "mama..", changeLine("mama.."))
	assert.Equal(t, "mama...", changeLine("mama..."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama ...."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama . .."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama ... ."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama . . . . .. ..."))
}

func TestNumberLetter(t *testing.T) {
	initRegexp()
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("a) mama"))
	assert.Equal(t, "aa) mama", changeLine("aa) mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("E) mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("a. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("B. mama"))
}

func TestNumber(t *testing.T) {
	initRegexp()
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("1. mama"))
	assert.Equal(t, "2005. mama", changeLine("2005. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("1.2.3. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("10.20.30. mama"))
}

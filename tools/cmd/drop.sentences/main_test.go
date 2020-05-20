package main

import (
	"strings"
	"testing"

	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/stretchr/testify/assert"
)

var lmCache *lema.Cache

func initCache(t *testing.T) {
	var err error
	lmCache, err = lema.NewTestCache(strings.NewReader(`
mama wR
Vilnius wPL
ir wR
kas wR
oliaa w
sutartis wR
h wA-h-
Tell. wA-Tell.-	
`))
	assert.Nil(t, err)
}

func TestNoChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "mama ir Vilnius", changeLine("mama ir Vilnius", lmCache))
	assert.NotEqual(t, "", changeLine("mama ir Vilnius", lmCache))
	assert.NotEqual(t, "", changeLine("mama, ir Vilnius.", lmCache))
	assert.NotEqual(t, "", changeLine("mama, ir Vilnius <PILDOMA>.", lmCache))
}

func TestDrop(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", changeLine("oliaa ir kas", lmCache))
}

func TestDropNonLtLetter(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", changeLine("ሀsutartis ir ir ir ir ir ir ", lmCache))
}

func TestDropJustNumbers(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", changeLine("0.0 1,000", lmCache))
}

func TestDropJustSpecial(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", changeLine("<PILDOMA>", lmCache))
}

func TestDropJustAbbreviation(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", changeLine("Tell.", lmCache))
	assert.Equal(t, "", changeLine("h", lmCache))
}

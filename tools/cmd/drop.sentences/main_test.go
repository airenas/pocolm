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
mama wL
Vilnius wPL
ir wL
kas w
olia w
sutartis w
KAUNAS wPR	
`))
	assert.Nil(t, err)
}

func TestNoChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "mama ir Vilnius", changeLine("mama ir Vilnius", lmCache))
	assert.Equal(t, "mama ir Vilnius", changeLine("mama ir Vilnius", lmCache))
	assert.Equal(t, "mama, ir Vilnius.", changeLine("mama, ir Vilnius.", lmCache))
}

func TestDrop(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", changeLine("olia ir kas", lmCache))
}

func TestDropNonLtLetter(t *testing.T) {
	assert.Equal(t, "", changeLine("ሀsutartis ir ir ir ir ir ir ", lmCache))
}

func TestDropJustNumbers(t *testing.T) {
	assert.Equal(t, "", changeLine("0.0 1,000", lmCache))
}

func TestDropJustSpecial(t *testing.T) {
	assert.Equal(t, "", changeLine("<PILDOMA>", lmCache))
}

func TestDropJustAbbreviation(t *testing.T) {
	assert.Equal(t, "", changeLine("h", lmCache))
	assert.Equal(t, "", changeLine("Tel.", lmCache))
}

func TestDropParenthesis(t *testing.T) {
	assert.Equal(t, "", changeLine("(olia ir", lmCache))
}

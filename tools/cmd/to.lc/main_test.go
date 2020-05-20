package main

import (
	"strings"
	"testing"

	"github.com/airenas/pocolm/tools/internal/pkg/abbr"
	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/stretchr/testify/assert"
)

var lmCache *lema.Cache

var ad *abbr.Abbreviations

func initCache(t *testing.T) {
	var err error
	lmCache, err = lema.NewTestCache(strings.NewReader(`
mama wR
Vilnius wP
Ir wR
Kas wR
Olia wR
sutartis w
KAUNAS wPR	
SUTARTIS wR
Kass w
LRT wA-LRT-
VDU. wA-Vdu-
`))
	assert.Nil(t, err)
	ad = newTestAbbreviation(t)
}

func TestNoChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "mama ir Vilnius", changeLine("mama Ir Vilnius", lmCache, ad))
}

func TestChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "olia ir kas", changeLine("Olia Ir Kas", lmCache, ad))
}

func TestQuoted(t *testing.T) {
	initCache(t)
	assert.Equal(t, "„Olia“ ir „Vilnius“", changeLine("„Olia“ Ir „VILNIUS“", lmCache, ad))
}

func TestStartWithPunct(t *testing.T) {
	initCache(t)
	assert.Equal(t, ".„Olia“ ,.ir (Vilnius", changeLine(".„Olia“ ,.Ir (VILNIUS", lmCache, ad))
}

func TestQuotedLeaveUpper(t *testing.T) {
	initCache(t)
	assert.Equal(t, "„Olia“ ir „Kaunas“", changeLine("„Olia“ Ir „KAUNAS“", lmCache, ad))
	assert.Equal(t, "olia ir Kaunas", changeLine("Olia Ir KAUNAS", lmCache, ad))
}

func TestChangeAll(t *testing.T) {
	initCache(t)
	assert.Equal(t, "sutartis Vilnius", changeLine("SUTARTIS VILNIUS", lmCache, ad))
}

func TestLeavesUnknown(t *testing.T) {
	initCache(t)
	assert.Equal(t, "Kass", changeLine("Kass", lmCache, ad))
}

func TestAbbreviationFile(t *testing.T) {
	initCache(t)
	assert.Equal(t, "olia. pp.y.", changeLine("Olia. Pp.Y.", lmCache, ad))
}

func TestAbbreviationLema(t *testing.T) {
	initCache(t)
	assert.Equal(t, "olia. LRT.,", changeLine("Olia. LRT.,", lmCache, ad))
	assert.Equal(t, "olia. Vdu.,", changeLine("Olia. VDU.,", lmCache, ad))
}

func newTestAbbreviation(t *testing.T) *abbr.Abbreviations {
	res, err := abbr.NewAbbrReader(strings.NewReader(`
t. y.
pp.y.	
`))
	assert.Nil(t, err)
	return res
}

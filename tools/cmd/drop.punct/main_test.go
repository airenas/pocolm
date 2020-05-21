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
VILNIUS wP
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
	assert.Equal(t, "mama Ir Vilnius", changeLine("mama Ir Vilnius", lmCache, ad))
}

func TestChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "a b", changeLine("a )- b", lmCache, ad))
	assert.Equal(t, "Olia Ir Kas", changeLine("Olia Ir Kas...", lmCache, ad))
}

func TestQuoted(t *testing.T) {
	initCache(t)
	assert.Equal(t, "„Olia“ Ir „VILNIUS“", changeLine("„Olia“ Ir „VILNIUS“,!?", lmCache, ad))
}

func TestStartWithPunct(t *testing.T) {
	initCache(t)
	assert.Equal(t, "„Olia“ Ir Vilnius", changeLine(".„Olia“ ,.Ir (((Vilnius", lmCache, ad))
}

func TestNumbers(t *testing.T) {
	initCache(t)
	assert.Equal(t, "132 -1,45.24", changeLine("132/ -1,45.24.,", lmCache, ad))
	assert.Equal(t, "132% 132", changeLine("132%., 132.%", lmCache, ad))
	assert.Equal(t, "-132 +123,123", changeLine(",/,-132  -=+123,123.", lmCache, ad))
	assert.Equal(t, "±132", changeLine("±132", lmCache, ad))
}

func TestChangeVarious(t *testing.T) {
	initCache(t)
	assert.Equal(t, "sutartis priedas", changeLine("sutartis | priedas", lmCache, ad))
	assert.Equal(t, "sutartis Vilnius", changeLine("* - sutartis - Vilnius %)()<>{}[]", lmCache, ad))
}

func TestAbbreviationFile(t *testing.T) {
	initCache(t)
	assert.Equal(t, "olia pp.y.", changeLine("olia. pp.y.", lmCache, ad))
}

func TestAbbreviationLema(t *testing.T) {
	initCache(t)
	assert.Equal(t, "olia LRT", changeLine("olia. LRT.,", lmCache, ad))
}

func newTestAbbreviation(t *testing.T) *abbr.Abbreviations {
	res, err := abbr.NewAbbrReader(strings.NewReader(`
t. y.
pp.y.	
`))
	assert.Nil(t, err)
	return res
}

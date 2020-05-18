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
mama w
Vilnius wP
ir w
kas w
olia w
sutartis w
KAUNAS wP	
`))
	assert.Nil(t, err)
}

func TestNoChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "mama ir Vilnius", changeLine("mama ir Vilnius", lmCache))
}

func TestChange(t *testing.T) {
	initCache(t)
	assert.Equal(t, "olia ir kas", changeLine("Olia Ir Kas", lmCache))
}

func TestQuoted(t *testing.T) {
	initCache(t)
	assert.Equal(t, "„olia“ ir „Vilnius“", changeLine("„Olia“ Ir „VILNIUS“", lmCache))
}

func TestQuotedLeaveUpper(t *testing.T) {
	initCache(t)
	assert.Equal(t, "„olia“ ir „Kaunas“", changeLine("„Olia“ Ir „KAUNAS“", lmCache))
	assert.Equal(t, "olia ir kaunas", changeLine("Olia Ir KAUNAS", lmCache))
}

func TestChangeAll(t *testing.T) {
	initCache(t)
	assert.Equal(t, "sutartis Vilnius", changeLine("SUTARTIS VILNIUS", lmCache))
}

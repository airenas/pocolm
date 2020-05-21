package lema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lmCache *Cache

func initCache(t *testing.T) {
	var err error
	lmCache, err = NewTestCache(strings.NewReader(`
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

func TestEmptyNoFail(t *testing.T) {
	initCache(t)
	assert.Equal(t, "", lmCache.AbbreviationString(""))
	assert.Equal(t, false, lmCache.Regular(""))
	assert.Equal(t, false, lmCache.Proper(""))
}

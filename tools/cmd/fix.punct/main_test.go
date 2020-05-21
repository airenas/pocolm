package main

import (
	"strings"
	"testing"

	"github.com/airenas/pocolm/tools/internal/pkg/abbr"
	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	assert.Equal(t, "mama, olia.", changePunct("mama, olia."))
}

func TestFixes(t *testing.T) {
	assert.Equal(t, "mama, olia", changePunct("mama,olia"))
	assert.Equal(t, "mama. olia", changePunct("mama.olia"))
	assert.Equal(t, "mama: olia", changePunct("mama:olia"))
	assert.Equal(t, "mama - olia", changePunct("mama-olia"))
	assert.Equal(t, "mama (olia", changePunct("mama(olia"))
	assert.Equal(t, "1) olia", changePunct("1)olia"))
}

func TestLeavesWithNumber(t *testing.T) {
	assert.Equal(t, "1-olia", changePunct("1-olia"))
	assert.Equal(t, "olia-2", changePunct("olia-2"))
}

func TestAbbreviations(t *testing.T) {
	ad := newTestAbbreviation(t)
	assert.Equal(t, "a t.y. olia", changeLine("a t.y. olia", ad))
	assert.Equal(t, "a. T.y. olia", changeLine("a. T.y. olia", ad))
	assert.Equal(t, "a t.y. olia", changeLine("a t. y. olia", ad))
}

func TestQuetedShort(t *testing.T) {
	ad := newTestAbbreviation(t)
	assert.Equal(t, "olia „P“", changeLine("olia „P“", ad))
}

func TestLeavesNumbers(t *testing.T) {
	assert.Equal(t, "2015.12.12", changePunct("2015.12.12"))
	assert.Equal(t, "2015-12-12", changePunct("2015-12-12"))
	assert.Equal(t, "2015,12,12", changePunct("2015,12,12"))
}

func newTestAbbreviation(t *testing.T) *abbr.Abbreviations {
	res, err := abbr.NewAbbrReader(strings.NewReader(`
t. y.
p	
`))
	assert.Nil(t, err)
	return res
}

func benchmarkRegexp(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		changePunct(s)
	}
}

func BenchmarkChange(b *testing.B) {
	benchmarkRegexp(b, "a.asdsad-dasdasd das dsad ds das")
}

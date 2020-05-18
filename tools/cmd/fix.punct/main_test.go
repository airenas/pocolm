package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	assert.Equal(t, "mama, olia.", changeLine("mama, olia."))
}

func TestFixes(t *testing.T) {
	assert.Equal(t, "mama, olia", changeLine("mama,olia"))
	assert.Equal(t, "mama. olia", changeLine("mama.olia"))
	assert.Equal(t, "mama: olia", changeLine("mama:olia"))
	assert.Equal(t, "mama - olia", changeLine("mama-olia"))
}

func TestLeavesInQuotes(t *testing.T) {
	assert.Equal(t, "„mama-olia“", changeLine("„mama-olia“"))
}

func TestLeavesNumbers(t *testing.T) {
	assert.Equal(t, "2015.12.12", changeLine("2015.12.12"))
	assert.Equal(t, "2015-12-12", changeLine("2015-12-12"))
	assert.Equal(t, "2015,12,12", changeLine("2015,12,12"))
}

func benchmarkRegexp(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		changeLine(s)
	}
}

func BenchmarkChange(b *testing.B) {
	benchmarkRegexp(b, "a.asdsad-dasdasd das dsad ds das")
}

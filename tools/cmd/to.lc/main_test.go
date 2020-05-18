package main

import (
	"testing"

	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	lm := lema.NewCache()
	defer lm.Close()
	assert.Equal(t, "mama ir Petras", changeLine("mama ir Petras", lm))
}

func TestChange(t *testing.T) {
	lm := lema.NewCache()
	defer lm.Close()
	assert.Equal(t, "olia ir kas", changeLine("Olia Ir Kas", lm))
}

func TestQuoted(t *testing.T) {
	lm := lema.NewCache()
	defer lm.Close()
	assert.Equal(t, "„olia“ ir „Vilnius“", changeLine("„Olia“ Ir „VILNIUS“", lm))
}

func TestChangeAll(t *testing.T) {
	lm := lema.NewCache()
	defer lm.Close()
	assert.Equal(t, "sutartis Vilnius", changeLine("SUTARTIS VILNIUS", lm))
}

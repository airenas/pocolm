package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var vocab map[string]bool

func initTestVocab(t *testing.T) {
	var err error
	vocab, err = loadVocabReader(strings.NewReader(`
olia 1
0 2
+ 3
tata 4
`))
	assert.Nil(t, err)
}

func TestFind(t *testing.T) {
	initTestVocab(t)
	nl, err := changeLine("10 olia", vocab)
	assert.Nil(t, err)
	assert.Equal(t, "10 olia", nl)
}

func TestNotFind(t *testing.T) {
	initTestVocab(t)
	nl, err := changeLine("10 olia2", vocab)
	assert.Nil(t, err)
	assert.Equal(t, "", nl)
}

func TestFail(t *testing.T) {
	initTestVocab(t)
	_, err := changeLine("10", vocab)
	assert.NotNil(t, err)
}

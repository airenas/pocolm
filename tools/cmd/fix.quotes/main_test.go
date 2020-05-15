package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	assert.Equal(t, "mama", changeLine("mama"))
}

func TestDrops(t *testing.T) {
	assert.Equal(t, "mama olia", changeLine("mama    olia"))
	assert.Equal(t, "mama tatata", changeLine("mama tatata"))
	assert.Equal(t, "mama tatata", changeLine("mama ,,tatata"))
	assert.Equal(t, "mama tatata", changeLine("\"mama tatata\""))
	assert.Equal(t, "mama, tatata", changeLine("\"mama, tatata\"     "))
}

func TestChages(t *testing.T) {
	assert.Equal(t, "„mama“", changeLine("\"mama\""))
	assert.Equal(t, "„mama“", changeLine(",,mama\""))
	assert.Equal(t, "„mama“", changeLine(",,mama''"))
}

func TestChangesWithSep(t *testing.T) {
	assert.Equal(t, "„mama“.", changeLine(",,mama\"."))
	assert.Equal(t, ".„mama“", changeLine(".,,mama\""))
	assert.Equal(t, ".„mama“.", changeLine(".,,mama\"."))
	assert.Equal(t, ".„mama“]>*&^.,-", changeLine(".,,mama\"]>*&^.,-"))
	assert.Equal(t, "„mama.mama“.", changeLine(",,mama.mama\"."))
}

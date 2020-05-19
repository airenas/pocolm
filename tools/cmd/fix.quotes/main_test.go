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
	assert.Equal(t, "Mama tatata", changeLine("\"Mama tatata\""))
	assert.Equal(t, "Mama, tatata", changeLine("\"Mama, tatata\"     "))
}

func TestChanges(t *testing.T) {
	assert.Equal(t, "„Mama“", changeLine("\"Mama\""))
	assert.Equal(t, "„Mama“", changeLine(",,Mama\""))
	assert.Equal(t, "„Mama“", changeLine(",,Mama''"))
}

func TestChangesWithSep(t *testing.T) {
	assert.Equal(t, "„Mama“.", changeLine(",,Mama\"."))
	assert.Equal(t, ".„Mama“", changeLine(".,,Mama\""))
	assert.Equal(t, ".„Mama“.", changeLine(".,,Mama\"."))
	assert.Equal(t, ".„Mama“]>*&^.,-", changeLine(".,,Mama\"]>*&^.,-"))
}

func TestDropsWithSepInside(t *testing.T) {
	assert.Equal(t, "Mama.mama.", changeLine(",,Mama.mama\"."))
	assert.Equal(t, "Mama-mama", changeLine("\"Mama-mama\""))
	assert.Equal(t, "Mama,mama", changeLine("\"Mama,mama\""))
	assert.Equal(t, "Mama,", changeLine("\"Mama,\""))
}

func TestNumbersLeaves(t *testing.T) {
	assert.Equal(t, "„Mama1231“", changeLine("\"Mama1231\""))
}

func TestDropsLower(t *testing.T) {
	assert.Equal(t, "mama", changeLine(",,mama\""))
}
func TestSepSpecialCase(t *testing.T) {
	assert.Equal(t, "„A.“.", changeLine(",,A.\"."))
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeFileName(t *testing.T) {
	assert.Equal(t, "olia/000001.txt", makeFileName("olia", 1))
	assert.Equal(t, "olia/012345.txt", makeFileName("olia", 12345))
}

func TestTitleFix(t *testing.T) {
	 assert.Equal(t, "a.\nb", getData(&article{Title: "a " , Body: "b"}))
	 assert.Equal(t, "a.\nb", getData(&article{Title: "a. ", Body: "b"}))
	 assert.Equal(t, "a?\nb", getData(&article{Title: "a?  ", Body: "b"}))
	 assert.Equal(t, "a!\nb", getData(&article{Title: "a!  ", Body: "b"}))
	 assert.Equal(t, "a,.\nb", getData(&article{Title: "a,  ", Body: "b"}))
	 assert.Equal(t, "b", getData(&article{Title: "", Body: "b"}))
}

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixSpaces(t *testing.T) {
	assert.Equal(t, "mama olia tata", FixSpaces(" mama  olia          tata "))
	assert.Equal(t, "mama olia,.", FixSpaces(" mama  olia,."))
}

func benchmarkRegexp(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		s = FixSpacesR(s)
	}
}

func benchmarkLoop(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		s = FixSpaces(s)
	}
}

func BenchmarkReg(b *testing.B) {
	benchmarkRegexp(b, "a.     asdsad          dasdasd   ")
}

func BenchmarkLoop(b *testing.B) {
	benchmarkLoop(b, "a.     asdsad          dasdasd   ")
}

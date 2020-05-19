package abbr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	ad := newTestAbbreviation(t, `t. y.`)
	assert.NotNil(t, ad)
	i, s := ad.Find("t. y.")
	assert.Equal(t, 5, i)
	assert.Equal(t, "t.y.", s)
	i, s = ad.Find("t.y.")
	assert.Equal(t, 4, i)
	assert.Equal(t, "t.y.", s)
	i, s = ad.Find("T.   y.")
	assert.Equal(t, 7, i)
	assert.Equal(t, "t.y.", s)
}

func TestNotFind(t *testing.T) {
	ad := newTestAbbreviation(t, `t. y.
t. kas`)
	assert.NotNil(t, ad)
	i, _ := ad.Find("tt. y.")
	assert.Equal(t, 0, i)
	i, _ = ad.Find("t.")
	assert.Equal(t, 0, i)
	i, _ = ad.Find("t.y")
	assert.Equal(t, 0, i)
	i, _ = ad.Find("t.yy.")
	assert.Equal(t, 0, i)
	i, _ = ad.Find("t. kass")
	assert.Equal(t, 0, i)
	i, _ = ad.Find("t. kas0")
	assert.Equal(t, 0, i)
}

func TestFind(t *testing.T) {
	ad := newTestAbbreviation(t, `t. y.
t. kas`)
	i, _ := ad.Find("t. kas,")
	assert.Equal(t, 6, i)
	i, _ = ad.Find("t. kas(),")
	assert.Equal(t, 6, i)
}

func TestShort(t *testing.T) {
	ad := newTestAbbreviation(t, `t. kas. tas.
t. kas
t.`)
	i, _ := ad.Find("t. kass,")
	assert.Equal(t, 2, i)
	i, _ = ad.Find("t. kas. tass.,")
	assert.Equal(t, 6, i)
}

func TestSlash(t *testing.T) {
	ad := newTestAbbreviation(t, `t/kg
t/g
t. g.
`)
	i, _ := ad.Find("t / kg")
	assert.Equal(t, 6, i)
	i, _ = ad.Find("t/g.")
	assert.Equal(t, 3, i)
	i, _ = ad.Find("t. g.")
	assert.Equal(t, 5, i)
}

func newTestAbbreviation(t *testing.T, s string) *Abbreviations {
	res, err := NewAbbrReader(strings.NewReader(s))
	assert.Nil(t, err)
	return res
}

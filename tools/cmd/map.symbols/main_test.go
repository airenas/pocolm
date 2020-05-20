package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChange(t *testing.T) {
	initRegexp(&params{})
	assert.Equal(t, "mama", changeLine("mama"))
}

func TestURL(t *testing.T) {
	initRegexp(&params{changeURL: true})
	assert.Equal(t, "mama <URL>", changeLine("mama http://www.delfi.lt"))
	assert.Equal(t, "mama <URL> ir", changeLine("mama http://delfi.lt?olia?tatata=tatat  ir"))
	assert.Equal(t, "mama www.delfi.lt", changeLine("mama www.delfi.lt"))
}

func TestEmail(t *testing.T) {
	initRegexp(&params{changeEmail: true})
	assert.Equal(t, "mama <EMAIL>", changeLine("mama a@a.lt"))
	assert.Equal(t, "mama <EMAIL> ir", changeLine("mama ai@delfi.lt  ir"))
	assert.Equal(t, "<EMAIL> olia", changeLine("aaa@adelfi.lt olia"))
}

func TestEmailFix(t *testing.T) {
	initRegexp(&params{changeEmail: true})
	assert.Equal(t, "mama <EMAIL>", changeLine("mama mailto:a@a.lt"))
}

func TestUnderscore(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "mama <PILDOMA> <PILDOMA>", changeLine("mama __ _______"))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama________"))
}

func TestDot(t *testing.T) {
	initRegexp(&params{})
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama ...."))
	assert.Equal(t, "mama..", changeLine("mama.."))
	assert.Equal(t, "mama...", changeLine("mama..."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama ...."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama . .."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama ... ."))
	assert.Equal(t, "mama <PILDOMA>", changeLine("mama . . . . .. ..."))
}

func TestNumberLetter(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("a) mama"))
	assert.Equal(t, "aa) mama", changeLine("aa) mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("E) mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("a. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("B. mama"))
}

func TestNumber(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("1. mama"))
	assert.Equal(t, "2005. mama", changeLine("2005. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("1.2.3. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("10.20.30. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("1) mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("500.300.30. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("500) mama"))
}

func TestStar(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<PUNKTAS> mama", changeLine("* mama"))
}

func TestDash(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<PUNKTAS> mama", changeLine("- mama"))
}

func TestNumberRoman(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("IX. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("IX) mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("ix. mama"))
	assert.Equal(t, "<NUMERACIJA> mama", changeLine("ix) mama"))
}

func TestNumberRomanPart(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<NUMERACIJA> skyrius. mama", changeLine("IX skyrius. mama"))
}

func TestSlash(t *testing.T) {
	initRegexp(&params{changeSlash: true})
	assert.Equal(t, "buvo ir gerai", changeLine("buvo ir/ar gerai"))
	assert.Equal(t, "buvo ir gerai", changeLine("buvo ir/ar/gal gerai"))
	assert.Equal(t, "buvo iš gerai", changeLine("buvo iš/į/gal gerai"))
	assert.Equal(t, "buvo iš gerai", changeLine("buvo iš/į gerai"))
	assert.Equal(t, "buvo, labai gerai", changeLine("buvo,labai/įgerai/gal gerai"))
}

func TestSlashCase(t *testing.T) {
	initRegexp(&params{changeSlash: true})
	assert.Equal(t, "buvo Ir gerai", changeLine("buvo Ir/ar gerai"))
	assert.Equal(t, "buvo ir gerai", changeLine("buvo ir/Ar/gal gerai"))
}

func TestSlashIgnore(t *testing.T) {
	initRegexp(&params{changeSlash: true})
	assert.Equal(t, "buvo ir/r1 gerai", changeLine("buvo ir/r1 gerai"))
	assert.Equal(t, "buvo 1ir/r gerai", changeLine("buvo 1ir/r gerai"))
	assert.Equal(t, "buvo 012/ir gerai", changeLine("buvo 012/ir gerai"))
	assert.Equal(t, "buvo ir/01 gerai", changeLine("buvo ir/01 gerai"))
	assert.Equal(t, "buvo 22/01 gerai", changeLine("buvo 22/01 gerai"))
}

func TestAppendix(t *testing.T) {
	initRegexp(&params{changeNumeration: true})
	assert.Equal(t, "<NUMERACIJA> priedas. mama", changeLine("8 priedas. mama"))
}

func TestDropHyperlink(t *testing.T) {
	initRegexp(&params{})
	assert.Equal(t, "mama", changeLine("mama HYPERLINK"))
}

func TestDropFormtext(t *testing.T) {
	initRegexp(&params{})
	assert.Equal(t, "mama", changeLine("mama FORMTEXT"))
}

func TestDropPageRef(t *testing.T) {
	initRegexp(&params{})
	assert.Equal(t, "mama", changeLine("mama PAGEREF _Toc305158161 \\h 29"))
	assert.Equal(t, "mama", changeLine("mama PAGEREF _Toc305158161 \\h 2"))
	assert.Equal(t, "mama aaaa", changeLine("mama PAGEREF _Toc305158161 \\h 2 aaaa"))
}

func TestDropMergefield(t *testing.T) {
	initRegexp(&params{})
	assert.Equal(t, "mama", changeLine("mama MERGEFIELD „Turto_aprašas“"))
	assert.Equal(t, "mama", changeLine("mama MERGEFIELD Turto_aprašas"))
	assert.Equal(t, "mama aaaa", changeLine("mama MERGEFIELD „Turto_aprašas“ aaaa"))
	assert.Equal(t, "mama aaaa", changeLine("mama MERGEFIELD „Turto_22_aprašas“ aaaa"))
}

func benchmarkRegexp(b *testing.B, replaces []*replace, s string) {
	for i := 0; i < b.N; i++ {
		for _, rep := range replaces {
			s = rep.regxp.ReplaceAllString(s, rep.str)
		}
	}
}

func BenchmarkTwo(b *testing.B) {
	replaces = make([]*replace, 0)
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[A-Z]\\)")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[a-z]\\)")})
	benchmarkRegexp(b, replaces, "a. asdsad dasdasd das dsad ds das")
}

func BenchmarkOne(b *testing.B) {
	replaces = make([]*replace, 0)
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("(?i)^[A-Z]\\)")})
	benchmarkRegexp(b, replaces, "a. asdsad dasdasd das dsad ds das")
}

func BenchmarkFour(b *testing.B) {
	replaces = make([]*replace, 0)
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[A-Z]\\)")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[a-z]\\)")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[A-Z]\\.")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[a-z]\\.")})
	benchmarkRegexp(b, replaces, "a. asdsad dasdasd das dsad ds das")
}

func BenchmarkOneAll(b *testing.B) {
	replaces = make([]*replace, 0)
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("(?i)^[A-Z]\\[\\)\\.]")})
	benchmarkRegexp(b, replaces, "a. asdsad dasdasd das dsad ds das")
}

package lema

import (
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/punct"
)

//IsNumber return true if string is number with punctuation marks
func IsNumber(w string) bool {
	rns := []rune(w)
	dg := 0
	for _, r := range rns {
		if punct.IsPunct(r) {
			continue
		} else if unicode.IsDigit(r) {
			dg++
			continue
		}
		return false
	}
	return dg > 0
}

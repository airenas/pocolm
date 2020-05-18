package lema

import (
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/punct"
)

//IsNumber return true if string is number with punctuation marks
func IsNumber(w string) bool {
	rns := []rune(w)
	for _, r := range rns {
		if (unicode.IsDigit(r)) || punct.IsPunct(r) {
			continue
		}
		return false
	}
	return true
}

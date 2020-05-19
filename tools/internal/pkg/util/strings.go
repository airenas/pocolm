package util

import (
	"strings"
	"unicode"
)

//FixSpacesR removes double spaces and trims
func FixSpacesR(s string) string {
	return strings.TrimSpace(MultiSpacesRegexp.ReplaceAllString(s, " "))
}

//FixSpaces removes double spaces and trims
func FixSpaces(s string) string {
	res := strings.Builder{}
	res.Grow(len(s))
	ins := false
	for _, r := range s {
		if r == ' ' {
			ins = true
			continue
		}
		if ins && res.Len() > 0 {
			res.WriteRune(' ')
		}
		res.WriteRune(r)
		ins = false
	}
	return res.String()
}

//DropSpaces removes spaces
func DropSpaces(s string) string {
	res := strings.Builder{}
	rns := []rune(s)
	res.Grow(len(s))
	ins := false
	pr := rune(-1)
	for _, r := range rns {
		if r == ' ' {
			ins = ins || LetterDigit(pr)
			continue
		}
		pr = r
		if ins && res.Len() > 0 {
			res.WriteRune(' ')
		}
		res.WriteRune(r)
		ins = false
	}
	return res.String()
}

//LetterDigit returns true on letter digit
func LetterDigit(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

package util

import "strings"

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

package util

import "regexp"

//MultiSpacesRegexp to or more spaces regexp
var MultiSpacesRegexp *regexp.Regexp

//EMailRegexp to match email
var EMailRegexp *regexp.Regexp

//SpecialWordRegexp to match <xxx>
var SpecialWordRegexp *regexp.Regexp

func init() {
	MultiSpacesRegexp = regexp.MustCompile("[ ]{2,}")
	EMailRegexp = regexp.MustCompile("[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}")
	SpecialWordRegexp = regexp.MustCompile("<[\\p{L}_\\d]+>")
}

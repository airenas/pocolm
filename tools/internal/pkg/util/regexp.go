package util

import "regexp"

//MultiSpacesRegexp to or more spaces regexp
var MultiSpacesRegexp *regexp.Regexp

func init() {
	MultiSpacesRegexp, _ = regexp.Compile("[ ]{2,}")
}

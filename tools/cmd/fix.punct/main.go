package main

import (
	"strings"
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

func main() {
	cmd.ProcessByLine(processLine)
}

func processLine(line string) (string, error) {
	return changeLine(strings.TrimSpace(line)), nil
}

func changeLine(line string) string {
	if len(line) == 0 {
		return line
	}
	rns := []rune(line)
	res := strings.Builder{}
	res.Grow(len(rns))
	pr := '\n'
	for _, r := range rns {
		rLetter := unicode.IsLetter(r)
		if pr == '.' && rLetter {
			res.WriteRune(' ')
		} else if pr == ',' && rLetter {
			res.WriteRune(' ')
		} else if pr == ':' && rLetter {
			res.WriteRune(' ')
		} else if r == '-' && unicode.IsLetter(pr) {
			res.WriteRune(' ')
		} else if pr == '-' && rLetter {
			res.WriteRune(' ')
		}
		res.WriteRune(r)
		pr = r
	}
	return util.MultiSpacesRegexp.ReplaceAllString(res.String(), " ")
}

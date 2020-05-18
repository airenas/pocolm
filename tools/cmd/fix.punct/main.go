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
	pr := '\n'
	for _, r := range rns {
		if pr == '.' && unicode.IsLetter(r) {
			res.WriteRune(' ')
		}
		if pr == ',' && unicode.IsLetter(r) {
			res.WriteRune(' ')
		}
		if pr == ':' && unicode.IsLetter(r) {
			res.WriteRune(' ')
		}
		if r == '-' && unicode.IsLetter(pr) {
			res.WriteRune(' ')
		}
		if pr == '-' && unicode.IsLetter(r) {
			res.WriteRune(' ')
		}
		res.WriteRune(r)
		pr = r
	}
	return util.MultiSpacesRegexp.ReplaceAllString(res.String(), " ")
}

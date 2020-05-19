package main

import (
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

var replaceableSymbols map[rune][]rune

func init() {
	replaceableSymbols = make(map[rune][]rune)
	for _, r := range []rune(" \t\r•\uFEFF\x00\u007f") {
		replaceableSymbols[r] = []rune(" ")
	}
	replaceableSymbols['–'] = []rune("-")
	replaceableSymbols['…'] = []rune("...")
}

func main() {
	cmd.InitApp()
	cmd.ProcessByLine(processLine)
}

func processLine(line string) (string, error) {
	return changeLine(strings.TrimSpace(line)), nil
}

func changeLine(line string) string {
	if len(line) == 0 {
		return line
	}
	runes := []rune(line)
	res := make([]rune, 0)
	for _, r := range runes {
		res = append(res, changeSymbol(r)...)
	}
	r := string(res)
	r = util.MultiSpacesRegexp.ReplaceAllString(r, " ")
	return strings.TrimSpace(r)
}

func changeSymbol(r rune) []rune {
	s, ok := replaceableSymbols[r]
	if ok {
		return s
	}
	return []rune{r}
}

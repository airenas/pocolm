package main

import (
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
)

var quoteSymbols map[rune]bool
var sepSymbols map[rune]bool

func init() {
	quoteSymbols = make(map[rune]bool)
	for _, r := range []rune("\"“”„“„”") {
		quoteSymbols[r] = true
	}
	sepSymbols = make(map[rune]bool)
	for _, r := range []rune(",./\\;:[]{}<>()*&^%$#@!~?-+") {
		sepSymbols[r] = true
	}
}

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
	ln := strings.ReplaceAll(line, ",,", "\"")
	ln = strings.ReplaceAll(ln, "''", "\"")
	ln = strings.ReplaceAll(ln, "``", "\"")
	f := strings.Fields(ln)
	res := make([]string, 0)
	for _, s := range f {
		s1, w, s2 := extractWordSep(s)
		rns := []rune(w)
		q := quoted(rns)
		rns = trimQuotes(rns)
		if q && len(rns) > 0 {
			rns = quoteLt(rns)
		}
		res = append(res, s1+string(rns)+s2)
	}
	return strings.Join(res, " ")
}

func quoted(rns []rune) bool {
	l := len(rns)
	if l < 3 {
		return false
	}
	return quoteSymbols[rns[0]] && quoteSymbols[rns[l-1]]
}

func trimQuotes(rns []rune) []rune {
	for len(rns) > 0 && quoteSymbols[rns[0]] {
		rns = rns[1:]
	}
	for len(rns) > 0 && quoteSymbols[rns[len(rns)-1]] {
		rns = rns[:len(rns)-1]
	}
	return rns
}

func extractWordSep(s string) (string, string, string) {
	rns := []rune(s)
	i1, i2 := 0, len(rns)-1
	for ; i1 < len(rns) && sepSymbols[rns[i1]]; i1++ {
	}
	for ; i2 > i1 && sepSymbols[rns[i2]]; i2-- {
	}
	// s2 := ""
	// if i2 < len(rns) {
	// 	s2 = string(rns[i2:])
	// }
	return string(rns[:i1]), string(rns[i1 : i2+1]), string(rns[i2+1:])
}

func quoteLt(rns []rune) []rune {
	rns = append([]rune{punct.StartQuote}, rns...)
	return append(rns, punct.EndQuote)
}

package main

import (
	"strings"
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
)

var quoteSymbols map[rune]bool

func init() {
	quoteSymbols = make(map[rune]bool)
	for _, r := range []rune("\"“”„") {
		quoteSymbols[r] = true
	}
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
	ln := strings.ReplaceAll(line, ",,", "\"")
	ln = strings.ReplaceAll(ln, "''", "\"")
	ln = strings.ReplaceAll(ln, "``", "\"")
	f := strings.Fields(ln)
	res := make([]string, 0)
	for _, s := range f {
		w, s1, s2 := punct.PureWord(s)
		rns := []rune(w)
		q := quoted(rns)
		rns = trimQuotes(rns)
		if q && len(rns) > 0 && isOKForQuote(rns) {
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

func quoteLt(rns []rune) []rune {
	rns = append([]rune{punct.StartQuote}, rns...)
	return append(rns, punct.EndQuote)
}

func isOKForQuote(rns []rune) bool {
	for i, r := range rns {
		if i == 0 && unicode.IsLower(r) {
			return false
		}
		if i > 0 && i == (len(rns)-1) && r == '.' {
			return true
		}
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			return false
		}
	}
	return true
}

package main

import (
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

var parenthesesSymbols map[rune]rune

func init() {
	parenthesesSymbols = make(map[rune]rune)
	parenthesesSymbols['('] = ')'
	parenthesesSymbols['['] = ']'
	parenthesesSymbols['{'] = '}'
	parenthesesSymbols['<'] = '>'
	parenthesesSymbols['«'] = '»'
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
	rns := []rune(line)
	res := strings.Builder{}
	for i := 0; i < len(rns); i++ {
		p, f := parenthesesSymbols[rns[i]]
		if f {
			pi, f := findParenthesis(rns[i:], p)
			if f {
				res.WriteString(" <SKLIAUSTUOSE> ")
				i += pi
				continue
			}
		}
		res.WriteRune(rns[i])
	}
	return util.FixSpaces(rs)
}

func findParenthesis(rns []rune, p rune) (int, bool) {
	ps := rns[0]
	c := 0
	for i, r := range rns {
		if r == ps {
			c++
		}
		if r == p {
			c--
			if c == 0 {
				return i, true
			}
		}
	}
	return 0, false
}

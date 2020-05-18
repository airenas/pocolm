package main

import (
	"strings"
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

type lemaProper interface {
	IsProper(string) bool
}

func main() {
	lm := lema.NewCache()
	defer lm.Close()
	cmd.ProcessByLine(func(line string) (string, error) { return changeLine(line, lm), nil })
}

func changeLine(line string, lm lemaProper) string {
	if len(line) == 0 {
		return line
	}
	strs := strings.Split(line, " ")
	res := strings.Builder{}
	for _, w := range strs {
		if w != "" {
			if res.Len() > 0 {
				res.WriteString(" ")
			}
			res.WriteString(changeWord(w, lm))
		}
	}
	return res.String()
}

func changeWord(w string, lm lemaProper) string {
	wc := punct.PureWord(w)
	if wc == "" {
		return w
	}
	if util.SpecialWordRegexp.MatchString(wc) {
		return w
	}
	if lm.IsProper(wc) {
		return changeTitle(w)
	}
	return strings.ToLower(w)
}

func changeTitle(w string) string {
	r := []rune(strings.ToLower(w))
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

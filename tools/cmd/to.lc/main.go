package main

import (
	"log"
	"strings"
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

type lemaProper interface {
	AlwaysProper(string) bool
	RegularAndProper(string) bool
}

func main() {
	cmd.InitApp()
	lm, err := lema.NewCache()
	if err != nil {
		log.Fatal(err)
	}
	defer lm.Close()
	cmd.ProcessByLine(func(line string) (string, error) { return changeLine(strings.TrimSpace(line), lm), nil })
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
	return util.MultiSpacesRegexp.ReplaceAllString(res.String(), " ")
}

func changeWord(w string, lm lemaProper) string {
	wc := punct.PureWord(w)
	if wc == "" {
		return w
	}
	if util.SpecialWordRegexp.MatchString(wc) {
		return w
	}
	if isNumber(wc) {
		return w
	}
	if isQuoted(w) && lm.RegularAndProper(wc) {
		return changeToTitle(w)
	}
	if lm.AlwaysProper(wc) {
		return changeToTitle(w)
	}
	return strings.ToLower(w)
}

func changeToTitle(w string) string {
	r := []rune(strings.ToLower(w))
	i := 0
	if isQuoted(w) {
		i = 1
	}
	r[i] = unicode.ToUpper(r[i])
	return string(r)
}

func isNumber(w string) bool {
	rns := []rune(w)
	for _, r := range rns {
		if (unicode.IsNumber(r)) || punct.IsPunct(r) {
			continue
		}
		return false
	}
	return true
}

func isQuoted(w string) bool {
	rns := []rune(w)
	return len(rns) > 0 && rns[0] == punct.StartQuote
}

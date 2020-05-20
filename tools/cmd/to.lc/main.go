package main

import (
	"flag"
	"log"
	"strings"
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/abbr"
	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

type lemaProper interface {
	Regular(string) bool
	Proper(string) bool
	AbbreviationString(string) string
}

type params struct {
	abbrFile string
}

func main() {
	prms := params{}
	flag.StringVar(&prms.abbrFile, "abbr", "", "Abbreviations file")
	cmd.InitApp()
	if prms.abbrFile == "" {
		flag.Usage()
		log.Fatal("")
	}
	ad, err := abbr.NewAbbr(prms.abbrFile)
	if err != nil {
		log.Fatal(err)
	}

	lm, err := lema.NewCache()
	if err != nil {
		log.Fatal(err)
	}
	defer lm.Close()
	cmd.ProcessByLine(func(line string) (string, error) { return changeLine(strings.TrimSpace(line), lm, ad), nil })
}

func changeLine(line string, lm lemaProper, ad *abbr.Abbreviations) string {
	if len(line) == 0 {
		return line
	}
	strs := strings.Fields(line)
	res := strings.Builder{}
	res.Grow(len(line))
	for _, w := range strs {
		if w != "" {
			res.WriteString(" ")
			res.WriteString(changeWord(w, lm, ad))
		}
	}
	return util.FixSpaces(res.String())
}

func changeWord(w string, lm lemaProper, ad *abbr.Abbreviations) string {
	wc, _, _ := punct.TrimWord(w, punct.IsPunct)
	if wc == "" {
		return w
	}
	if util.SpecialWordRegexp.MatchString(wc) {
		return w
	}
	wcq, bg, end := punct.TrimWord(w, punct.IsAllSep)
	wc = punct.TrimQuote(wcq)
	if lema.IsNumber(wc) {
		return w
	}
	dt := ""
	if strings.HasPrefix(end, ".") {
		dt = "."
	}
	ai, aw := ad.Find(wc + dt)
	if ai > 0 {
		return aw + w[ai:]
	}

	aw = lm.AbbreviationString(wc + dt)
	if aw != "" {
		return aw + w[len(aw):]
	}

	if !lm.Regular(wc) && lm.Proper(wc) {
		return bg + changeToTitle(wcq) + end
	}
	if lm.Regular(wc) && !lm.Proper(wc) && !isQuoted(wcq) {
		return strings.ToLower(w)
	}
	if lm.Regular(wc) && lm.Proper(wc) && allUpper(wc) {
		return bg + changeToTitle(wcq) + end
	}
	return w
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

func allUpper(w string) bool {
	for _, r := range w {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func isQuoted(w string) bool {
	rns := []rune(w)
	return len(rns) > 0 && rns[0] == punct.StartQuote
}

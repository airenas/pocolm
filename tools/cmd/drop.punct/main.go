package main

import (
	"flag"
	"log"
	"strings"

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
			res.WriteString(clearWord(w, lm, ad))
		}
	}
	return util.FixSpaces(res.String())
}

func clearWord(w string, lm lemaProper, ad *abbr.Abbreviations) string {
	wc, _, _ := punct.TrimWord(w, punct.IsPunct)
	if wc == "" {
		return ""
	}
	if util.SpecialWordRegexp.MatchString(wc) {
		return wc
	}
	wcg, _, end := punct.TrimWord(w, punct.IsAllSep)
	wc = punct.TrimQuote(wcg)
	if lema.IsNumber(wc) {
		return trimNumber(w)
	}
	dt := ""
	if strings.HasPrefix(end, ".") {
		dt = "."
	}
	ai, aw := ad.Find(wc + dt)
	if ai > 0 {
		return aw
	}

	aw = lm.AbbreviationString(wc + dt)
	if aw != "" {
		return aw
	}

	return wcg
}

func trimNumber(w string) string {
	rw, b, e := punct.TrimWord(w, punct.IsAllSep)
	if strings.HasSuffix(b, "-") || strings.HasSuffix(b, "+") {
		rw = string(b[len(b)-1]) + rw
	}
	if strings.HasPrefix(e, "%") {
		rw = rw + string(e[0])
	}
	return rw
}

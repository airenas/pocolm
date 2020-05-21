package main

import (
	"flag"
	"log"
	"strings"
	"unicode"

	"github.com/airenas/pocolm/tools/internal/pkg/abbr"
	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

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
	cmd.ProcessByLine(func(l string) (string, error) {
		return processLine(l, ad)
	})
}

func processLine(line string, ad *abbr.Abbreviations) (string, error) {
	line = strings.TrimSpace(line)
	if line != "" {
		line = changeLine(strings.TrimSpace(line), ad)
	}
	return line, nil
}

func changeLine(line string, ad *abbr.Abbreviations) string {
	return groupAbbr(changePunct(line), ad)
}

func groupAbbr(line string, ad *abbr.Abbreviations) string {
	res := strings.Builder{}
	res.Grow(len(line))
	rns := []rune(line)
	for i := 0; i < len(rns); i++ {
		if i == 0 || !(util.LetterDigit(rns[i-1]) ||
			rns[i-1] == punct.StartQuote || rns[i-1] == punct.EndQuote) {
			ia, _ := ad.Find(string(rns[i:]))
			if ia != 0 {
				res.WriteRune(' ')
				res.WriteString(util.DropSpaces(string(rns[i : i+ia])))
				res.WriteRune(' ')
				i = i + ia - 1
				continue
			}
		}
		res.WriteRune(rns[i])
	}
	return util.FixSpaces(res.String())
}

func changePunct(line string) string {
	rns := []rune(line)
	res := strings.Builder{}
	res.Grow(len(rns))
	for i, r := range rns {
		pr := getByIndex(rns, i-1)
		rLetter := unicode.IsLetter(r)
		if pr == '.' && rLetter {
			res.WriteRune(' ')
		} else if pr == ',' && rLetter {
			res.WriteRune(' ')
		} else if pr == ':' && rLetter {
			res.WriteRune(' ')
		} else if r == '-' && unicode.IsLetter(pr) && unicode.IsLetter(getByIndex(rns, i+1)) {
			res.WriteRune(' ')
		} else if pr == '-' && rLetter && unicode.IsLetter(getByIndex(rns, i-2)) {
			res.WriteRune(' ')
		} else if pr == ')' && !unicode.IsSpace(r) {
			res.WriteRune(' ')
		} else if !unicode.IsSpace(pr) && r == '(' {
			res.WriteRune(' ')
		}
		res.WriteRune(r)
		pr = r
	}
	return util.FixSpaces(res.String())
}

func getByIndex(rns []rune, i int) rune {
	if i >= 0 && i < len(rns) {
		return rns[i]
	}
	return '\n'
}

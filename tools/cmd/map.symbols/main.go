package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
	"mvdan.cc/xurls/v2"
)

type replace struct {
	str   string
	regxp *regexp.Regexp
}

var replaces []*replace
var regDots *regexp.Regexp

func main() {
	cmd.InitApp()
	initRegexp()
	cmd.ProcessByLine(processLine)
}

func initRegexp() {
	replaces = make([]*replace, 0)
	replaces = append(replaces, &replace{str: "<EMAIL>", regxp: util.EMailRegexp})
	replaces = append(replaces, &replace{str: "<EMAIL>", regxp: newRegexp("mailto:<EMAIL>")})
	replaces = append(replaces, &replace{str: "<URL>", regxp: xurls.Relaxed()})
	replaces = append(replaces, &replace{str: " <PILDOMA> ", regxp: newRegexp("[_]{2,}|[\\.]{4,}| [\\.]{3,}")})

	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("(?i)^[A-Z][\\)\\.]")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("(?i)^M{0,4}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})[\\.)]")})

	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^(([0-9]){1,2}\\.){1,}")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^([0-9]){1,2}\\)")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^\\* ")})

	replaces = append(replaces, &replace{str: " ",
		regxp: newRegexp("HYPERLINK|FORMTEXT|PAGEREF _Toc\\d+ \\\\h \\d+|MERGEFIELD [„]?[\\p{L}_\\d]+[“]?")})

	regDots = newRegexp("\\. \\.")
}

func newRegexp(line string) *regexp.Regexp {
	r, err := regexp.Compile(line)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func processLine(line string) (string, error) {
	return changeLine(strings.TrimSpace(line)), nil
}

func changeLine(line string) string {
	if len(line) == 0 {
		return line
	}
	//fix dots
	for regDots.MatchString(line) {
		line = regDots.ReplaceAllString(line, "..")
	}

	for _, rep := range replaces {
		line = rep.regxp.ReplaceAllString(line, rep.str)
	}
	line = util.MultiSpacesRegexp.ReplaceAllString(line, " ")
	return strings.TrimSpace(line)
}

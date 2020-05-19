package main

import (
	"flag"
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

type params struct {
	changeEmail      bool
	changeURL        bool
	changeNumeration bool
}

func main() {
	prms := params{}
	flag.BoolVar(&prms.changeEmail, "email", false, "change emails to <EMAIL>")
	flag.BoolVar(&prms.changeURL, "url", false, "change URLs to <URL>")
	flag.BoolVar(&prms.changeNumeration, "num", false, "change initial numeration to <NUMERACIJA>")
	cmd.InitApp()
	initRegexp(&prms)
	cmd.ProcessByLine(processLine)
}

func initRegexp(p *params) {
	romanNumbers := "M{0,4}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})"
	replaces = make([]*replace, 0)
	if p.changeEmail {
		replaces = append(replaces, &replace{str: "<EMAIL>", regxp: util.EMailRegexp})
		replaces = append(replaces, &replace{str: "<EMAIL>", regxp: newRegexp("mailto:<EMAIL>")})
	}
	if p.changeURL {
		replaces = append(replaces, &replace{str: "<URL>", regxp: xurls.Strict()})
	}
	if p.changeNumeration {
		replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("(?i)^[A-Z][\\)\\.]")})
		replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("(?i)^" + romanNumbers + "[\\.)]")})

		replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^(([0-9]){1,3}\\.){1,}")})
		replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^([0-9]){1,3}\\)")})
		replaces = append(replaces, &replace{str: " <NUMERACIJA> $3", regxp: newRegexp("^(([0-9]){1,3})( priedas.)")})
		replaces = append(replaces, &replace{str: " <NUMERACIJA> $5", regxp: newRegexp("^(" + romanNumbers + ")( skyrius.)")})
		replaces = append(replaces, &replace{str: " <PUNKTAS> ", regxp: newRegexp("^[\\*-] ")})
	}

	replaces = append(replaces, &replace{str: " <PILDOMA> ", regxp: newRegexp("[_]{2,}|[\\.]{4,}| [\\.]{3,}")})
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

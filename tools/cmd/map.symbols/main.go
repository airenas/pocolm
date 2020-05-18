package main

import (
	"log"
	"os"
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
	log.SetOutput(os.Stderr)
	initRegexp()
	cmd.ProcessByLine(processLine)
}

func initRegexp() {
	replaces = make([]*replace, 0)
	replaces = append(replaces, &replace{str: "<EMAIL>", regxp: util.EMailRegexp})
	replaces = append(replaces, &replace{str: "<URL>", regxp: xurls.Strict()})
	replaces = append(replaces, &replace{str: " <PILDOMA> ", regxp: newRegexp("[_]{2,}")})
	replaces = append(replaces, &replace{str: " <PILDOMA> ", regxp: newRegexp("[\\.]{4,}")})
	replaces = append(replaces, &replace{str: " <PILDOMA> ", regxp: newRegexp(" [\\.]{3,}")})

	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[A-Z]\\)")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[a-z]\\)")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[A-Z]\\.")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[a-z]\\.")})

	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^(([0-9]){1,2}\\.){1,}")})
	replaces = append(replaces, &replace{str: " <NUMERACIJA> ", regxp: newRegexp("^[a-z]\\)")})

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

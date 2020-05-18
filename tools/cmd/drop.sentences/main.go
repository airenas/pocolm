package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/lema"
	"github.com/airenas/pocolm/tools/internal/pkg/punct"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
)

type lemaLt interface {
	IsLt(string) bool
	IsAbbreviation(string) bool
}

func main() {
	cmd.InitApp()
	cmd.Config.PrintF = printLine
	lm, err := lema.NewCache()
	if err != nil {
		log.Fatal(err)
	}
	defer lm.Close()
	cmd.ProcessByLine(func(line string) (string, error) { return changeLine(strings.TrimSpace(line), lm), nil })
}

func printLine(dest io.Writer, line string) error {
	if len(line) > 0 {
		_, err := fmt.Fprintf(dest, "%s\n", line)
		return err
	}
	return nil
}

func changeLine(line string, lm lemaLt) string {
	if len(line) == 0 {
		return line
	}
	if leave(line, lm) {
		return line
	}
	return ""
}

func leave(l string, lm lemaLt) bool {
	lt, nlt := calc(l, lm)
	if float64(nlt) >= 0.2*float64(lt+nlt) {
		return false
	}
	return true
}

func calc(line string, lm lemaLt) (int, int) {
	strs := strings.Split(line, " ")
	lt := 0
	nlt := 0
	for _, w := range strs {
		if lema.HasNonLT(w) {
			return 0, 1
		}

		t := getType(w, lm)
		if t == "lt" {
			lt++
		} else if t == "nlt" {
			nlt++
		}
	}
	return lt, nlt
}

func getType(w string, lm lemaLt) string {
	wc := punct.PureWord(w)
	if util.SpecialWordRegexp.MatchString(wc) {
		return "spec"
	}
	if lema.IsNumber(wc) {
		return "num"
	}
	if lm.IsLt(wc) {
		return "lt"
	}
	if lm.IsAbbreviation(wc) {
		return "abrv"
	}
	return "nlt"
}

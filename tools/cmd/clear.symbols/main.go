package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

var replaceableSymbols map[rune]rune

func init() {
	replaceableSymbols = make(map[rune]rune)
	for _, r := range []rune(" \t\r") {
		replaceableSymbols[r] = ' '
	}
	replaceableSymbols['–'] = '-'
}

func main() {
	log.SetOutput(os.Stderr)
	fs := flag.CommandLine
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s: [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[1:])
	f, err := util.NewReadWrapper(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	destination, err := util.NewWriteWrapper(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	rd := bufio.NewReader(f)
	ln := 0
	for {
		ln++
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(errors.Wrapf(err, "Line %d", ln))
		}
		line = strings.TrimSpace(line)
		line = strings.TrimSpace(line)
		nLine := changeLine(line)
		_, err = fmt.Fprintf(destination, "%s\n", nLine)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Cant write file"))
		}
	}
	log.Printf("Done")
}

func changeLine(line string) string {
	if len(line) == 0 {
		return line
	}
	runes := []rune(line)
	res := make([]rune, 0)
	for _, r := range runes {
		res = append(res, changeSymbol(r))
	}
	r := string(res)
	r = strings.TrimSpace(r)
	r = strings.ReplaceAll(r, "  ", " ")
	r = strings.ReplaceAll(r, "...", ".")
	r = strings.ReplaceAll(r, "..", ".")
	return string(r)
}

func changeSymbol(r rune) rune {
	s, ok := replaceableSymbols[r]
	if ok {
		return s
	}
	return r
}

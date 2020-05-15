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

var parenthesesSymbols map[rune]rune

func init() {
	parenthesesSymbols = make(map[rune]rune)
	parenthesesSymbols['('] = ')'
	parenthesesSymbols['['] = ']'
	parenthesesSymbols['{'] = '}'
	parenthesesSymbols['<'] = '>'
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
	rns := []rune(line)
	res := strings.Builder{}
	for i := 0; i < len(rns); i++ {
		p, f := parenthesesSymbols[rns[i]]
		if f {
			pi, f := findParenthesis(rns[i:], p)
			if f {
				res.WriteString(" <SKL> ")
				i += pi
				continue
			}
		}
		res.WriteRune(rns[i])
	}
	rs := res.String()
	rs = strings.ReplaceAll(rs, "  ", " ")
	return strings.TrimSpace(rs)
}

func findParenthesis(rns []rune, p rune) (int, bool) {
	ps := rns[0]
	c := 0
	for i, r := range rns {
		if r == ps {
			c++
		}
		if r == p {
			c--
			if c == 0 {
				return i, true
			}
		}
	}
	return 0, false
}

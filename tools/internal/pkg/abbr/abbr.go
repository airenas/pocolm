package abbr

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type abbr struct {
	next  map[rune]*abbr
	ch    rune
	final bool
}

//Abbreviations keeps map of possible values
type Abbreviations struct {
	list *abbr
}

//NewAbbrReader loads from reader
func NewAbbrReader(reader io.Reader) (*Abbreviations, error) {
	res := &Abbreviations{}
	res.list = newAbbr()

	rd := bufio.NewReader(reader)
	end := false
	for !end {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				end = true
			} else {
				return nil, errors.Wrap(err, "Read file line error")
			}
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			res.add(line)
		}
	}
	return res, nil
}

//NewAbbr loads from file
func NewAbbr(f string) (*Abbreviations, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewAbbrReader(file)
}

func newAbbr() *abbr {
	res := &abbr{}
	res.next = make(map[rune]*abbr)
	return res
}

//Find return index of matchin abbreiation
func (a *Abbreviations) Find(line string) (int, string) {
	ca := a.list
	fi := 0
	rs := strings.Builder{}
	rsNext := strings.Builder{}
	rns := []rune(line)
	for i, r := range rns {
		if !unicode.IsSpace(r) {
			na, f := ca.next[unicode.ToLower(r)]
			if f {
				if na.final && wordEnd(rns, i) {
					fi = i + 1
					rs.WriteString(rsNext.String())
					rs.WriteRune(na.ch)
					rsNext.Reset()
				} else {
					rsNext.WriteRune(na.ch)
				}
				ca = na
			} else {
				break
			}
		}
	}
	return fi, rs.String()
}

func wordEnd(rns []rune, i int) bool {
	return i+1 >= len(rns) || isNotLetterDigit(rns[i]) || isNotLetterDigit(rns[i+1])
}

func isNotLetterDigit(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func (a *Abbreviations) add(line string) {
	ca := a.list
	var na *abbr
	var f bool
	for _, r := range line {
		if !unicode.IsSpace(r) {
			wt := unicode.ToLower(r)
			na, f = ca.next[wt]
			if !f {
				na = newAbbr()
				na.ch = r
				ca.next[wt] = na
			}
			ca = na
		}
	}
	if na != nil {
		na.final = true
	}
}

func split(l string) []string {
	res := []string{}
	pr := 0
	for i, r := range l {
		if r == '.' || r == '/' || r == ' ' {
			res = append(res, l[pr:i+1])
			pr = i + 1
		}
	}
	if pr < len(l) {
		res = append(res, l[pr:])
	}
	return res
}

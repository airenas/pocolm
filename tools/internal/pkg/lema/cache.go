package lema

import (
	"bufio"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
)

type lInfo struct {
	abbreviation string //non empty if it is abbreviationa and has abbreviation form
	proper       bool
	regular      bool
	number       bool
}

//Cache helper class
type Cache struct {
	words     map[string]*lInfo
	path      string
	m         sync.Mutex
	save      bool
	vFileName string
}

//AbbreviationString returns abbreviation form it it has or empty
func (l *Cache) AbbreviationString(w string) string {
	if w == "" {
		return ""
	}
	r := l.getData(w)
	return r.abbreviation
}

//Proper returns true if word can be proper
func (l *Cache) Proper(w string) bool {
	if w == "" {
		return false
	}
	r := l.getData(w)
	return r.proper
}

//Regular returns true if word is regular
func (l *Cache) Regular(w string) bool {
	if w == "" {
		return false
	}
	r := l.getData(w)
	return r.regular
}

//NewCache creates lema cache
func NewCache() (*Cache, error) {
	l := Cache{}
	l.words = make(map[string]*lInfo)
	l.vFileName = l.vocabFile()
	err := l.loadMap()
	if err != nil {
		return nil, err
	}
	go l.runSave()
	return &l, nil
}

func (l *Cache) getData(w string) *lInfo {
	l.m.Lock()
	defer l.m.Unlock()

	r, ok := l.words[w]
	if ok {
		return r
	}
	r = l.getDataFromServer(w)
	l.words[w] = r
	l.save = true
	return r
}

func (l *Cache) getDataFromServer(w string) *lInfo {
	var res lInfo
	if unicode.IsLetter([]rune(w)[0]) && !HasNonLT(w) && !hasHTTPSymbols(w) {
		r, err := Analyze(w)

		if err != nil {
			panic(errors.Wrap(err, "Can't analyze '"+w+"'"))
		}
		res.proper = isProper(r)
		res.regular = isRegular(r)
		res.abbreviation = abbreviation(r)
	}

	return &res
}

func (l *Cache) vocabFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := path.Join(home, ".lema", "cache", "vocab")
	return dir
}

func (l *Cache) loadMap() error {
	l.m.Lock()
	defer l.m.Unlock()

	_, err := os.Stat(l.vFileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return l.loadVocab(l.vFileName)
}

//Close finalizes cache - saves to disk
func (l *Cache) Close() {
	l.saveVocab()
}

func (l *Cache) runSave() {
	for {
		time.Sleep(30 * time.Second)
		l.saveVocab()
	}
}
func (l *Cache) saveVocab() {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.save {
		return
	}
	dir := filepath.Dir(l.vFileName)
	os.MkdirAll(dir, 0770)
	file, err := os.OpenFile(l.vFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for k, v := range l.words {
		file.Write([]byte(k + " " + toStr(v) + "\n"))
	}
}

func toStr(l *lInfo) string {
	res := "w"
	if l.proper {
		res = res + "P"
	}
	if l.regular {
		res = res + "R"
	}
	if l.number {
		res = res + "N"
	}
	if l.abbreviation != "" {
		res = res + "A-" + l.abbreviation + "-"
	}
	return res
}

func (l *Cache) loadVocab(f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()

	return l.loadReader(file)
}

func (l *Cache) loadReader(reader io.Reader) error {
	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "Read file line error")
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			strs := strings.Split(line, " ")
			if len(strs) < 2 {
				return errors.Errorf("Wrong line '%s'", line)
			}
			li := lInfo{}
			s := strs[1]
			ia := strings.Index(s, "A-")
			if ia > 0 {
				li.abbreviation = strings.TrimRight(s[ia+2:], "-")
				s = s[:ia]
			}
			li.proper = strings.Index(s, "P") > -1
			li.regular = strings.Index(s, "R") > -1
			li.number = strings.Index(s, "N") > -1
			l.words[strs[0]] = &li
		}
	}
	return nil
}

func isProper(r *Result) bool {
	if r.Suffix != "" { // ignore our suffix check
		return false
	}
	for _, mi := range r.Mi {
		if strings.HasPrefix(mi.Mi, "I") {
			return true
		}
	}
	return false
}

func isRegular(r *Result) bool {
	for _, mi := range r.Mi {
		if !strings.HasPrefix(mi.Mi, "I") && !strings.HasPrefix(mi.Mi, "Y") {
			return true
		}
	}
	return false
}

func abbreviation(r *Result) string {
	res := ""
	for _, mi := range r.Mi {
		if !strings.HasPrefix(mi.MiVdu, "Y") {
			return ""
		}
		res = mi.MF
	}
	return res
}

var encoder = charmap.ISO8859_13.NewEncoder()

//HasNonLT non lt letters
func HasNonLT(w string) bool {
	_, err := encoder.String(w)
	return err != nil
}

func hasHTTPSymbols(w string) bool {
	return strings.Index(w, "/") > 0 || strings.Index(w, "%") > 0 || strings.Index(w, ">") > 0
}

//NewTestCache creates lema cache for testing
func NewTestCache(rd io.Reader) (*Cache, error) {
	c := Cache{}
	c.words = make(map[string]*lInfo)
	return &c, c.loadReader(rd)
}

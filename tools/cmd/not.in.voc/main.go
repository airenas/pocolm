package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/pkg/errors"
)

type params struct {
	vocabFile string
}

func main() {
	prms := params{}
	flag.StringVar(&prms.vocabFile, "v", "", "Vocabulary to compare with")
	cmd.InitApp()
	cmd.Config.PrintF = printLine
	if prms.vocabFile == "" {
		flag.Usage()
		log.Fatal("")
	}

	voc, err := loadVocab(prms.vocabFile)
	if err != nil {
		log.Fatal(err)
	}
	cmd.ProcessByLine(func(line string) (string, error) { return changeLine(strings.TrimSpace(line), voc) })
}

func printLine(dest io.Writer, line string) error {
	if len(line) > 0 {
		_, err := fmt.Fprintf(dest, "%s\n", line)
		return err
	}
	return nil
}

func changeLine(line string, voc map[string]bool) (string, error) {
	if len(line) == 0 {
		return line, nil
	}
	strs := strings.Fields(line)
	if len(strs) < 2 {
		return "", errors.Errorf("Wrong line '%s'", line)
	}
	if (len(strs) > 1) && voc[strs[1]] {
		return line, nil
	}
	return "", nil
}

func loadVocab(f string) (map[string]bool, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return loadVocabReader(file)
}

func loadVocabReader(reader io.Reader) (map[string]bool, error) {
	res := make(map[string]bool)
	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "Read file line error")
		}
		line = strings.TrimSpace(line)
		if line != "" {
			strs := strings.Fields(line)
			if len(strs) < 2 {
				return nil, errors.Errorf("Wrong line '%s'", line)
			}
			res[strs[0]] = true
		}
	}
	return res, nil
}

package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/airenas/pocolm/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

//Config app configuration
var Config cmdData

type cmdData struct {
	PrintF   func(io.Writer, string) error
	ProcessF func(string) (string, error)
}

func init() {
	Config.PrintF = func(dest io.Writer, line string) error { _, err := fmt.Fprintf(dest, "%s\n", line); return err }
}

//InitApp main function for tool to init
func InitApp() {
	log.SetOutput(os.Stderr)
	fs := flag.CommandLine
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s: [input-file | stdin] [output-file | stdout]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[1:])
}

//ProcessByLine main function for tool to read file and write file
func ProcessByLine(procF func(string) (string, error)) {
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
	end := false
	for !end {
		ln++
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				end = true
			} else {
				log.Fatal(errors.Wrapf(err, "Line %d", ln))
			}
		}
		nLine, err := procF(line)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Line %d", ln))
		}
		err = Config.PrintF(destination, nLine)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Cant write file"))
		}
	}
	log.Printf("Done")
}

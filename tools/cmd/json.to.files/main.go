package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/airenas/pocolm/tools/internal/pkg/cmd"
	"github.com/airenas/pocolm/tools/internal/pkg/util"
	"github.com/pkg/errors"
)

type article struct {
	Date     string `json:"date,omitempty"`
	Template string `json:"template"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	URL      string `json:"url"`
	Hash     string `json:"hash"`
}

type params struct {
	outDir    string
	overWrite bool
}

func main() {
	prms := params{}
	flag.StringVar(&prms.outDir, "d", "", "Output director")
	flag.BoolVar(&prms.overWrite, "f", false, "Force owerwrite files if exists")
	cmd.InitApp()
	if prms.outDir == "" {
		flag.Usage()
		log.Fatal("No outpur dir")
	}

	f, err := util.NewReadWrapper(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	info, err := os.Stat(prms.outDir)
	if os.IsNotExist(err) || info == nil {
		log.Printf("Creating dir %s", prms.outDir)
		err = os.MkdirAll(prms.outDir, 0755)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Can't create %s", prms.outDir))
		}
	}

	jDecoder := json.NewDecoder(f)
	jDecoder.DisallowUnknownFields()

	_, err = jDecoder.Token()
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Error reading json"))
	}
	count := 0
	for ; jDecoder.More(); count++ {
		var a article
		err := jDecoder.Decode(&a)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Error reading json"))
		}
		as := getData(&a)
		if as != "" {
			err = saveFile(makeFileName(prms.outDir, count+1), as, !prms.overWrite)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "Can't save file"))
			}
		}
	}
	log.Printf("Processed %d items", count)
	_, err = jDecoder.Token()
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Error reading json"))
	}
}

func getData(a *article) string {
	title := strings.TrimSpace(a.Title)
	r := []rune(title)
	if len(r) > 0 {
		lr := r[len(r)-1]
		if !(lr == '.' || lr == '!' || lr == '?') {
			title = title + "."
		}
	}
	return strings.TrimSpace(title + "\n" + strings.TrimSpace(a.Body))
}

func makeFileName(dir string, num int) string {
	fn := fmt.Sprintf("%06d.txt", num)
	return path.Join(dir, fn)
}

func saveFile(fileName string, data string, check bool) error {
	if check {
		info, err := os.Stat(fileName)
		if os.IsExist(err) || info != nil {
			return errors.Errorf("File %s exists.", fileName)
		}
	}
	return ioutil.WriteFile(fileName, []byte(data), 0664)
}

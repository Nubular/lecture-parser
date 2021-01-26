package parser

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	absPath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	inPath := filepath.Join(absPath, "input")
	xmlPath := filepath.Join(inPath, "simple_lec.xml")
	meta, err := GetMeta(xmlPath)
	if err != nil {
		log.Panic(err)
	}

	_, _ = GetSections(meta, xmlPath)
}

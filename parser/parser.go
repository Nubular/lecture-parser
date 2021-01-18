package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inPath := filepath.Join(absPath, "input")
	xmlPath := filepath.Join(inPath, "simple_lec.xml")
	meta, err := getMeta(xmlPath)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = getSections(meta, xmlPath)
}

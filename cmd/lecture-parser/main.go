package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser/parser"
)

var sections []parser.Section
var config Config

// @todo having an image tag followed by an audio tag could cause a crash

func main() {
	// if err := extractor.GetPDFPage("slides1.pdf", "image.jpeg", 1); err != nil {
	// 	log.Fatal(err)
	// }

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inPath := filepath.Join(absPath, "altInput")
	outPath := filepath.Join(absPath, "output")
	xmlPath := filepath.Join(inPath, "simple_lec.xml")

	start(xmlPath, inPath, outPath)
}

func printSections() {
	fmt.Print("\n[")
	for _, section := range sections {
		j, _ := json.MarshalIndent(section, "", "	")
		// fmt.Println(unsafe.Sizeof(section))
		fmt.Print(string(j), ",\n")
	}
	fmt.Print("]\n")
}

func start(xmlPath string, inPath string, outPath string) {
	var err error
	config, err = loadConfig()
	if err != nil {
		log.Println(err)
	}

	meta, err := parser.GetMeta(xmlPath)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println(meta)

	sections, err = parser.GetSections(meta, xmlPath)
	if err != nil {
		log.Panic(err)
	}

	printSections()
	err = getFrames(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}

	err = getAudio(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}

	err = getClips(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}
	serializeSections(outPath)
}

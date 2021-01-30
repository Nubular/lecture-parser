package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/nubular/lecture-parser/parser"
)

var sections []parser.Section
var config Config

func main() {
	// if err := extractor.GetPDFPage("slides1.pdf", "image.jpeg", 1); err != nil {
	// 	log.Fatal(err)
	// }

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inPath := filepath.Join(absPath, "input")
	outPath := filepath.Join(absPath, "output")
	xmlPath := filepath.Join(inPath, "simple_lec.xml")

	start(xmlPath, inPath, outPath)
}

func printSections() {
	fmt.Print("\n[")
	for _, section := range sections {
		j, _ := json.MarshalIndent(section, "", "	")
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
	fmt.Println(meta)
	yee := fmt.Sprintf("<does this print>")
	fmt.Println(yee, reflect.TypeOf(yee))
	sections, err = parser.GetSections(meta, xmlPath)
	if err != nil {
		log.Panic(err)
	}

	err = getFrames(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}

	err = getAudio(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}
	serializeSections(outPath)
	// printSections()
}

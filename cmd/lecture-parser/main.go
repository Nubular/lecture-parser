package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser/parser"
)

// @TODO: Remove these and use pointers instead
var sections []parser.Section
var config Config

// check for existence of files, extract file metadata (length of video, number of pages). Convert external videos to usable format.
func main() {

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inPath := filepath.Join(absPath, "input")
	outPath := filepath.Join(absPath, "output")
	xmlPath := filepath.Join(inPath, "example_lec.xml")

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

	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

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

	// printSections()
	err = getFrames(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}
	printSections()
	return
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

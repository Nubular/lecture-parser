package lecture

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nubular/lecture-parser/parser"
)

// @TODO: Remove these and use pointers instead
var sections []parser.Section
var config Config

func printSections() {
	fmt.Print("\n[")
	for _, section := range sections {
		j, _ := json.MarshalIndent(section, "", "	")
		// fmt.Println(unsafe.Sizeof(section))
		fmt.Print(string(j), ",\n")
	}
	fmt.Print("]\n")
}

func Start(xmlPath string, inPath string, outPath string) {

	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	var err error
	// config, err = loadConfig()
	// if err != nil {
	// 	log.Println(err)
	// }

	meta, err := parser.GetMeta(xmlPath)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println(meta)

	sections, err = parser.GetSections(meta, xmlPath)
	if err != nil {
		log.Panic(err)
	}

	err = getFrames(inPath, outPath)
	if err != nil {
		log.Panic(err)
	}
	// printSections()
	// return
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

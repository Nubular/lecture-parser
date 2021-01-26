package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser/extractor"
	"github.com/nubular/lecture-parser/parser"
)

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

func checkDupPage(frames []extractor.Frame, page int) bool {
	for _, frame := range frames {
		if frame.Page == page {
			return true
		}
	}
	return false
}

// Make sections global?
func getFrames(inPath string, outPath string, sections []parser.Section) error {

	if len(sections) == 0 {
		return errors.New("No sections received")
	}

	frames := make(map[string][]extractor.Frame)
	// prevFramePath := ""
	for _, section := range sections {
		// fmt.Println(section)
		if section.FrameType == "slide" {
			if _, exists := frames[section.SlideDeck.Src]; !exists {
				frames[section.SlideDeck.Src] = make([]extractor.Frame, 0)
			}
			if checkDupPage(frames[section.SlideDeck.Src], section.Page) {
				continue
			}

			frame := extractor.Frame{ImageName: fmt.Sprintf("%s_%d", section.SlideDeck.ID, section.Page), Page: section.Page}
			frames[section.SlideDeck.Src] = append(frames[section.SlideDeck.Src], frame)
		}
	}
	for src, frame := range frames {
		filepath := filepath.Join(inPath, src)
		err := extractor.GetPDFPages(filepath, outPath, frame)
		if err != nil {
			return err
		}
	}

	return nil
}

func start(xmlPath string, inPath string, outPath string) {

	config, err := loadConfig()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(config)
	meta, err := parser.GetMeta(xmlPath)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(meta)

	sections, err := parser.GetSections(meta, xmlPath)

	fmt.Println(getFrames(inPath, outPath, sections))

}

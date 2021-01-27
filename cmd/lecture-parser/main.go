package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser/parser"
	"github.com/nubular/lecture-parser/util"
)

var sections []parser.Section

func main() {
	// if err := extractor.GetPDFPage("slides1.pdf", "image.jpeg", 1); err != nil {
	// 	log.Fatal(err)
	// }

	absPath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	inPath := filepath.Join(absPath, "input")
	outPath := filepath.Join(absPath, "output")
	xmlPath := filepath.Join(inPath, "simple_lec.xml")

	err = start(xmlPath, inPath, outPath)
	if err != nil {
		log.Panic(err)
	}
}

func checkDupPage(frames []util.Frame, page int) bool {
	for _, frame := range frames {
		if frame.Page == page {
			return true
		}
	}
	return false
}

func printSections() {
	fmt.Print("\n[")
	for _, section := range sections {
		j, _ := json.MarshalIndent(section, "", "	")
		fmt.Print(string(j), ",\n")
	}
	fmt.Print("]\n")
}

func getFrames(inPath string, outPath string) error {

	if len(sections) == 0 {
		return errors.New("No sections received")
	}

	// map from filename to array of frames to be extracted.
	slides := make(map[string][]util.Frame)
	images := make([]util.Frame, 0)
	audio := make([]util.Frame, 0)
	video := make([]util.Frame, 0)

	for i, section := range sections {

		if section.FrameType == "slide" {
			if _, exists := slides[section.SlideDeck.Src]; !exists {
				slides[section.SlideDeck.Src] = make([]util.Frame, 0)
			}
			if checkDupPage(slides[section.SlideDeck.Src], section.Page) {
				continue
			}
			filename := fmt.Sprintf("%s_%d.jpg", section.SlideDeck.ID, section.Page)
			frame := util.Frame{FileName: filename, Page: section.Page}

			sections[i].FrameSrc.ImageSrc = filename

			slides[section.SlideDeck.Src] = append(slides[section.SlideDeck.Src], frame)
		} else if section.FrameType == "image" {
			filename := section.ResourceSrc
			sections[i].FrameSrc.ImageSrc = filename
			images = append(images, util.Frame{FileName: filename})

		} else if section.FrameType == "audio" {

			// I could take the previous frame displayed, but this way if you have a slide or image (with no corresponding ssml,
			//doesn't generate a section) followed by an audio tag it'll display stuff correctly
			filename := section.ResourceAttr["frameSrc"]
			sections[i].FrameSrc.ImageSrc = filename

			audio = append(audio, util.Frame{FileName: filename})

		} else if section.FrameType == "video" {

			filename := section.ResourceSrc
			sections[i].FrameSrc.VideoSrc = filename
			video = append(video, util.Frame{FileName: filename})
		} else {
			log.Println("Unidentified Section Type: ", section.FrameType)
		}

	}
	// Make this more async Awaitey?
	for src, frame := range slides {
		srcPath := filepath.Join(inPath, src)
		clipPath := filepath.Join(outPath, "FRAMES")
		err := util.GetPDFPages(srcPath, clipPath, frame)
		if err != nil {
			return err
		}
	}

	err := util.AsyncCopyFrames(inPath, filepath.Join(outPath, "FRAMES"), images)
	if err != nil {
		return err
	}
	err = util.AsyncCopyFrames(inPath, filepath.Join(outPath, "AUDIO"), audio)
	if err != nil {
		return err
	}
	err = util.AsyncCopyFrames(inPath, filepath.Join(outPath, "CLIPS"), video)
	if err != nil {
		return err
	}

	return nil
}

func start(xmlPath string, inPath string, outPath string) error {

	config, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Println(config)
	meta, err := parser.GetMeta(xmlPath)
	if err != nil {
		return err
	}
	fmt.Println(meta)

	sections, err = parser.GetSections(meta, xmlPath)
	if err != nil {
		return err
	}

	err = getFrames(inPath, outPath)
	if err != nil {
		return err
	}
	printSections()
	return nil
}

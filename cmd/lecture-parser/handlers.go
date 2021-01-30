package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/nubular/lecture-parser/parser"
	"github.com/nubular/lecture-parser/util"
)

func checkDupPage(frames []util.Frame, page int) bool {
	for _, frame := range frames {
		if frame.Page == page {
			return true
		}
	}
	return false
}

func getFrames(inPath string, outPath string) error {

	if len(sections) == 0 {
		return errors.New("No sections received")
	}

	// map from filename to array of frames to be extracted.
	slides := make(map[string][]util.Frame)
	images := make([]util.Frame, 0)
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

	err = util.AsyncCopyFrames(inPath, filepath.Join(outPath, "CLIPS"), video)
	if err != nil {
		return err
	}

	return nil
}

func getAudio(inPath string, outPath string) error {

	if len(sections) == 0 {
		return errors.New("No sections received")
	}

	ssml := make([]util.Frame, 0)
	audio := make([]util.Frame, 0)

	for i, section := range sections {

		if section.FrameType == "slide" {
			filename := fmt.Sprintf("%s_%d.mp3", section.SlideDeck.ID, section.ID)
			frame := util.Frame{FileName: filename, SSML: section.SSML}

			sections[i].FrameSrc.AudioSrc = filename

			ssml = append(ssml, frame)
		} else if section.FrameType == "audio" {

			// I could take the previous frame displayed, but this way if you have a slide or image (with no corresponding ssml,
			//doesn't generate a section) followed by an audio tag it'll display stuff correctly
			filename := section.ResourceAttr["frameSrc"]
			sections[i].FrameSrc.ImageSrc = filename

			audio = append(audio, util.Frame{FileName: filename})

		}
	}
	err := util.CreateMP3(filepath.Join(outPath, "AUDIO"), ssml, config.CacheFiles)
	if err != nil {
		return err
	}
	err = util.AsyncCopyFrames(inPath, filepath.Join(outPath, "AUDIO"), audio)
	if err != nil {
		return err
	}
	return nil
}

func serializeSections(outPath string) error {
	type sectionFile struct {
		yee []parser.Section `json:"yee"`
	}

	toWrite, _ := json.Marshal(sections)
	jsonPath := filepath.Join(outPath, "sections.json")
	return ioutil.WriteFile(jsonPath, toWrite, 0644)

}

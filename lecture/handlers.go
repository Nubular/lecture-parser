package lecture

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/nubular/lecture-parser/highlight"
	"github.com/nubular/lecture-parser/parser"
	"github.com/nubular/lecture-parser/util"
)

//:p
func checkDupPage(frames []parser.Section, page int) bool {
	for _, frame := range frames {
		if frame.Page == page {
			return true
		}
	}
	return false
}

// @todo delete sections of images which don't have ssml or audio
func getFrames(inPath string, outPath string) error {

	if len(sections) == 0 {
		return errors.New("no frame sections received")
	}

	// map from filename to array of frames to be extracted.

	slides := make(map[string][]parser.Section)
	highlightFrames := make([]parser.Section, 0)
	images := make([]parser.Section, 0)
	video := make([]parser.Section, 0)
	prevFramePath := ""
	// Folders the correspondng media can go in. Use empty string if no subfolder to be made
	imageFolder := "FRAMES"
	videoFolder := "CLIPS"
	for i, section := range sections {

		if section.FrameType == "slide" {
			// check if map has been made, if not then make it
			if _, exists := slides[section.SlideDeck.Src]; !exists {
				slides[section.SlideDeck.Src] = make([]parser.Section, 0)
			}

			// Check if page already in list to be rendered
			if !checkDupPage(slides[section.SlideDeck.Src], section.Page) {
				filename := fmt.Sprintf("%s_%d.jpg", section.SlideDeck.ID, section.Page)

				sections[i].FrameSrc.ImageSrc = filename
				slides[section.SlideDeck.Src] = append(slides[section.SlideDeck.Src], sections[i])
			}
			// requires original slide image to be in FRAMES
			if section.FrameSubType == "highlight" {

				resourceName := fmt.Sprintf("%s_%d.jpg", section.SlideDeck.ID, section.Page)
				// Need a better filename for highlighted slides
				filename := fmt.Sprintf("%s_%d_[%s].jpg", section.SlideDeck.ID, section.Page, section.ResourceAttr["points"])

				sections[i].ResourceAttr["srcImage"] = resourceName
				sections[i].FrameSrc.ImageSrc = filename

				highlightFrames = append(highlightFrames, sections[i])
			}
			prevFramePath = sections[i].FrameSrc.ImageSrc

		} else if section.FrameType == "image" {
			sections[i].FrameSrc.ImageSrc = section.ResourceSrc
			images = append(images, sections[i])

			prevFramePath = sections[i].FrameSrc.ImageSrc

		} else if section.FrameType == "video" {

			sections[i].FrameSrc.VideoSrc = section.ResourceSrc
			video = append(video, sections[i])
		} else if section.FrameType == "audio" {
			// I could take the previous frame displayed, but this way if you have a slide or image (with no corresponding ssml,
			//doesn't generate a section) followed by an audio tag it'll display stuff correctly
			if section.ResourceAttr["frameSrc"] != "" {
				sections[i].FrameSrc.ImageSrc = section.ResourceAttr["frameSrc"]
			} else {
				sections[i].FrameSrc.ImageSrc = prevFramePath
			}
		} else {
			log.Println("Unidentified Section Type: ", section.FrameType)
		}

	}

	// Make this more async Awaitey?
	for src, frame := range slides {
		srcPath := filepath.Join(inPath, src)
		outPath := filepath.Join(outPath, "FRAMES")
		err := util.GetPDFPages(srcPath, outPath, frame)
		if err != nil {
			return err
		}
	}
	if len(images) != 0 {
		err := util.AsyncCopyFrames(inPath, filepath.Join(outPath, imageFolder), images)
		if err != nil {
			return err
		}
	} else {
		log.Println("Could not identify any images to be transferred")
	}

	if config.ScriptPath == "" {
		return fmt.Errorf("no .py script path provided")
	}
	// until I figure out the syscall path issue
	if len(highlightFrames) != 0 {
		imageFolderPath := filepath.Join(outPath, "FRAMES")
		err := highlight.AsyncHighlightImage(imageFolderPath, filepath.Join(outPath, imageFolder), config.ScriptPath, highlightFrames)
		if err != nil {
			return err
		}
	} else {
		log.Println("Could not identify any images to be transferred")
	}
	if len(video) != 0 {
		err := util.AsyncCopyFrames(inPath, filepath.Join(outPath, videoFolder), video)
		if err != nil {
			return err
		}
	} else {
		log.Println("Could not identify any videos to be transferred")
	}
	return nil
}

func getAudio(inPath string, outPath string) error {

	if len(sections) == 0 {
		return errors.New("no audio sections received")
	}

	ssml := make([]parser.Section, 0)
	audio := make([]parser.Section, 0)
	// Folder the corresponding media goes in.
	audioFolder := "AUDIO"
	for i, section := range sections {

		if section.FrameType == "slide" {
			filename := fmt.Sprintf("%03d.mp3", section.ID)

			sections[i].FrameSrc.AudioSrc = filename

			ssml = append(ssml, sections[i])

		} else if section.FrameType == "image" {
			filename := fmt.Sprintf("%03d.mp3", section.ID)

			sections[i].FrameSrc.AudioSrc = filename

			ssml = append(ssml, sections[i])

		} else if section.FrameType == "audio" {

			sections[i].FrameSrc.AudioSrc = section.ResourceSrc
			audio = append(audio, sections[i])
		}
	}
	if len(ssml) != 0 {
		err := util.CreateMP3(filepath.Join(outPath, audioFolder), ssml, config.CacheFiles)
		if err != nil {
			return err
		}
	} else {
		log.Println("Could not identify any ssml to be generated")
	}
	if len(audio) != 0 {
		err := util.AsyncCopyFrames(inPath, filepath.Join(outPath, audioFolder), audio)
		if err != nil {
			return err
		}
	} else {
		log.Println("Could not identify any audio files to be transferred")
	}
	return nil
}

func serializeSections(outPath string) error {
	toWrite, _ := json.Marshal(sections)
	jsonPath := filepath.Join(outPath, "sections.json")
	return ioutil.WriteFile(jsonPath, toWrite, 0644)

}

func getClips(inPath string, outPath string) error {

	// Folder the corresponding media goes in
	imageFolder := "FRAMES"
	audioFolder := "AUDIO"
	videoFolder := "CLIPS"

	clips := make([]parser.Section, 0)

	for i, section := range sections {
		switch section.FrameType {
		case "slide", "audio", "image":
			filename := fmt.Sprintf("%03d.mp4", section.ID)

			sections[i].FrameSrc.VideoSrc = filename
			clips = append(clips, sections[i])
		case "video":

		default:
			log.Printf("Ignoring Unidentified section of type %s\n", section.FrameType)
		}

	}

	if len(clips) == 0 {
		log.Println("Could not identify any clips to be created")
		return nil
	}
	err := util.AsyncCombineImageAudio(inPath, outPath, imageFolder, audioFolder, videoFolder, clips)
	if err != nil {
		return err
	}
	return nil
}

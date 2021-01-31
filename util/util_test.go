package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/nubular/lecture-parser/parser"
)

// func TestCopyFrames(t *testing.T) {
// 	absPath, _ := os.Getwd()
// 	absPath = filepath.Dir(absPath)
// 	inPath := filepath.Join(absPath, "input")
// 	outPath := filepath.Join(absPath, "output")
// 	var frames []parser.Section
// 	// frames = append(frames, Frame{FileName: "Shutter.Island.2010.1080p.BluRay.x264.mp4"})
// 	frames = append(frames, parser.Section{FrameSrc.ImageSrc: "slides1.pdf"})
// 	frames = append(frames, Frame{FileName: "slides2.pdf"})

// 	err := AsyncCopyFrames(inPath, outPath, frames)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func Test
func TestCreateMP3(t *testing.T) {

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	absPath = filepath.Dir(absPath)

	outPath := filepath.Join(absPath, "output", "AUDIO")

	frames := make([]parser.Section, 0)
	frames = append(frames, parser.Section{SSML: `<speak>
You can also write custom ssml which works with specific platforms!
You say, <phoneme alphabet="ipa" ph="pɪˈkɑːn">pecan</phoneme>. 
I say, <phoneme alphabet="ipa" ph="ˈpi.kæn">pecan</phoneme>.
Sometimes it can be useful to <prosody volume="loud">increase the volume 
for a specific speech.</prosody>
 
</speak>`, FrameSrc: struct {
		ImageSrc string "json:\"imageSrc,omitempty\""
		AudioSrc string "json:\"audioSrc,omitempty\""
		VideoSrc string "json:\"videoSrc,omitempty\""
	}{AudioSrc: "test3.mp3"}})
	frames = append(frames, parser.Section{SSML: `<speak>
Now we are on the second page of slide 1. 
We can also change the slide deck that is open!
</speak>`, FrameSrc: struct {
		ImageSrc string "json:\"imageSrc,omitempty\""
		AudioSrc string "json:\"audioSrc,omitempty\""
		VideoSrc string "json:\"videoSrc,omitempty\""
	}{AudioSrc: "test2.mp3"}})

	err = CreateMP3(outPath, frames, true) // change the flag to send request to aws
	if err != nil {
		fmt.Println(err)
	}
}

func TestCobmineImageAudio(i *testing.T) {
	// path, err := exec.LookPath("ffmpeg")
	// if err != nil {
	// 	log.Fatal("installing fortune is in your future")
	// }
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	imagePath := filepath.Join(absPath, "output", "FRAMES", "slides1_1.jpg")
	audioPath := filepath.Join(absPath, "output", "AUDIO", "001.mp3")
	absPath = filepath.Dir(absPath)
	err = combineImageAudio(imagePath, audioPath, "res.mp4")
	if err != nil {
		fmt.Println(err)
	}
}

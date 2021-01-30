package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFrames(t *testing.T) {
	absPath, _ := os.Getwd()
	absPath = filepath.Dir(absPath)
	inPath := filepath.Join(absPath, "input")
	outPath := filepath.Join(absPath, "output")
	var frames []Frame
	// frames = append(frames, Frame{FileName: "Shutter.Island.2010.1080p.BluRay.x264.mp4"})
	frames = append(frames, Frame{FileName: "slides1.pdf"})
	frames = append(frames, Frame{FileName: "slides2.pdf"})

	err := AsyncCopyFrames(inPath, outPath, frames)
	if err != nil {
		t.Error(err)
	}
}

// func Test
func TestCreateMP3(t *testing.T) {

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	absPath = filepath.Dir(absPath)

	outPath := filepath.Join(absPath, "output", "AUDIO")

	frames := make([]Frame, 0)
	frames = append(frames, Frame{SSML: `<speak>
You can also write custom ssml which works with specific platforms!
You say, <phoneme alphabet="ipa" ph="pɪˈkɑːn">pecan</phoneme>. 
I say, <phoneme alphabet="ipa" ph="ˈpi.kæn">pecan</phoneme>.
Sometimes it can be useful to <prosody volume="loud">increase the volume 
for a specific speech.</prosody>
 
</speak>`, FileName: "test3.mp3"})
	frames = append(frames, Frame{SSML: `<speak>
Now we are on the second page of slide 1. 
We can also change the slide deck that is open!
</speak>`, FileName: "test2.mp3"})
	frames = append(frames, Frame{SSML: `<speak>
Welcome to a example that will demonstrate the basic features this ssml extension

Currently, slide 1 is set as the active slide deck
 
</speak>`, FileName: "test1.mp3"})
	err = CreateMP3(outPath, frames, false)
	if err != nil {
		fmt.Println(err)
	}
}

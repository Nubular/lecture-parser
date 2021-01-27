package util

import (
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
	frames = append(frames, Frame{FileName: "Shutter.Island.2010.1080p.BluRay.x264.mp4"})
	frames = append(frames, Frame{FileName: "slides1.pdf"})
	frames = append(frames, Frame{FileName: "slides2.pdf"})

	err := AsyncCopyFrames(inPath, outPath, frames)
	if err != nil {
		t.Error(err)
	}
}

// func Test

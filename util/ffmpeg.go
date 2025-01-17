package util

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nubular/lecture-parser/parser"
)

func combineImageAudio(imagePath string, audioPath string, outPath string) error {
	cmd := exec.Command("ffmpeg",
		"-hide_banner",
		"-loglevel", "warning",
		"-y",
		"-loop", "1", "-i", imagePath,
		"-i", audioPath,
		"-tune", "stillimage",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-vf", "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2,setsar=sar=1:1",
		"-shortest",
		"-fflags",
		"+shortest",
		"-max_interleave_delta", "100M",
		outPath)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("[FFMPEG combine] %s for image: %s and audio %s", stdoutStderr, imagePath, audioPath)
	}
	return nil
}

func AsyncCombineImageAudio(inPath string, outPath string, imageFolder string, audioFolder string, videoFolder string, frames []parser.Section) error {
	if len(frames) == 0 {
		return errors.New("no video sections received")
	}
	checkPath := filepath.Join(outPath, videoFolder)
	if _, err := os.Stat(checkPath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", checkPath)
		os.Mkdir(checkPath, os.ModePerm)
	}

	concurrency := 4
	sem := make(chan bool, concurrency)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, frame := range frames {
		sem <- true
		imagePath := filepath.Join(outPath, imageFolder, frame.FrameSrc.ImageSrc)
		audioPath := filepath.Join(outPath, audioFolder, frame.FrameSrc.AudioSrc)
		videoOutPath := filepath.Join(outPath, videoFolder, frame.FrameSrc.VideoSrc)
		go func(imagePath string, audioPath string, videoOutPath string) {
			defer func() { <-sem }()
			select {
			case <-ctx.Done():
				return
			default:
			}
			err := combineImageAudio(imagePath, audioPath, videoOutPath)
			// log.Printf("Done combining to %s", videoOutPath)
			if err != nil {
				log.Println(err)
				cancel()
				return
			}
		}(imagePath, audioPath, videoOutPath)

	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return ctx.Err()
}

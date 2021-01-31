package util

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sync"

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
		"-vf", "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2",
		"-shortest",
		"-fflags",
		"+shortest",
		"-max_interleave_delta", "100M",
		outPath)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s for image: %s and audio %s", stdoutStderr, imagePath, audioPath)
	}
	return nil
}

func AsyncCombineImageAudio(inPath string, outPath string, imageFolder string, audioFolder string, videoFolder string, frames []parser.Section) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, frame := range frames {
		imagePath := filepath.Join(outPath, imageFolder, frame.FrameSrc.ImageSrc)
		audioPath := filepath.Join(outPath, audioFolder, frame.FrameSrc.AudioSrc)
		videoOutPath := filepath.Join(outPath, videoFolder, frame.FrameSrc.VideoSrc)
		wg.Add(1)
		go func(imagePath string, audioPath string, videoOutPath string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			err := combineImageAudio(imagePath, audioPath, videoOutPath)
			log.Printf("Done combining to %s", videoOutPath)
			if err != nil {
				log.Println(err)
				cancel()
				return
			}
		}(imagePath, audioPath, videoOutPath)

	}
	wg.Wait()
	return ctx.Err()
}

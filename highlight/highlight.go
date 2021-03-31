package highlight

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nubular/lecture-parser/parser"
)

func highlightImage(imagePath string, points string, outPath string, scriptPath string) error {
	list := strings.Split(points, " ")

	// absPath, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	time.Sleep(3 * time.Second)
	// absPath = filepath.Dir(absPath)
	// pythonPath := filepath.Join(absPath, "highlight/east.py")
	command := []string{scriptPath, "--input", string(imagePath), "-gr", "--output", string(outPath), "--list"}
	command = append(command, list...)
	cmd := exec.Command("python", command...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s for image: %s and outPath %s with list %v", stdoutStderr, imagePath, outPath, list)
	}
	return nil
}

func AsyncHighlightImage(inPath string, outPath string, scriptPath string, frames []parser.Section) error {

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", outPath)
		os.Mkdir(outPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("[Highlighter] %s", err)
		}
	}

	concurrency := 6
	sem := make(chan bool, concurrency)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, frame := range frames {
		sem <- true
		resourceSrc := frame.ResourceAttr["srcImage"]
		fileSrc := frame.FrameSrc.ImageSrc
		frameInPath := filepath.Join(inPath, resourceSrc)
		frameOutPath := filepath.Join(outPath, fileSrc)
		points := frame.ResourceAttr["points"]
		go func(frameInPath string, points string, frameOutPath string) {
			defer func() { <-sem }()
			select {
			case <-ctx.Done():
				return
			default:
			}
			err := highlightImage(frameInPath, points, frameOutPath, scriptPath)

			if err != nil {
				log.Println(err)
				cancel()
				return
			}
		}(frameInPath, points, frameOutPath)

	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return ctx.Err()
}

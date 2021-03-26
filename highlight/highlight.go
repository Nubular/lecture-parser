package highlight

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nubular/lecture-parser/parser"
)

func highlightImage(imagePath string, points string, outPath string) error {
	list := strings.Split(points, " ")
	command := []string{"highlight/east.py", "--input", string(imagePath), "--output", string(outPath), "--list"}
	command = append(command, list...)
	cmd := exec.Command("python", command...)
	// cmd := exec.Command("python", "test.py")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s for image: %s and outPath %s", stdoutStderr, imagePath, outPath)
	}
	return nil
}

// func main() {
// 	highlightImage("FRAMES/CN1_1.jpg", "11 22", "yeeee2.jpg")
// }

func AsyncHighlightImage(inPath, outPath string, frames []parser.Section) error {

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", outPath)
		os.Mkdir(outPath, os.ModePerm)
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, frame := range frames {
		resourceSrc := frame.ResourceAttr["srcImage"]
		fileSrc := frame.FrameSrc.ImageSrc
		print(fileSrc)
		frameInPath := filepath.Join(inPath, resourceSrc)
		frameOutPath := filepath.Join(outPath, fileSrc)
		log.Println(frameInPath, frameOutPath)
		points := frame.ResourceAttr["points"]
		wg.Add(1)

		go func(frameInPath string, points string, frameOutPath string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			err := highlightImage(frameInPath, points, frameOutPath)
			if err != nil {
				log.Println(err)
				cancel()
				return
			}
		}(frameInPath, points, frameOutPath)

	}
	wg.Wait()
	return ctx.Err()
}

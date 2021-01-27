package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func copyFile(frameInPath, frameOutPath string) (int64, error) {
	sourceFileStat, err := os.Stat(frameInPath)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", frameInPath)
	}
	source, err := os.Open(frameInPath)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.Create(frameOutPath)
	if err != nil {
		return 0, err
	}
	defer dest.Close()

	nBytes, err := io.Copy(dest, source)
	if err != nil {
		return 0, err
	}

	return nBytes, nil
}

// AsyncCopyFrames copies all files in the FileName Field in the Frame struct
// https://gist.github.com/lucassha/9ffd60225790bdf071e7969e91cbbdb5
func AsyncCopyFrames(inPath, outPath string, frames []Frame) error {

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", outPath)
		os.Mkdir(outPath, os.ModePerm)
	}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, frame := range frames {
		frameInPath := filepath.Join(inPath, frame.FileName)
		frameOutPath := filepath.Join(outPath, frame.FileName)
		wg.Add(1)

		go func(frameInPath, frameOutPath string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			_, err := copyFile(frameInPath, frameOutPath)
			if err != nil {
				log.Println(err)
				cancel()
				return
			}
			// log.Println("Done Writing to ", frameOutPath)
		}(frameInPath, frameOutPath)

	}
	wg.Wait()
	return ctx.Err()
}

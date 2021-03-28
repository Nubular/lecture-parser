package lecture_test

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser"
)

func ExampleStart() {
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inPath := filepath.Join(absPath, "input")
	outPath := filepath.Join(absPath, "output")
	xmlPath := filepath.Join(inPath, "example_lec.xml")

	lecture.Start(xmlPath, inPath, outPath)
	// Output: test
}

// func TestStart(t *testing.T) {
// 	absPath, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	inPath := filepath.Join(absPath, "input")
// 	outPath := filepath.Join(absPath, "output")
// 	xmlPath := filepath.Join(inPath, "example_lec.xml")

// 	lecture.Start(xmlPath, inPath, outPath)
// }

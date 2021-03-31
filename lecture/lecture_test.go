package lecture_test

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser/lecture"
)

func ExampleStart() {
	conf := lecture.Config{
		CacheFiles: false,
		ScriptPath: "/run/media/nubular/Shared/repo/lecture-parser/test/east.py",
	}
	lecture.SetConfig(conf)
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	absPath = filepath.Dir(absPath)
	// inPath := filepath.Join(absPath, "input")
	// outPath := filepath.Join(absPath, "output")
	// xmlPath := filepath.Join(inPath, "example_lec.xml")
	inPath := filepath.Join(absPath, "example_input")
	outPath := filepath.Join(absPath, "output")
	xmlPath := filepath.Join(inPath, "SampleLec.xml")

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

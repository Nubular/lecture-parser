package main

import (
	"log"

	"github.com/nubular/lecture-parser/extractor"
)

func main() {
	if err := extractor.ConvertPdfToJpg("test.pdf", "image.jpeg"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/nubular/lecture-parser/extractor"
)

func main() {
	if err := extractor.GetPDFPage("test.pdf", "image.jpeg", 2); err != nil {
		log.Fatal(err)
	}
}

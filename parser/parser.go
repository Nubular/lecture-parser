package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func controlTag(tag string) bool {
	ignore := []string{"info", "settings", "deck", "slide", "audio", "image", "video", "lecture"}
	flag := true
	for _, element := range ignore {
		if tag == element {
			flag = false
		}
	}
	return flag
}

func main() {
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	Inpath := filepath.Join(absPath, "input")
	xmlFile, err := os.Open(filepath.Join(Inpath, "simple_lec.xml"))
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	_ = byteValue
	r := strings.NewReader(string(byteValue))

	decoder := xml.NewDecoder(r)

	var start, end int64
	foundStart := false
	_ = foundStart

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "slide" {
				fmt.Println("yes my sweet nigga")
			}

			if controlTag(se.Name.Local) {
				break
			}
			fmt.Println("New Section: ")
			fmt.Println("<speak>\n", string(byteValue)[start:end])
			start = end

			start = decoder.InputOffset()
		}
		end = decoder.InputOffset()
	}
}

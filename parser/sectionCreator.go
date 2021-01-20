package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// SECTIONS Contains all generated sections
var SECTIONS []Section

// CURRENT can't be capital for some reason
var CURRENT Section

// META I like uppercase variables
var META *LectureMeta

// Section is a Struct that holds the information about the current section.
type Section struct {
	ID           int               `json:"id"`
	Voice        string            `json:"voice"`
	FrameType    string            `json:"frameType"`
	ResourceSrc  string            `json:"resourceSrc,omitempty"`
	ResourceAttr map[string]string `json:"resourceAttr,omitempty"`
	FrameFit     string            `json:"frameFit,omitempty"`
	SlideDeck    Deck              `json:"slideDeck,omitempty"`
	Page         int               `json:"page"`
	SSML         string            `json:"ssml,omitempty"`
}

func customTag(tag string) bool {
	ignore := []string{"info", "settings", "deck", "slide", "audio", "image", "video", "lecture"}
	flag := false
	for _, element := range ignore {
		if tag == element {
			flag = true
		}
	}
	return flag
}

func infoTag(tag string) bool {
	ignore := []string{"info", "settings", "deck"}
	flag := false
	for _, element := range ignore {
		if tag == element {
			flag = true
		}
	}
	return flag
}

// handling it like this can lead to ordering problems <slide page="x" deck="x"/> vs. <slide deck="x" page="x" />
func handleSlide(tag xml.StartElement) {
	if len(tag.Attr) == 0 {
		log.Println("Ignoring <slide/> with no attr")
	}
	CURRENT.FrameType = "slide"
	for _, attr := range tag.Attr {

		if attr.Name.Local == "deck" {
			// CURRENT.SlideDeck = attr.Value
			for _, deck := range META.Deck {
				if deck.ID == attr.Value {
					CURRENT.SlideDeck = deck
				}
			}
		} else if attr.Name.Local == "page" {

			num, err := strconv.Atoi(attr.Value)
			if err != nil {
				log.Println(err)
			}

			match, _ := regexp.MatchString("^[+|-][0-9]+$", attr.Value)
			if match {
				CURRENT.Page = CURRENT.Page + num
			} else {
				CURRENT.Page = num
			}
			if CURRENT.Page < 0 {
				fmt.Println("Page Index is less than 0")
			}
		}
	}
}

func handleImage(tag xml.StartElement) {
	if len(tag.Attr) == 0 {
		log.Println("Ignoring <image/> with no attr")
	}

	CURRENT.FrameType = "image"

	for _, attr := range tag.Attr {
		if attr.Name.Local == "src" {
			CURRENT.ResourceSrc = attr.Value
		}
		if attr.Name.Local == "fit" {
			CURRENT.FrameFit = attr.Value
		}
	}
}

func addResourceSection(tag xml.StartElement) {
	if len(tag.Attr) == 0 {
		log.Printf("Ignoring <%s/> with no attr\n", tag.Name.Local)
	}

	section := Section{ID: CURRENT.ID}
	section.ResourceAttr = make(map[string]string)

	for _, attr := range tag.Attr {
		if attr.Name.Local == "src" {
			section.ResourceSrc = attr.Value
		}
		section.ResourceAttr[attr.Name.Local] = attr.Value
	}
	if tag.Name.Local == "audio" {
		section.FrameType = "audio"
	} else if tag.Name.Local == "video" {
		section.FrameType = "video"
	}

	SECTIONS = append(SECTIONS, section)
	CURRENT.ID = CURRENT.ID + 1

}

func handleControlTag(tag xml.StartElement) {
	switch tag.Name.Local {

	case "slide":
		handleSlide(tag)
	case "audio", "video":
		addResourceSection(tag)
	case "image":
		handleImage(tag)

	}
}

func addSSMLSection(ssml string) {
	ssml = strings.TrimSpace(ssml)
	if ssml == "" {
		return
	}

	ssml = fmt.Sprint("<speak>\n", ssml, "\n<speak/>")
	section := Section{
		ID:        CURRENT.ID,
		Voice:     CURRENT.Voice,
		FrameType: CURRENT.FrameType,
		SSML:      ssml,
	}

	switch CURRENT.FrameType {

	case "slide":
		section.SlideDeck = CURRENT.SlideDeck
		section.Page = CURRENT.Page
	case "image":
		section.ResourceSrc = CURRENT.ResourceSrc
		section.FrameFit = CURRENT.FrameFit
	}
	// j, _ := json.Marshal(section)
	// fmt.Println(string(j))
	SECTIONS = append(SECTIONS, section)

	CURRENT.ID = CURRENT.ID + 1
}

func getSections(meta *LectureMeta, xmlPath string) (string, error) {

	// metajson, _ := json.Marshal(meta)
	// fmt.Println(string(metajson))

	META = meta

	CURRENT = Section{
		ID:        1,
		Voice:     meta.Settings.Voice,
		FrameType: "slide",
		SlideDeck: meta.ActiveDeck,
		Page:      1,
	}

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	r := strings.NewReader(string(byteValue))

	decoder := xml.NewDecoder(r)

	var start, end int64

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			// Ignore any tags that might belong to ssml
			if customTag(se.Name.Local) {

				if !infoTag(se.Name.Local) {
					addSSMLSection(string(byteValue)[start:end])
				}

				handleControlTag(se)

				start = end
				start = decoder.InputOffset()
			}

		}
		end = decoder.InputOffset()
	}

	printSections()
	return "", nil
}

func printSections() {
	fmt.Print("\n[")
	for _, section := range SECTIONS {
		j, _ := json.MarshalIndent(section, "", "	")
		fmt.Print(string(j), ",\n")
	}
	fmt.Print("]\n")
}

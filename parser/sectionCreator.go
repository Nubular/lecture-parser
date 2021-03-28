package parser

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

// handle various attributes in tags missing
// no page number bound checking

// sections Contains all generated sections
var sections []Section

// current can't be capital for some reason
var current Section

// meta I like uppercase variables
var meta *LectureMeta

// Section is a Struct that holds the information about the current section.
type Section struct {
	ID           int    `json:"id"`
	Voice        string `json:"voice"`
	FrameType    string `json:"frameType"`
	FrameSubType string `json:"frameSubType,omitempty"`
	FrameSrc     struct {
		ImageSrc string `json:"imageSrc,omitempty"`
		AudioSrc string `json:"audioSrc,omitempty"`
		VideoSrc string `json:"videoSrc,omitempty"`
	} `json:"frameSrc,omitempty"`
	ResourceSrc  string            `json:"resourceSrc,omitempty"`
	ResourceAttr map[string]string `json:"resourceAttr,omitempty"`
	FrameFit     string            `json:"frameFit,omitempty"`
	SlideDeck    Deck              `json:"slideDeck,omitempty"`
	Page         int               `json:"page"`
	SSML         string            `json:"ssml,omitempty"`
}

func customTag(tag string) bool {
	ignore := []string{"info", "settings", "deck", "slide", "highlight", "audio", "image", "video", "lecture"}
	flag := false
	for _, element := range ignore {
		if tag == element {
			flag = true
		}
	}
	return flag
}

func infoTag(tag string) bool {
	ignore := []string{"info", "settings", "deck", "lecture"}
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
	current.FrameType = "slide"
	if tag.Name.Local == "highlight" {
		current.FrameSubType = "highlight"
	}
	current.ResourceAttr = make(map[string]string)
	for _, attr := range tag.Attr {
		if attr.Name.Local == "deck" {
			if current.SlideDeck.ID == attr.Value {
				continue
			}
			for _, deck := range meta.Deck {
				if deck.ID == attr.Value {
					current.SlideDeck = deck
					current.Page = 1
				}
			}
		} else if attr.Name.Local == "page" {

			num, err := strconv.Atoi(attr.Value)
			if err != nil {
				log.Println(err)
			}

			match, _ := regexp.MatchString("^[+|-][0-9]+$", attr.Value)
			if match {
				current.Page = current.Page + num
			} else {
				current.Page = num
			}
			if current.Page < 0 {
				fmt.Println("Page Index is less than 0")
			}
		}
		current.ResourceAttr[attr.Name.Local] = strings.TrimSpace(attr.Value)
	}
}

func handleImage(tag xml.StartElement) {
	if len(tag.Attr) == 0 {
		log.Println("Ignoring <image/> with no attr")
	}

	current.FrameType = "image"
	current.ResourceAttr = make(map[string]string)

	for _, attr := range tag.Attr {
		current.ResourceAttr[attr.Name.Local] = strings.TrimSpace(attr.Value)
		if attr.Name.Local == "src" {
			current.ResourceSrc = attr.Value
		}
		if attr.Name.Local == "fit" {
			current.FrameFit = attr.Value
		}
	}
}

func addResourceSection(tag xml.StartElement) {
	if len(tag.Attr) == 0 {
		log.Printf("Ignoring <%s/> with no attr\n", tag.Name.Local)
	}

	section := Section{ID: current.ID}
	section.ResourceAttr = make(map[string]string)

	for _, attr := range tag.Attr {
		if attr.Name.Local == "src" {
			section.ResourceSrc = attr.Value
		}
		section.ResourceAttr[attr.Name.Local] = strings.TrimSpace(attr.Value)
	}

	if tag.Name.Local == "audio" {

		section.FrameType = "audio"
		if _, exists := section.ResourceAttr["frameSrc"]; !exists {
			section.ResourceAttr["frameSrc"] = current.ResourceSrc
		}
	} else if tag.Name.Local == "video" {
		section.FrameType = "video"
	}

	sections = append(sections, section)
	current.ID = current.ID + 1

}

func handleControlTag(tag xml.StartElement) {
	switch tag.Name.Local {

	case "slide", "highlight":
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

	ssml = fmt.Sprintf("<speak>\n %s \n</speak>", ssml)
	section := Section{
		ID:        current.ID,
		Voice:     current.Voice,
		FrameType: current.FrameType,
		SSML:      ssml,
	}

	switch current.FrameType {

	case "slide":
		section.SlideDeck = current.SlideDeck
		section.Page = current.Page
		// Attributes don't carry over between tags.
		if current.FrameSubType == "highlight" {
			section.FrameSubType = current.FrameSubType
			section.ResourceAttr = current.ResourceAttr
			current.FrameSubType = ""
		}
	case "image":
		section.ResourceSrc = current.ResourceSrc
		section.FrameFit = current.FrameFit
		section.ResourceAttr = current.ResourceAttr
	}
	// j, _ := json.Marshal(section)
	// fmt.Println(string(j))
	sections = append(sections, section)

	current.ID = current.ID + 1
}

func GetSections(metadata *LectureMeta, xmlPath string) ([]Section, error) {

	meta = metadata

	current = Section{
		ID:        1,
		Voice:     meta.Settings.Voice,
		FrameType: "slide",
		SlideDeck: meta.ActiveDeck,
		Page:      1,
	}

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		log.Println("Error opening XML: ", xmlPath)
		return nil, err
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
		s := string(byteValue)[start:end]
		s = strings.ReplaceAll(s, "\n", `\n`)
		switch se := t.(type) {
		case xml.StartElement:
			// Ignore any tags that might belong to ssml
			if customTag(se.Name.Local) {

				if !infoTag(se.Name.Local) {
					addSSMLSection(string(byteValue)[start:end])
				}
				handleControlTag(se)

				start = decoder.InputOffset()
			}
		case xml.EndElement:
			if customTag(se.Name.Local) {
				// handle closing lecture tag.
				if se.Name.Local == "lecture" {
					addSSMLSection(string(byteValue)[start:end])
				}

				start = decoder.InputOffset()
			}

		}
		end = decoder.InputOffset()
	}

	return sections, nil
}

func printSections() {
	fmt.Print("\n[")
	for _, section := range sections {
		j, _ := json.MarshalIndent(section, "", "	")
		fmt.Print(string(j), ",\n")
	}
	fmt.Print("]\n")
}

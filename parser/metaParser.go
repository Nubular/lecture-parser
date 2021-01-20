package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

// LectureMeta represents the relevant metadata in the xml
type LectureMeta struct {
	XMLName xml.Name `xml:"lecture" json:"-"`
	Info    struct {
		Title       string `xml:"title,attr" json:"-title"`
		Description string `xml:"description,attr" json:"-description"`
		Authors     string `xml:"authors,attr" json:"-authors"`
	} `xml:"info"`
	Settings struct {
		Voice               string `xml:"voice,attr" json:"-voice"`
		Resolution          string `xml:"resolution,attr" json:"-resolution"`
		Fps                 string `xml:"fps,attr" json:"-fps"`
		BreakAfterSlide     string `xml:"breakAfterSlide,attr" json:"-breakAfterSlide"`
		BreakAfterParagraph string `xml:"breakAfterParagraph,attr" json:"-breakAfterParagraph"`
	} `xml:"settings"`
	Deck []Deck `xml:"deck"`
	Mark []struct {
		Name    string `xml:"name,attr" json:"-name"`
		Chapter string `xml:"chapter,attr" json:"-chapter,omitempty"`
	} `xml:"mark"`
	Audio []struct {
		Src     string `xml:"src,attr" json:"-src"`
		ClipEnd string `xml:"clipEnd,attr" json:"-clipEnd,omitempty"`
	} `xml:"audio"`
	Image []struct {
		Src string `xml:"src,attr" json:"-src"`
		Fit string `xml:"fit,attr" json:"-fit,omitempty"`
	} `xml:"image"`
	Voice struct {
		Name string `xml:"name,attr" json:"-name"`
	} `xml:"voice"`
	Video []struct {
		Src       string `xml:"src,attr" json:"-src"`
		ClipEnd   string `xml:"clipEnd,attr" json:"-clipEnd,omitempty"`
		KeepFrame string `xml:"keepFrame,attr" json:"-keepFrame,omitempty"`
		ClipBegin string `xml:"clipBegin,attr" json:"-clipBegin,omitempty"`
	} `xml:"video"`
	ActiveDeck Deck `json:"activeDeck"`
}

type Deck struct {
	ID     string `xml:"id,attr" json:"-id,omitempty"`
	Src    string `xml:"src,attr" json:"-src,omitempty"`
	Active bool   `xml:"active,attr,omitempty" json:"-active,omitempty"`
}

/*
@todo error checking for empty info and lexeme tags, along with checks for deck page numbers.
*/
func getMeta(xmlPath string) (*LectureMeta, error) {

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	_ = byteValue

	var lecturemeta LectureMeta
	if err := xml.Unmarshal(byteValue, &lecturemeta); err != nil {
		log.Fatal(err)
	}

	for _, deck := range lecturemeta.Deck {
		if deck.Active {
			lecturemeta.ActiveDeck = deck
		}
	}

	return &lecturemeta, nil
}

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type AltLecture struct {
	XMLName  xml.Name  `xml:"lecture" json:"-"`
	Elements []Element `xml:",any"`
}

type Text struct {
	Type string
	Text string
}

type Element struct {
	Type       string `json:",omitempty"`
	Name       string
	Value      interface{}       `json:",omitempty"`
	Attributes map[string]string `json:",omitempty"`
	Elements   []Element         `xml:",any"`
}

func (m *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	default: //really shouldn't do this :p
		var e interface{}
		// if err := d.DecodeElement(&e, &start); err != nil {
		// 	return err
		// }
		fmt.Println(start, e)

	}
	return nil
}

func (m *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	default: //really shouldn't do this :p
		var e interface{}
		if err := d.DecodeElement(&e, &start); err != nil {
			return err
		}
		// j, _ := json.Marshal(start)
		// fmt.Println(string(j), e)
		m.Value = e
		mp := make(map[string]string)
		for _, attr := range start.Attr {
			mp[attr.Name.Local] = attr.Value
		}
		m.Attributes = mp
		if len(start.Attr) > 0 {
			m.Type = "element"
			m.Name = start.Name.Local
		} else {
			m.Type = start.Name.Local
		}

	}
	return nil
}

func parser() {

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	Inpath := filepath.Join(absPath, "input")
	xmlFile, err := os.Open(filepath.Join(Inpath, "lec.xml"))
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	_ = byteValue

	// var lecture Lecture
	// if err := xml.Unmarshal([]byte(blob), &lecture); err != nil {
	// 	log.Fatal(err)
	// }
	// j, _ := json.Marshal(lecture)
	// fmt.Println(lecture.Deck)
	// fmt.Println(string(j))

	altblob := `
	<lecture>
text
<slide/>
text2
<slide/>
text3
</lecture>
`
	_ = altblob
	var altlecture AltLecture
	if err := xml.Unmarshal(byteValue, &altlecture); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(altlecture)

	j, _ := json.Marshal(altlecture)
	fmt.Println(string(j))
}

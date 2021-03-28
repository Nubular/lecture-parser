package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nubular/lecture-parser/parser"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// Implement max page count using pdfinfo
// func main() {

// 	pdfName := "test.pdf"
// 	imageName := "test.jpg"

// 	if err := GetPDF(pdfName, imageName); err != nil {
// 		log.Fatal(err)
// 	}
// }

// Frame defines the page and it's imageName
type Frame struct {
	Page      int    `json:"page"`
	FileName  string `json:"imageName"`
	SSML      string `json:"ssml"`
	ImagePath string `json:"imagePath"`
	AudioPath string `json:"audioPath"`
}

// GetPDFPage extracts the specified page number from the supplied file
func GetPDFPage(pdfName string, imageName string, pageNum int) error {

	absPath, err := os.Getwd()
	if err != nil {
		return err
	}

	Inpath := filepath.Join(absPath, "input", pdfName)
	Outpath := filepath.Join(absPath, "output")

	if _, err := os.Stat(Outpath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", Outpath)
		os.Mkdir(Outpath, os.ModePerm)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(72, 72); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(Inpath); err != nil {
		return err
	}

	mw.SetIteratorIndex(pageNum - 1)

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(100); err != nil {
		return err
	}

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}

	// Save File
	return mw.WriteImage(filepath.Join(Outpath, imageName))
}

//GetPDFPages extracts the pages in the supplied array
func GetPDFPages(inPath string, outPath string, frames []parser.Section) error {

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", outPath)
		os.Mkdir(outPath, os.ModePerm)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(100, 100); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(inPath); err != nil {
		return err
	}

	// fmt.Println(mw.GetNumberImages())

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(100); err != nil {
		return err
	}

	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}
	var err error
	for _, frame := range frames {

		mw.SetIteratorIndex(frame.Page - 1)

		// if err := mw.SetImageFormat("jpeg"); err != nil {
		// 	return err
		// }

		// Save File
		err = mw.WriteImage(filepath.Join(outPath, frame.FrameSrc.ImageSrc))

		if err != nil {
			return fmt.Errorf("%s [PDF extraction: page: %d, loc: %s, to %s] ", err, frame.Page, inPath, filepath.Join(outPath, frame.FrameSrc.ImageSrc))
		}

	}
	return nil
}

// GetPDF will take a filename of a pdf file and convert the file into an
// image which will be saved back to the same location. It will save the image as a
// high resolution jpg file with minimal compression.
func GetPDF(pdfName string, imageName string) error {

	absPath, err := os.Getwd()
	if err != nil {
		return err
	}

	Inpath := filepath.Join(absPath, "input", pdfName)
	Outpath := filepath.Join(absPath, "output")

	if _, err := os.Stat(Outpath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", Outpath)
		os.Mkdir(Outpath, os.ModePerm)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(72, 72); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(Inpath); err != nil {
		return err
	}

	mw.SetIteratorIndex(0)

	log.Print("Converting ", pdfName, absPath)

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(50); err != nil {
		return err
	}

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}

	// Save File
	return mw.WriteImages(filepath.Join(Outpath, imageName), false)
}

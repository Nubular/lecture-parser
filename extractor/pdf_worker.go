package extractor

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {

	pdfName := "test.pdf"
	imageName := "test.jpg"

	if err := GetPDF(pdfName, imageName); err != nil {
		log.Fatal(err)
	}
}

/*
Implement Max page count using pdfinfo or imagemagick.
*/
// GetPDFPage extracts the specified page number from the supplied file
func GetPDFPage(pdfName string, imageName string, pageNum int) error {
	// path := filepath.Join("../../input", pdfName)
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

	// Must be *before* ReadImageFile
	if err := mw.SetResolution(72, 72); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(Inpath); err != nil {
		return err
	}

	mw.SetIteratorIndex(pageNum)

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
	return mw.WriteImage(filepath.Join(Outpath, imageName))
}

// GetPDF will take a filename of a pdf file and convert the file into an
// image which will be saved back to the same location. It will save the image as a
// high resolution jpg file with minimal compression.
func GetPDF(pdfName string, imageName string) error {
	// path := filepath.Join("../../input", pdfName)
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

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(72, 72); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(Inpath); err != nil {
		return err
	}

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

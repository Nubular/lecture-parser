package main

import (
	"fmt"
	"os"

	"gocv.io/x/gocv"
)

func pain() {
	window := gocv.NewWindow("Face Blur")
	defer window.Close()
	img := gocv.IMRead("FRAMES/CN1_1.jpg", gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Could not read image")
		os.Exit(1)
	}

	// blur := gocv.NewMat()
	// gocv.GaussianBlur(gray_image, &blur, image.Pt(75, 75), 0, 0, gocv.BorderDefault)
	window.IMShow(img)
	window.WaitKey(0)
}

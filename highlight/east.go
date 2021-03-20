package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"gocv.io/x/gocv"
)

func resize_image_blob(src gocv.Mat, max_side_len float64) (gocv.Mat, float64, float64) {
	h := float64(src.Size()[0])
	w := float64(src.Size()[1])
	var ratio float64

	if math.Max(h, w) > max_side_len {
		ratio = max_side_len / math.Max(h, w)
	} else {
		ratio = 1
	}

	resize_w := int(w * ratio)
	resize_h := int(h * ratio)

	resize_w = (resize_w / 32) * 32
	if resize_w < 32 {
		resize_w = 32
	}
	resize_h = (resize_h / 32) * 32
	if resize_h < 32 {
		resize_h = 32
	}

	scalar := gocv.NewScalar(123.68, 116.78, 103.94, 0)
	blob := gocv.BlobFromImage(src, 1.0, image.Pt(resize_w, resize_h), scalar, true, false)

	rw := w / float64(resize_w)
	rh := h / float64(resize_h)

	return blob, rh, rw
}

func east(detector gocv.Net, blob gocv.Mat, inpWidth int, inpHeight int) []gocv.Mat {
	outNames := []string{"feature_fusion/Conv_7/Sigmoid", "feature_fusion/concat_3"}
	detector.SetInput(blob, "")

	return detector.ForwardLayers(outNames)
	// return detector.Forward("")
}

// https://docs.opencv.org/3.4/d3/d63/classcv_1_1Mat.html
func decodeBoundingBoxes(scores, geometry gocv.Mat, scoreThresh int) {
	// dim := geometry.Size()[1]
	height := scores.Size()[2]
	width := scores.Size()[3]

	for i := 0; i < height; i++ {
		conf := make([]float32, 0)
		for j := 0; j < width; j++ {
			scoresData := scores.GetFloatAt(0, i*width+j)
			// x0_data := geometry.GetFloatAt(0, i*width+j)
			// x1_data := geometry.GetFloatAt(0, height*width+i*width+j)
			// x2_data := geometry.GetFloatAt(0, height*width*2+i*width+j)
			// x3_data := geometry.GetFloatAt(0, height*width*3+i*width+j)
			// angles_data := geometry.GetFloatAt(0, height*width*4+i*width+j)

			// log.Printf("[%f %f %f %f %f %f]", scoresData, x0_data, x1_data, x2_data, x3_data, angles_data)
			conf = append(conf, scoresData)
			// 	score := scoresData[j]

			// 	if score < scoreThresh {
			// 		continue
			// 	}

			// 	offset_x := j * 4
			// 	offset_y := i * 4
			// 	angle := angles_data[j]

			// 	cosA := math.Cos(angle)
			// 	sinA := math.Sin(angle)
		}
		fmt.Println(conf)
	}
	// yee, _ := geometry.DataPtrFloat32()
	// fmt.Println(yee)

}

func main() {
	// confThreshold := 0.5
	// nmsThreshold := 0.4
	inpWidth := 320
	inpHeight := 320
	// modelDetector := "frozen_east_text_detection.pb"
	// x_bias := 0.1
	// y_bias := 10
	// minD := 30
	// alpha := 0.5

	// alpha_slider_max := 100
	// x_bias_max := 100
	// y_bias_max := 10

	window := gocv.NewWindow("DNN Detection")
	defer window.Close()

	detector := gocv.ReadNet("frozen_east_text_detection.pb", "")

	// detector := gocv.ReadNet("model.caffemodel", "proto.txt")
	// ln := detector.GetLayerNames()
	// log.Println(ln, len(ln))

	img := gocv.IMRead("FRAMES/CN1_1.jpg", gocv.IMReadColor)
	defer img.Close()
	// return
	if img.Empty() {
		fmt.Println("Could not read image")
		os.Exit(1)
	}

	// height_ := img.Size()[0]
	// width_ := img.Size()[1]
	blob, rH, rW := resize_image_blob(img, 24000)
	defer blob.Close()
	_, _ = rW, rH
	start := time.Now()
	outs := east(detector, blob, inpWidth, inpHeight)
	elapsed := time.Since(start)
	log.Println("Inference time: ", elapsed)
	// performDetection(&img, outs)
	// window.IMShow(img)
	// window.WaitKey(0)
	scores := outs[0]
	geometry := outs[1]
	log.Println(scores.Size(), geometry.Size())
	decodeBoundingBoxes(scores, geometry, 0)
}

func performDetection(frame *gocv.Mat, results gocv.Mat) {
	conf := make([]float32, 0)
	for i := 0; i < results.Total(); i += 7 {
		confidence := results.GetFloatAt(0, i+2)
		conf = append(conf, confidence)
		if confidence > 0.5 {
			left := int(results.GetFloatAt(0, i+3) * float32(frame.Cols()))
			top := int(results.GetFloatAt(0, i+4) * float32(frame.Rows()))
			right := int(results.GetFloatAt(0, i+5) * float32(frame.Cols()))
			bottom := int(results.GetFloatAt(0, i+6) * float32(frame.Rows()))
			gocv.Rectangle(frame, image.Rect(left, top, right, bottom), color.RGBA{0, 255, 0, 0}, 2)
		}
	}
	log.Println(conf[0])
}

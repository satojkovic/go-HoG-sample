package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"math"

	"github.com/nfnt/resize"
)

const (
	FileName = "gmap_pin.jpg"
	ResizeX  = 30
	ResizeY  = 60
)

func ExtractHoG(img image.Image, imgw, imgh int) error {
	fmt.Println("--- Extract HoG Feature ---")

	// Resize image
	fmt.Println(" * Resize image")
	if imgw != ResizeX || imgh != ResizeY {
		img = resize.Resize(ResizeX, ResizeY, img, resize.Lanczos3)
	}

	// Convert to grayscale image
	fmt.Println(" * Convert to gray scale image")
	rect := image.Rect(0, 0, ResizeX, ResizeY)
	grayImg := image.NewGray(rect)
	for x := 0; x < ResizeX; x++ {
		for y := 0; y < ResizeY; y++ {
			oldColor := img.At(x, y)
			grayColor := img.ColorModel().Convert(oldColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	// Compute gradient
	fmt.Println(" * Compute gradient")
	gradImg := image.NewGray(rect)
	dirImg := image.NewGray(rect)
	for x := 1; x < ResizeX-1; x++ {
		for y := 1; y < ResizeY-1; y++ {
			stride := grayImg.Stride

			fu := float64(grayImg.Pix[(y*stride)+(x+1)] - grayImg.Pix[(y*stride)+(x-1)])
			fv := float64(grayImg.Pix[((y+1)*stride)+x] - grayImg.Pix[((y-1)*stride)+x])

			m := math.Sqrt(fu*fu + fv*fv)
			theta := 0.0
			if fu != 0.0 || fv != 0.0 {
				theta = (math.Atan(fv/fu) * 180.0 / math.Pi) + 90 // 0 - 180.0
			}

			gradImg.Pix[(y*gradImg.Stride)+x] = uint8(m)
			dirImg.Pix[(y*dirImg.Stride)+x] = uint8(theta)
		}
	}

	return nil
}

func main() {
	// Open image file
	b, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Fatal(err)
	}

	img, str, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(str)
		log.Fatal(err)
	}

	imgconf, _, err := image.DecodeConfig(bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	// Extract HoG feature
	err = ExtractHoG(img, imgconf.Width, imgconf.Height)
	if err != nil {
		log.Fatal(err)
	}

	// Show HoG Feature
}

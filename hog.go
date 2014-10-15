package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"

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
	if imgw != ResizeX || imgh != ResizeY {
		img = resize.Resize(ResizeX, ResizeY, img, resize.Lanczos3)
		fmt.Println("* Resized image")
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

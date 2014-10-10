package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
)

const (
	FileName = "gmap_pin.jpg"
)

type HoG struct {
	Feat  []float64
	Ndims uint64
}

func NewHoG(width, height int) *HoG {
	hog := &HoG{}
	return hog
}

func (hog *HoG) extract() error {
	fmt.Println("--- Extract HoG feature ---")

	return nil
}

func main() {
	// Open image file
	file, err := os.Open(FileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Successfully opened: %s\n", FileName)
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Image size: (%d, %d)\n", img.Width, img.Height)
	}

	// Compute HoG feature
	hog := NewHoG(img.Width, img.Height)
	err = hog.extract()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("--- Done ---")
	}

}

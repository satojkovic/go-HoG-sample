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

type Cell struct {
	Width, Height int
	NumGrad       int
}

type Block struct {
	Width, Height int
	NumX, NumY    int
	NumDim        int
	Cells         []Cell
}

type HoG struct {
	NumX, NumY int
	Blocks     []Block
	Descriptor []float64
}

func NewHoG(imgw, imgh int) *HoG {
	hog := &HoG{}
	hog.Initialize(imgw, imgh)
	return hog
}

func (self *HoG) Initialize(imgw, imgh int) {
}

func (self *Block) Initialize(imgw, imgh int) {
}

func (self *Cell) Initialize(imgw, imgh int) {
}

func (hog *HoG) Extract() error {
	fmt.Println("--- Extract HoG Feature ---")
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
	err = hog.Extract()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("--- Done ---")
	}
}

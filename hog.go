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
	Width, Height int // the number of pixels in a cell
	NumGrad       int
}

type Block struct {
	Width, Height int // the number of cells in a block
	NumDim        int
	Cells         []Cell // cells in a block
}

type HoG struct {
	NumCellX, NumCellY   int // the number of cells in an image
	NumBlockX, NumBlockY int // the number of blocks in an image
	NumDim               int
	Blocks               []Block // blocks in an image
	Descriptor           []float64
}

func NewHoG(imgw, imgh int) *HoG {
	hog := &HoG{}
	hog.Initialize(imgw, imgh)
	return hog
}

func (self *HoG) Initialize(imgw, imgh int) {
	cell := &Cell{}
	cell.Initialize()
	self.NumCellX = imgw / cell.Width
	self.NumCellY = imgh / cell.Height
	fmt.Printf("The number of pixels in a cell: (%d, %d)\n",
		cell.Width, cell.Height)
	fmt.Printf("The number of cells in an image: (%d, %d)\n",
		self.NumCellX, self.NumCellY)

	block := &Block{}
	block.Initialize(imgw, imgh, cell.NumGrad)
	self.NumBlockX = self.NumCellX - block.Width + 1
	self.NumBlockY = self.NumCellY - block.Height + 1
	fmt.Printf("The number of cells in a block: (%d, %d)\n",
		block.Width, block.Height)
	fmt.Printf("The number of blocks in an image: (%d, %d)\n",
		self.NumBlockX, self.NumBlockY)

	self.NumDim = self.NumBlockX * self.NumBlockY * block.NumDim
	fmt.Printf("Total dimensions: %d\n", self.NumDim)
}

func (self *Block) Initialize(imgw, imgh, celldim int) {
	self.Width = 3
	self.Height = 3
	self.NumDim = self.Width * self.Height * celldim
}

func (self *Cell) Initialize() {
	self.Width = 5
	self.Height = 5
	self.NumGrad = 9
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

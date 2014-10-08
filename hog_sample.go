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

func main() {
	// Open image file
	file, err := os.Open(FileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("successfully opened: %s\n", FileName)
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("image size: (%d, %d)\n", img.Width, img.Height)
	}

}

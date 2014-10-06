package main

import (
	"fmt"
	"os"
)
import "log"

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

}

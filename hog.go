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
	FileName  = "gmap_pin.jpg"
	ResizeX   = 30
	ResizeY   = 60
	CellSize  = 5 // [pixel]
	BinRange  = 20
	BlockSize = 3 // [cell]
	Epsilon   = 1.0
)

type Cell struct {
	Hist []float64
}

func NewCell() Cell {
	cell := Cell{}
	cell.Hist = make([]float64, 9)
	return cell
}

func resizeImage(img image.Image, imgw, imgh int) (image.Image, image.Rectangle) {
	if imgw != ResizeX || imgh != ResizeY {
		img = resize.Resize(ResizeX, ResizeY, img, resize.Lanczos3)
	}
	rect := image.Rect(0, 0, ResizeX, ResizeY)

	return img, rect
}

func convertGray(img image.Image, rect image.Rectangle) *image.Gray {
	grayImg := image.NewGray(rect)
	for y := 0; y < ResizeY; y++ {
		for x := 0; x < ResizeX; x++ {
			oldColor := img.At(x, y)
			grayColor := img.ColorModel().Convert(oldColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

func computeCellHist(grayImg *image.Gray) [][]Cell {
	fu, fv := 0.0, 0.0
	cells := make([][]Cell, (ResizeX / CellSize))
	for ci := 0; ci < (ResizeX / CellSize); ci++ {
		cells[ci] = make([]Cell, (ResizeY / CellSize))
	}
	for cy := 0; cy < (ResizeY / CellSize); cy++ {
		for cx := 0; cx < (ResizeX / CellSize); cx++ {
			cells[cx][cy] = NewCell()
		}
	}

	for y := 0; y < ResizeY; y++ {
		for x := 0; x < ResizeX; x++ {
			stride := grayImg.Stride

			if x == 0 {
				fu = float64(grayImg.Pix[(y*stride)+(x+1)] - grayImg.Pix[(y*stride)+(x)])
			} else if x == ResizeX-1 {
				fu = float64(grayImg.Pix[(y*stride)+(x)] - grayImg.Pix[(y*stride)+(x-1)])
			} else {
				fu = float64(grayImg.Pix[(y*stride)+(x+1)] - grayImg.Pix[(y*stride)+(x-1)])
			}

			if y == 0 {
				fv = float64(grayImg.Pix[((y+1)*stride)+x] - grayImg.Pix[((y)*stride)+x])
			} else if y == ResizeY-1 {
				fv = float64(grayImg.Pix[((y)*stride)+x] - grayImg.Pix[((y-1)*stride)+x])
			} else {
				fv = float64(grayImg.Pix[((y+1)*stride)+x] - grayImg.Pix[((y-1)*stride)+x])
			}

			m := math.Sqrt(fu*fu + fv*fv)
			theta := 0.0
			if fu != 0.0 || fv != 0.0 {
				theta = (math.Atan(fv/fu) * 180.0 / math.Pi) + 90 // 0 - 180.0
			}

			// Cell histogram
			bin := int(theta / BinRange)
			if bin == 9 {
				bin -= 1
			}
			cells[int(x/CellSize)][int(y/CellSize)].Hist[bin] += m
		}
	}

	return cells
}

func computeBlockNorm(cells [][]Cell) []float64 {
	hog := make([]float64, 3240)
	hogidx := 0
	cellnumx, cellnumy := ResizeX/CellSize, ResizeY/CellSize
	for cy := 0; cy < cellnumy; cy++ {
		for cx := 0; cx < cellnumx; cx++ {

			if cx+2 >= cellnumx || cy+2 >= cellnumy {
				continue
			}

			v := blockL2Norm(cx, cy, cellnumx, cells)

			// block normalization
			for iny := 0; iny < BlockSize; iny++ {
				for inx := 0; inx < BlockSize; inx++ {
					for b := 0; b < 9; b++ {
						val := cells[cx][cy].Hist[b]
						hog[hogidx] = val / math.Sqrt(v+Epsilon*Epsilon)
						hogidx++
					}
				}
			}
		}
	}

	return hog
}

func ExtractHoG(img image.Image, imgw, imgh int) ([]float64, error) {
	fmt.Println("--- Extract HoG Feature ---")

	// Resize image
	fmt.Println(" * Resize image")
	img, rect := resizeImage(img, imgw, imgh)

	// Convert to grayscale image
	fmt.Println(" * Convert to gray scale image")
	grayImg := convertGray(img, rect)

	// Compute cell histogram
	fmt.Println(" * Compute gradient")
	cells := computeCellHist(grayImg)

	// Compute block normalization
	fmt.Println(" * Compute block normalization")
	hog := computeBlockNorm(cells)

	return hog, nil
}

func blockL2Norm(cx, cy, cellnumx int, cells [][]Cell) float64 {
	v := 0.0
	for iny := 0; iny < BlockSize; iny++ {
		for inx := 0; inx < BlockSize; inx++ {
			for b := 0; b < 9; b++ {
				val := cells[cx+inx][cy+iny].Hist[b]
				v += (val * val)
			}
		}
	}

	return v
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
	hog, err := ExtractHoG(img, imgconf.Width, imgconf.Height)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- Extracted HoG Feature ---")
	fmt.Println(" * Len:", len(hog))
}

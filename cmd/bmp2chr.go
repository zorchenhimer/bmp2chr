package main

import (
	"path/filepath"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

    ".."
)


//const inputFilename string = "before.bmp"

func main() {
    var doubleHigh bool
	var inputFilename string
	var outputFilename string

	flag.StringVar(&inputFilename, "i", "", "Input BMP file")
	flag.StringVar(&outputFilename, "o", "", "Output filename (optional)")
	flag.BoolVar(&doubleHigh, "16", false, "8x16 tiles")
	flag.Parse()

	if len(inputFilename) == 0 {
		fmt.Println("Missing input file")
		os.Exit(1)
	}

	// Default the same name but with .chr extension
	if len(outputFilename) == 0 {
		outputFilename = inputFilename
		ext := filepath.Ext(inputFilename)
		outputFilename = outputFilename[0: len(outputFilename) - len(ext)] + ".chr"
	}

	if doubleHigh {
		fmt.Println("8x16 tiles are not yet supported")
		os.Exit(1)
	}

	// Read input file
	rawBmp, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse some headers
	fileHeader, err := bmp2chr.ParseFileHeader(rawBmp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	imageHeader, err := bmp2chr.ParseImageHeader(rawBmp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Validate image dimensions
	if imageHeader.Width != 128 {
		fmt.Println("Image width must be 128")
		os.Exit(1)
	}

	if imageHeader.Height%8 != 0 {
		fmt.Println("Image height must be a multiple of 8")
		os.Exit(1)
	}

	// Isolate the pixel data
	rawBmpPixels := rawBmp[fileHeader.Offset:len(rawBmp)]

	// Invert rows; They're stored top to bottom in BMP
	row := (len(rawBmpPixels) / 128) - 1
	uprightRows := []byte{}

	for row > -1 {
		// Get the row
		rawRow := rawBmpPixels[row*128 : row*128+128]
		// normalize each pixel's palette index
		for _, b := range rawRow {
			uprightRows = append(uprightRows, byte(int(b)%4))
		}
		row--
	}

	// split out the 8x8 or 8x16 tiles
	tileID := 0
	tiles := []*bmp2chr.RawTile{}
	numRows := 8
	if doubleHigh {
		numRows = 16
	}

	for tileID < (len(uprightRows) / 64) {
		// The first pixel offset in the current tile
		startOffset := (tileID/16)*(128*8) + (tileID%16)*8
		if doubleHigh {
			// From SlashLife for 8x16 tiles: lookupTileId = (tileId / 32) * 32 + (tileId % 32) / 2 + (tileId % 2) * 16
			// TODO: this isn't tested
			startOffset = (tileID / 32) * 32 + (tileID % 32) / 2 + (tileID % 2) * 16
		}

		tileBytes := []byte{}
		for y := 0; y < numRows; y++ {
			for x := 0; x < 8; x++ {
				tileBytes = append(tileBytes, uprightRows[startOffset+x+128*y])
			}
		}
		tiles = append(tiles, &bmp2chr.RawTile{Data: tileBytes})
		tileID++
	}

	chrFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer chrFile.Close()

	for _, tile := range tiles {
		_, err = chrFile.Write(tile.ToChr(doubleHigh))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}


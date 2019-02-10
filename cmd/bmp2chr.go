package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zorchenhimer/bmp2chr"
)

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
		outputFilename = outputFilename[0:len(outputFilename)-len(ext)] + ".chr"
	}

	if doubleHigh {
		fmt.Println("8x16 tiles are not yet supported")
		os.Exit(1)
	}

	bitmap, err := bmp2chr.OpenBitmap(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Invert rows; They're stored top to bottom in BMP
	row := (len(bitmap.Data) / 128) - 1
	uprightRows := []byte{}

	for row > -1 {
		// Get the row
		rawRow := bitmap.Data[row*128 : row*128+128]
		// normalize each pixel's palette index
		for _, b := range rawRow {
			uprightRows = append(uprightRows, byte(int(b)%4))
		}
		row--
	}

	// split out the 8x8 or 8x16 tiles
	tileID := 0
	tiles := []bmp2chr.Tile{}
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
			startOffset = (tileID/32)*32 + (tileID%32)/2 + (tileID%2)*16
		}

		var tileBytes bmp2chr.Tile
		for y := 0; y < numRows; y++ {
			for x := 0; x < 8; x++ {
				//tileBytes = append(tileBytes, uprightRows[startOffset+x+128*y])
				tileBytes[x+(8*y)] = uprightRows[startOffset+x+128*y]
			}
		}

		tiles = append(tiles, tileBytes)
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

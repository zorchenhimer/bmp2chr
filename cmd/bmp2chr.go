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

	bitmap, err := bmp2chr.OpenBitmap(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rect := bitmap.Rect()

	// Invert row order. They're stored top to bottom in BMP.
	uprightRows := []byte{}
	for row := (len(bitmap.Data) / rect.Max.X) - 1; row > -1; row-- {
		// Get the row
		rawRow := bitmap.Data[row*rect.Max.X : row*rect.Max.X+rect.Max.X]

		// normalize each pixel's palette index
		for _, b := range rawRow {
			uprightRows = append(uprightRows, byte(int(b)%4))
		}
	}

	// Cut out the 8x8 tiles
	tileID := 0
	tiles := []bmp2chr.Tile{}
	//fmt.Printf("uprightRows length: %d (%d)\n", len(uprightRows), len(uprightRows)/64)

	for tileID < (len(uprightRows) / 64) {
		// The first pixel offset in the current tile
		startOffset := (tileID/16)*(rect.Max.X*8) + (tileID%16)*8
		//fmt.Printf("tileID: %d startOffset: %d\n", tileID, startOffset)

		var tileBytes *bmp2chr.Tile
		tileBytes = bmp2chr.NewTile(tileID)
		for y := 0; y < 8; y++ {
			tileY := y

			// Wrap rows at 8 pixels
			if tileY >= 8 {
				tileY -= 8
			}

			// Get the pixels for the row.
			for x := 0; x < 8; x++ {
				tileBytes.Pix[x+(8*tileY)] = uprightRows[startOffset+x+rect.Max.X*y]
			}
		}

		tiles = append(tiles, *tileBytes)
		tileID++
	}

	chrFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer chrFile.Close()

	// If it's 8x16 mode, transform tiles.  Tiles on odd rows will be put after
	// the tile directly above them. The first four tiles would be $00, $10, $01, $11.
	if doubleHigh {
		newtiles := []bmp2chr.Tile{}
		for i := 0; i < len(tiles)/2; i++ {
			if i%16 == 0 && i > 0 {
				i += 16
			}

			newtiles = append(newtiles, tiles[i])
			newtiles = append(newtiles, tiles[i+16])
		}
		tiles = newtiles
	}

	for _, tile := range tiles {
		tchr := tile.ToChr()
		_, err = chrFile.Write(tchr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//fmt.Printf("number of tiles: %d\n", len(tiles))
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zorchenhimer/bmp2chr"
)

var supportedInputFormats []string = []string{".bmp"}

func main() {
	var doubleHigh bool
	var outputFilename string
	var debug bool

	//flag.StringVar(&inputFilename, "i", "", "Input BMP file")
	flag.StringVar(&outputFilename, "o", "", "Output filename")
	flag.BoolVar(&doubleHigh, "16", false, "8x16 tiles")
	flag.BoolVar(&debug, "debug", false, "Debug printing")
	flag.Parse()

	fileList := []string{}

	if len(flag.Args()) > 0 {
		for _, target := range flag.Args() {
			found, err := filepath.Glob(target)
			if err == nil && len(found) > 0 {
				fileList = append(fileList, found...)
			} else {
				fmt.Printf("%q not found\n", target)
				os.Exit(1)
			}
		}
	}

	if len(fileList) == 0 {
		fmt.Println("Missing input file(s)")
		os.Exit(1)
	}

	// Require an output filename if there's more than one input.
	if len(outputFilename) == 0 {
		if len(fileList) == 1 {
			outputFilename = fileList[0]
			ext := filepath.Ext(fileList[0])
			outputFilename = outputFilename[0:len(outputFilename)-len(ext)] + ".chr"
		} else {
			fmt.Println("Missing output filename")
			os.Exit(1)
		}
	}

	for _, file := range fileList {
		ext := filepath.Ext(file)
		found := false
		for _, supp := range supportedInputFormats {
			if ext == supp {
				found = true
			}
		}
		if !found {
			fmt.Printf("Unsupported input format for file %q\n", file)
			os.Exit(1)
		}
	}

	chrFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer chrFile.Close()

	for _, inputfile := range fileList {
		bitmap, err := bmp2chr.OpenBitmap(inputfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if debug {
			err := ioutil.WriteFile("upright.dat", bitmap.RawData, 0777)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// If it's 8x16 mode, transform tiles.  Tiles on odd rows will be put after
		// the tile directly above them. The first four tiles would be $00, $10, $01, $11.
		if doubleHigh {
			newtiles := []bmp2chr.Tile{}
			for i := 0; i < len(bitmap.Tiles)/2; i++ {
				if i%bitmap.TilesPerRow == 0 && i > 0 {
					i += bitmap.TilesPerRow
				}

				newtiles = append(newtiles, bitmap.Tiles[i])
				newtiles = append(newtiles, bitmap.Tiles[i+bitmap.TilesPerRow])
			}
			bitmap.Tiles = newtiles
		}

		for _, tile := range bitmap.Tiles {
			if debug {
				fmt.Println(tile.ASCII())
			}

			tchr := tile.ToChr()
			_, err = chrFile.Write(tchr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	//fmt.Printf("number of tiles: %d\n", len(tiles))
}

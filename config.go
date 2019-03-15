package bmp2chr

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	InputFile	string
	Indices	[]int
}

//func ReadConfig(filename string) (*Configuration, error) {
func ReadConfig(filename string) (*Bitmap, error) {
	rawjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Configuration
	if err := json.Unmarshal(rawjson, &config); err != nil {
		return nil, err
	}

	bitmap, err := OpenBitmap(config.InputFile)
	if err != nil {
		return nil, err
	}

	bitmap.Config = &config
	return bitmap, nil
}

// exists returns whether the given file or directory exists or not.
// Taken from https://stackoverflow.com/a/10510783
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

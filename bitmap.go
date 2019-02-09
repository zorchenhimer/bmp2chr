package bmp2chr

import (
    "encoding/binary"
    "fmt"
)

type FileHeader struct {
	Size   int // size of file in bytes
	Offset int // offset to start of pixel data
}

func (f FileHeader) String() string {
	return fmt.Sprintf("Size: %d Offset: %d", f.Size, f.Offset)
}

// Size, offset, error
func ParseFileHeader(input []byte) (*FileHeader, error) {
	if len(input) < 4 {
		return nil, fmt.Errorf("Data too short for header")
	}
	header := input[0:14]

	size := binary.LittleEndian.Uint32(header[2:6])
	offset := binary.LittleEndian.Uint32(header[10:14])
	return &FileHeader{Size: int(size), Offset: int(offset)}, nil
}

type ImageHeader struct {
	headerSize  int
	Width       int
	Height      int
	BitDepth    int
	Compression int
	Size        int // image size

	// "Pixels per meter"
	ppmX int
	ppmY int

	ColorMapEntries   int
	SignificantColors int
}

func (i ImageHeader) String() string {
	return fmt.Sprintf("(%d, %d) %d bpp @ %d bytes", i.Width, i.Height, i.BitDepth, i.Size)
}

func ParseImageHeader(input []byte) (*ImageHeader, error) {
	if len(input) < (14 + 12) {
		return nil, fmt.Errorf("Data too short for image header")
	}

	header := &ImageHeader{}
	header.headerSize = int(binary.LittleEndian.Uint32(input[14:18]))

	//headerRaw := input[14 : 14+header.Size]

	header.Width = int(binary.LittleEndian.Uint32(input[18:22]))
	header.Height = int(binary.LittleEndian.Uint32(input[22:26]))
	header.BitDepth = int(binary.LittleEndian.Uint16(input[28:30]))

	header.Size = int(binary.LittleEndian.Uint32(input[38:42]))

	return header, nil
}

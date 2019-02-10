package bmp2chr

// Assembled meta sprites and tiles.  These will be unwrapped to the specified
// layout (eg, 8*16 vs 8x8)
type MetaTile struct {
	Tiles []Tile

	// Width and Hight in tiles, not pixels
	Width  int
	Height int

	// Layout of tiles in the destination CHR
	Layout TileLayout
}

// Data is a list of palette indexes.  One ID per pixel.  A single tile is
// always 8x8 pixels.  Larger meta tiles (eg, 8*16) will be made up of multiple
// tiles of 64 total pixels.
type Tile [64]byte

// Ideally, each tile or object will be in its own input file and is assembled
// into the final CHR layout during assemble time.
type TileLayout int

const (
	TL_SINGLE = iota // Default.  A single 8x8 tile.
	TL_8X16          // 8x16 sprites.
	TL_ROW           // Row sequential
	TL_COLUMN        // Column sequential
	TL_ASIS          // Don't transform.  This will break things if there's meta tiles that are not the same size.
)

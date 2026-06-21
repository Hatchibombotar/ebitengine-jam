package main

type Tile struct {
	// "wall", "box"
	Type string
}

type Tilemap = [15][20]*Tile

func TileIsSolid(tile *Tile) bool {
	if tile == nil {
		return false
	}
	return tile.Type == "wall"
}

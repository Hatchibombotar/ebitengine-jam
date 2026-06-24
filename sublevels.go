package main

import (
	"encoding/json"
	"fmt"
)

type Tile struct {
	// "wall", "box"
	Type string `json:"type"`
	// used to identify different machines.
	SubType string `json:"subtype"`
}

type Tilemap = [15][20]*Tile

func TileIsSolid(tile *Tile) bool {
	if tile == nil {
		return false
	}
	return tile.Type == "wall" || tile.Type == "machine"
}

func TileIsAccended(tile *Tile) bool {
	if tile == nil {
		return false
	}
	return tile.Type == "box" || tile.Type == "conveyor_left" || tile.Type == "conveyor_down"
}

func PrintSublevel(g *Game) {
	v, err := json.Marshal(g.CurrentSublevel().tileMap)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	fmt.Println(string(v))
}

func LoadSubLevel(s string) *Tilemap {
	var t *Tilemap = &[15][20]*Tile{}
	err := json.Unmarshal([]byte(s), t)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return t
}

func createSublevels() map[string]*Sublevel {
	return map[string]*Sublevel{
		"sewer": {
			tileMap: LoadSubLevel(`[[{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null]]`),
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "tape"}, X: 10, Y: 3},
				{itemType: &Item{id: "tape"}, X: 5, Y: 6},
			},
		},
		"factory": {
			tileMap: LoadSubLevel(`[[{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"}]]`),
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "tape"}, X: 10, Y: 3},
				{itemType: &Item{id: "tape"}, X: 5, Y: 6},
			},
			Enemies: []*Enemy{
				createEnemy1(),
			},
		},
	}
}

type Sublevel struct {
	tileMap     *Tilemap
	inGameItems []*InGameItem

	conveyorItems []*ConveyorItem

	Enemies []*Enemy
}

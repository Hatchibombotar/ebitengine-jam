package main

import (
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	// "wall", "box", "vent_down"
	Type string `json:"type"`
	// used to identify different machines.
	SubType string `json:"subtype"`

	Inventory map[string]int

	Damage float64
}

type Sublevel struct {
	Title string

	tileMap     *Tilemap
	inGameItems []*InGameItem

	conveyorItems []*ConveyorItem

	adjacentSpaces AdjacentSpaces

	Enemies []*Enemy

	Background *ebiten.Image
	Overlay    *ebiten.Image
	Vignette   *ebiten.Image

	Spawners []*Spawner
}

type AdjacentSpaces struct {
	North string
	East  string
	South string
	West  string
	Up    string
	Down  string
}

type Spawner struct {
	X    int
	Y    int
	item string
}

type Tilemap = [15][20]*Tile

func TileIsSolid(tile *Tile) bool {
	if tile == nil {
		return false
	}
	return tile.Type == "wall" || tile.Type == "machine" || tile.Type == "electrical_switch"
}

func GetTileFromTileMap(tileMap *Tilemap, x, y int) *Tile {
	if x < 0 || x >= 20 || y < 0 || y >= 15 {
		return nil
	}
	return tileMap[y][x]
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

func LoadTileMapJSON(s string) *Tilemap {
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
		"sewer_0": {
			Title: "Sewer: West 2",
			// tileMap: LoadSubLevel(`[[{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall"}],[{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"},{"type":"wall"}],[{"type":"wall"},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null]]`),
			tileMap: LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null]]`),
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "string"}, X: 10, Y: 3},
				{itemType: &Item{id: "string"}, X: 5, Y: 6},
			},
			adjacentSpaces: AdjacentSpaces{
				West: "final_assembly",
				Up:   "final_assembly",
				East: "sewer_1",
			},
			Background: sewer_left,
			Vignette:   vignette,
		},
		"sewer_1": {
			Title:       "Sewer: West 1",
			tileMap:     LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null]]`),
			inGameItems: []*InGameItem{},
			adjacentSpaces: AdjacentSpaces{
				West: "sewer_0",
				East: "sewer_2",
			},

			Background: sewer_middle,
			Vignette:   vignette,
		},
		"sewer_2": {
			Title:       "Sewer: Center",
			tileMap:     LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null]]`),
			inGameItems: []*InGameItem{},
			adjacentSpaces: AdjacentSpaces{
				West:  "sewer_1",
				North: "sewer_3",
				East:  "sewer_5",
			},

			Background: sewer_middle_with_top_connection,
			Vignette:   vignette,
		},
		"sewer_3": {
			Title:       "Sewer: North 1",
			tileMap:     LoadTileMapJSON(`[[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null]]`),
			inGameItems: []*InGameItem{},
			adjacentSpaces: AdjacentSpaces{
				South: "sewer_2",
				North: "sewer_4",
			},

			Background: sewer_going_up,
			Vignette:   vignette,
		},
		"sewer_4": {
			Title:       "Sewer: North 2",
			tileMap:     LoadTileMapJSON(`[[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null]]`),
			inGameItems: []*InGameItem{},
			adjacentSpaces: AdjacentSpaces{
				South: "sewer_3",
				Up:    "storage_b",
			},

			Background: sewer_top,
			Vignette:   vignette,
		},

		"sewer_5": {
			Title:   "Sewer: East 1",
			tileMap: &[15][20]*Tile{},
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "string"}, X: 10, Y: 3},
				{itemType: &Item{id: "string"}, X: 5, Y: 6},
			},
			adjacentSpaces: AdjacentSpaces{
				West: "sewer_2",
				East: "sewer_entrance",
			},
			Background: sewer_middle,
			Vignette:   vignette,
		},

		"sewer_entrance": {
			Title:   "Sewer: Entrance",
			tileMap: &[15][20]*Tile{},
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "string"}, X: 10, Y: 3},
				{itemType: &Item{id: "string"}, X: 5, Y: 6},
			},
			adjacentSpaces: AdjacentSpaces{
				West: "sewer_5",
				Up:   "factory_floor",
			},
			Background: sewer_entrance,
			Vignette:   vignette,
		},
		"factory_floor": {
			Title:       "Factory Floor Q",
			tileMap:     LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"vent_down","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				Down: "sewer_entrance",
			},
			Background: metal_room,
			Overlay:    metal_room_overlay,
			Vignette:   vignette_mild,
			Spawners: []*Spawner{
				{
					X: 19, Y: 7, item: "metal_sheet",
				},
			},
		},
		"final_assembly": {
			Title:       "Final Assembly",
			tileMap:     LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"conveyor_down","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"seal_board_in_casing","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"program_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				East: "component_pnp",
				Down: "sewer_0",
			},
			Background: final_production,
			Overlay:    final_production_overlay,
			Vignette:   vignette_mild,
			Spawners: []*Spawner{
				{
					X: 19, Y: 7, item: "circuit_board_finished",
				},
				{
					X: 5, Y: 0, item: "casing",
				},
			},
		},
		"component_pnp": {
			Title:       "Component Pick and Place",
			tileMap:     LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"conveyor_down","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"conveyor_down","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"conveyor_down","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"conveyor_down","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"add_component_to_finished_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"add_component_to_finished_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"add_component_to_finished_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"add_component_to_finished_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				West:  "final_assembly",
				East:  "pcb_manufacture",
				South: "electrical_corridor",
			},
			Background: pick_and_place,
			Overlay:    pick_and_place_overlay,
			Vignette:   vignette_mild,
			Spawners: []*Spawner{
				{
					X: 19, Y: 7, item: "circuit_board_finished",
				},
				{
					X: 5, Y: 0, item: "battery",
				},
				{
					X: 8, Y: 0, item: "led",
				},
				{
					X: 11, Y: 0, item: "antenna",
				},
				{
					X: 14, Y: 0, item: "chip",
				},
			},
		},
		"pcb_manufacture": {
			Title:   "PCB Manufacture",
			tileMap: LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"conveyor_down","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"conveyor_down","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"apply_template_to_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"machine","subtype":"combine_copper_and_resin_board","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null},{"type":"conveyor_left","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),

			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				West: "component_pnp",
			},
			Background: pcb_manufacture,
			Overlay:    pcb_manufacture_overlay,
			Vignette:   vignette_mild,
			Spawners: []*Spawner{
				{
					X: 19, Y: 7, item: "resin_board",
				},
				{
					X: 14, Y: 0, item: "copper_sheet",
				},
			},
		},
		"storage_b": {
			Title:       "Storage B",
			tileMap:     LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"vent_down","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"vent_down","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				Down:  "sewer_4",
				North: "control_room",
			},
			Background: storage_b,
			Overlay:    storage_b_overlay,
			Vignette:   vignette_mild,
		},
		"control_room": {
			Title:   "Control Room",
			tileMap: LoadTileMapJSON(`[[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			// tileMap: &[15][20]*Tile{},
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "string"}, X: 10, Y: 3},
				{itemType: &Item{id: "string"}, X: 5, Y: 6},
			},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				South: "storage_b",
			},
			Background: control_room,
			Overlay:    control_room_overlay,
			Vignette:   vignette_mild,
		},
		"electrical_corridor": {
			Title:       "Electrical Corridor",
			tileMap:     LoadTileMapJSON(`[[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				North: "component_pnp",
				East:  "electrical",
			},
			Background: electrical_corridor,
			// Overlay:    control_room_overlay,
			Vignette: vignette_mild,
		},
		"electrical": {
			Title:       "Electrical Room",
			tileMap:     LoadTileMapJSON(`[[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"electrical_switch","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,{"type":"wall","subtype":"","Inventory":null}],[{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null},{"type":"wall","subtype":"","Inventory":null}]]`),
			inGameItems: []*InGameItem{},
			Enemies: []*Enemy{
				createEnemy1(),
			},
			adjacentSpaces: AdjacentSpaces{
				West: "electrical_corridor",
			},
			Background: electrical,
			Overlay:    electrical_overlay,
			Vignette:   vignette_mild,
		},
	}
}

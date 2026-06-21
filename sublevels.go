package main

func openCraftingUI(g *Game) {
	g.selectedRecipe = -1
}

func createSublevels() map[string]*Sublevel {
	return map[string]*Sublevel{
		"sewer": {
			tileMap: [15][20]*Tile{},
			inGameItems: []*InGameItem{
				{itemType: &Item{id: "tape"}, X: 10, Y: 3},
				{itemType: &Item{id: "tape"}, X: 5, Y: 6},
			},
		},
	}
}

type Sublevel struct {
	tileMap     Tilemap
	inGameItems []*InGameItem
}

package main

import "github.com/hajimehoshi/ebiten/v2"

type Item struct {
	id string
}

type ItemData struct {
	name  string
	image *ebiten.Image
}

type Recipe struct {
	result      string
	ingredients []string
}

var itemData map[string]*ItemData
var recipeData []*Recipe

func dropItemSlot(g *Game, dropItemSlot int) {
	slot := g.inventory[dropItemSlot]
	if slot == nil {
		return
	}

	g.CurrentSublevel().inGameItems = append(g.CurrentSublevel().inGameItems, &InGameItem{
		X:        g.player.position.X + g.player.facingDirection.X*0.8,
		Y:        g.player.position.Y + g.player.facingDirection.Y*0.8,
		itemType: &Item{id: slot.id},
	})

	g.inventory[dropItemSlot] = nil
}

func init() {
	itemData = map[string]*ItemData{
		"tape": {
			name:  "Tape",
			image: item_tape,
		},
		"bundle": {
			name:  "Bundle",
			image: item_bundle,
		},
		"box": {
			name:  "Box",
			image: item_box,
		},
		"holy_grail": {
			name:  "Holy Grain",
			image: target,
		},
	}

	recipeData = []*Recipe{
		{
			result:      "bundle",
			ingredients: []string{"tape", "tape"},
		},
		{
			result:      "bundle",
			ingredients: []string{"tape", "bundle"},
		},
		{
			result:      "tape",
			ingredients: []string{"bundle", "bundle"},
		},
		{
			result:      "tape",
			ingredients: []string{"tape", "tape", "tape"},
		},
		{
			result:      "tape",
			ingredients: []string{"tape", "tape"},
		},
		{
			result:      "tape",
			ingredients: []string{"tape", "holy_grail"},
		},
	}
}

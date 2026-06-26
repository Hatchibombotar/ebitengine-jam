package main

import "github.com/hajimehoshi/ebiten/v2"

type Item struct {
	id string
}

type ItemData struct {
	name     string
	image    *ebiten.Image
	heldItem *ebiten.Image
}

func GetHeldItemSprite(i *ItemData) *ebiten.Image {
	if i.heldItem != nil {
		return i.heldItem
	} else {
		return i.image
	}
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
		"string": {
			name:  "String",
			image: item_string,
		},
		"rod": {
			name:  "Rod",
			image: item_rod,
		},
		"screwdriver": {
			name:  "Screwdriver",
			image: item_screwdriver,
		},
		"hammer": {
			name:  "Hammer",
			image: item_hammer,
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
		"circuit_board_finished": {
			name:  "Circuit Board",
			image: circuit_board_finished,
		},
		"circuit_board_programmed": {
			name:  "Circuit Board (programmed)",
			image: circuit_board_finished,
		},
		"copper_sheet": {
			name:  "Copper Sheet",
			image: copper_sheet,
		},
		"metal_sheet": {
			name:     "metal Sheet",
			image:    metal_sheet,
			heldItem: metal_sheet_held,
		},
		"resin_board": {
			name:  "Resin Sheet",
			image: resin_board,
		},
		"uncut_circuit_board": {
			name:  "Raw Circuit Board",
			image: circuit_board_uncut,
		},
		"battery": {
			name:  "Battery",
			image: battery,
		},
		"led": {
			name:  "LED",
			image: led,
		},
		"chip": {
			name:  "Chip",
			image: chip,
		},
		"antenna": {
			name:  "Antenna",
			image: antenna,
		},
		"casing": {
			name:  "Unsealed Casing",
			image: casing,
		},
		"final_chip": {
			name:  "Evil Chip",
			image: final_chip,
		},
	}

	recipeData = []*Recipe{
		{
			result:      "screwdriver",
			ingredients: []string{"string", "rod"},
		},
		{
			result:      "hammer",
			ingredients: []string{"string", "rod", "metal_sheet"},
		},
		{
			result:      "bundle",
			ingredients: []string{"string", "bundle"},
		},
		{
			result:      "string",
			ingredients: []string{"bundle", "bundle"},
		},
		{
			result:      "string",
			ingredients: []string{"string", "string", "string"},
		},
		{
			result:      "string",
			ingredients: []string{"string", "string"},
		},
		{
			result:      "string",
			ingredients: []string{"string", "holy_grail"},
		},
	}
}

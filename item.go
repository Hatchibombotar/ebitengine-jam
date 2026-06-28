package main

import "github.com/hajimehoshi/ebiten/v2"

type Item struct {
	id         string
	resultData string
}

// template
// "Reprogrammer", "Hacking USB", "Mind Control"

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
	resultData  string
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
		"auth_chip": {
			name:  "Authentication Chip",
			image: item_auth_chip,
		},
		"auth_card": {
			name:  "Access Card",
			image: item_auth_card,
		},
		"wire_cutters": {
			name:  "Wire Cutters",
			image: item_wire_cutters,
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
		"broken_chip": {
			name:  "Empty Chip",
			image: broken_chip,
		},
		"template_machine": {
			name:  "Template Rewriter",
			image: item_template_machine,
		},
		"template": {
			name:  "Template",
			image: item_template,
		},
		"hacking_usb": {
			name:  "Hacking USB",
			image: item_hacking_usb,
		},
		"hacking_chip": {
			name:  "Hacking Chip",
			image: item_hacking_chip,
		},
		"reprogramming_chip": {
			name:  "Reprogramming Chip",
			image: reprogramming_chip,
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
			result:      "wire_cutters",
			ingredients: []string{"string", "rod", "rod", "metal_sheet"},
		},
		{
			result:      "auth_card",
			ingredients: []string{"auth_chip", "metal_sheet"},
		},
		{
			result:      "template_machine",
			ingredients: []string{"led", "chip", "copper_sheet"},
		},
		{
			result:      "template",
			resultData:  "Hacking USB",
			ingredients: []string{"template_machine", "template"},
		},
		{
			result:      "template",
			resultData:  "Reprogrammer",
			ingredients: []string{"template_machine", "template"},
		},
		{
			result:      "hacking_usb",
			ingredients: []string{"hacking_chip", "metal_sheet"},
		},
	}
}

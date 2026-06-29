package main

type ConveyorItem struct {
	// Position aligned to grid
	X, Y     float64
	itemType *Item
}

func UpdateConveyors(g *Game, sublevel *Sublevel) {
	if g.t%10 != 0 {
		return
	}

	for i, item := range sublevel.conveyorItems {
		tileX, tileY := int(item.X), int(item.Y)
		// var nextSpaceOccupied = false
		// for _, otherItem := range sublevel.conveyorItems {
		// 	if item == otherItem {
		// 		continue
		// 	}
		// 	otherTileX, otherTileY := int(item.X), int(item.Y)
		// 	if tileX == (otherTileX-1) && tileY == otherTileY {
		// 		nextSpaceOccupied = true
		// 	}
		// }
		// if nextSpaceOccupied {
		// 	continue
		// }

		if tileX < 0 {
			adjacentSpaceId := sublevel.adjacentSpaces.West
			if adjacentSpaceId != "" {
				adjacentSpace := g.spaces[adjacentSpaceId]

				item.X = 19

				adjacentSpace.conveyorItems = append(adjacentSpace.conveyorItems, item)
				sublevel.conveyorItems[i] = nil
			} else {
				sublevel.conveyorItems[i] = nil
			}
			continue
		}

		if sublevel.tileMap[tileY][tileX] == nil {
			continue
		}

		switch sublevel.tileMap[tileY][tileX].Type {
		case "conveyor_left", "machine":
			item.X -= CONVEYOR_SPEED / 16.0
		case "conveyor_down":
			item.Y += CONVEYOR_SPEED / 16.0
		}

		if sublevel.tileMap[tileY][tileX].Type == "machine" {
			if (item.X - float64(tileX)) > 1.0/16 {
				continue
			}

			keepItem := HandleProcessor(g, sublevel, item, sublevel.tileMap[tileY][tileX])
			if keepItem == false {
				sublevel.conveyorItems[i] = nil
			}
		}
	}

	result := []*ConveyorItem{}
	for _, item := range sublevel.conveyorItems {
		if item == nil {
			continue
		}
		result = append(result, item)
	}
	sublevel.conveyorItems = result
}

var machineTypeRecipes map[string]*Recipe

func init() {
	machineTypeRecipes = map[string]*Recipe{
		"seal_board_in_casing": {
			result:      "final_chip",
			ingredients: []string{"casing", "circuit_board_programmed"},
		},
		"program_board": {
			result:      "circuit_board_programmed",
			ingredients: []string{"circuit_board_finished"},
		},
		"add_component_to_finished_board": {
			result:      "circuit_board_finished",
			ingredients: []string{"circuit_board_finished"},
		},
		"combine_copper_and_resin_board": {
			result:      "uncut_circuit_board",
			ingredients: []string{"resin_board", "copper_sheet"},
		},
		"apply_template_to_board": {
			result:      "circuit_board_finished",
			ingredients: []string{"uncut_circuit_board"},
		},
	}
}

// MACHINE TYPES
// seal_board_in_casing
// program_board
// add_component_to_finished_board
// apply_template_to_board

// returns false if the item should be consumed.
func HandleProcessor(g *Game, s *Sublevel, item *ConveyorItem, tile *Tile) bool {
	if tile.Type != "machine" {
		panic("EXPECTED MACHINE")
	}

	recipe, exists := machineTypeRecipes[tile.SubType]
	if !exists {
		return true
	}

	var output string = recipe.result

	if item.itemType.id == output {
		return true
	}

	if item.itemType.id == "reprogramming_chip" || item.itemType.id == "hacking_chip" || item.itemType.id == "final_chip" || item.itemType.id == "broken_chip" {
		return true
	}

	if item.itemType.resultData != "" {
		tile.StoredSubtype = item.itemType.resultData
	}

	if tile.SubType == "seal_board_in_casing" {
		switch tile.StoredSubtype {
		case "Reprogrammer":
			output = "reprogramming_chip"
		case "Hacking USB":
			output = "hacking_chip"
		case "Mind Control":
			output = "final_chip"
		default:
			output = "broken_chip"
		}
	}

	if tile.Inventory == nil {
		tile.Inventory = map[string]int{}
	}
	tile.Inventory[item.itemType.id] += 1

	canCraft := true
	for _, ingredient := range recipe.ingredients {
		if tile.Inventory[ingredient] == 0 {
			canCraft = false
		}
	}

	if !canCraft {
		return false
	}

	for _, ingredient := range recipe.ingredients {
		tile.Inventory[ingredient] -= 1
	}

	item.itemType = &Item{
		id:         output,
		resultData: tile.StoredSubtype,
	}

	if tile.SubType == "apply_template_to_board" {
		if g.templateItem != nil {
			item.itemType.resultData = g.templateItem.resultData
		} else {
			item.itemType.resultData = "Empty"
		}
	}
	return true
}

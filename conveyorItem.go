package main

type ConveyorItem struct {
	// Position aligned to grid
	X, Y     float64
	itemType *Item
}

func UpdateConveyors(g *Game) {
	sublevel := g.CurrentSublevel()
	if g.t%10 != 0 {
		return
	}

	for _, item := range sublevel.conveyorItems {
		var nextSpaceOccupied = false
		tileX, tileY := int(item.X), int(item.Y)
		for _, otherItem := range sublevel.conveyorItems {
			if item == otherItem {
				continue
			}
			otherTileX, otherTileY := int(item.X), int(item.Y)
			if tileX == (otherTileX-1) && tileY == otherTileY {
				nextSpaceOccupied = true
			}
		}
		if nextSpaceOccupied {
			continue
		}

		if tileX < 0 || sublevel.tileMap[tileY][tileX] == nil {
			continue
		}

		switch sublevel.tileMap[tileY][tileX].Type {
		case "conveyor_left", "machine":
			item.X -= CONVEYOR_SPEED / 16.0
		case "conveyor_down":
			item.Y += CONVEYOR_SPEED / 16.0
		}

	}
}

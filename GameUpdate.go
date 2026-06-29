package main

import (
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func StartLockdown(g *Game) {
	if g.isLockDown {
		return
	}
	g.isElectricityDown = false
	g.isLockDown = true

	PlaySound(
		g.audioContext,
		lockdownSound,
		2,
	)

	g.CurrentSublevel().Enemies = append(g.CurrentSublevel().Enemies, createEnemy(), createEnemy())
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		CompleteTask(g, g.testIndex)
		g.testIndex += 1
	}
	g.uiContext.PreUpdate()

	if g.inEndScreen {
		g.endScreenUI.Update()
		g.uiContext.Update()
		return nil
	}

	if g.inCraftingUi {
		g.craftingUI.Update()
	}
	if g.inTodoUI {
		g.todoUI.Update()
	}
	if g.dayEndedInDeath {
		g.lossScreen.Update()
	}
	if g.inIntroScreen {
		g.introScreen.Update()
	}
	if g.inMainMenu {
		g.mainMenuUi.Update()
	}
	if g.inWinScreen {
		g.winScreen.Update()
	}
	g.hudUI.Update()
	g.uiContext.Update()

	// Update called each tick (1/60 s)
	g.t += 1
	if g.timeRemaining > 0 {
		g.timeRemaining -= 1
	}

	if g.isElectricityDown && g.electricityDownRemainingTime > 0 {
		g.electricityDownRemainingTime -= 1
	}

	if g.isElectricityDown && g.electricityDownRemainingTime == 0 {
		StartLockdown(g)
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyB) {
	// 	StartLockdown(g)
	// }

	if g.isLockDown && g.t%60 == 0 {
		PlaySound(
			g.audioContext,
			lockdownSound,
			2,
		)
	}
	if g.isLockDown && g.t%60 == 0 {
		if g.CurrentSublevel().isSafeArea && len(g.CurrentSublevel().Enemies) < 2 {
			g.CurrentSublevel().Enemies = append(g.CurrentSublevel().Enemies, createEnemy())
		} else if len(g.CurrentSublevel().Enemies) < 4 {
			g.CurrentSublevel().Enemies = append(g.CurrentSublevel().Enemies, createEnemy())
		}
	} else if !g.isElectricityDown && !g.CurrentSublevel().isSafeArea && len(g.CurrentSublevel().Enemies) < 1 && g.t%(60*60*2) == 0 {
		g.CurrentSublevel().Enemies = append(g.CurrentSublevel().Enemies, createEnemy())
	}

	for _, space := range g.spaces {
		if space.Spawners != nil {
			for _, spawner := range space.Spawners {
				if g.t%(CONVEYOR_SPEED*16*4) != 0 {
					continue
				}

				space.conveyorItems = append(space.conveyorItems,
					&ConveyorItem{
						X: float64(spawner.X),
						Y: float64(spawner.Y),
						itemType: &Item{
							id: spawner.item,
						},
					},
				)
			}
		}

		UpdateConveyors(g, space)
	}

	if g.inExclusiveUIMode {
		return nil
	}

	// Craft
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.inCraftingUi = !g.inCraftingUi
		if g.inCraftingUi {
			g.inTodoUI = false
			openCraftingUI(g)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.inTodoUI = !g.inTodoUI
		if g.inTodoUI {
			g.inCraftingUi = false
		}
	}

	cursorX, cursorY := ebiten.CursorPosition()
	targetX, targetY := cursorX/16, cursorY/16
	gridCursorX, gridCursorY := float64(cursorX)/16, float64(cursorY)/16

	playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	cursorSelectionTile := GetTileFromTileMap(g.CurrentSublevel().tileMap, targetX, targetY)
	tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

	if tileImOn != nil && tileImOn.Type == "conveyor_left" && g.t%10 == 0 {
		g.player.position.X -= CONVEYOR_SPEED / 16.0
	}

	if tileImOn != nil && tileImOn.Type == "conveyor_down" && g.t%10 == 0 {
		g.player.position.Y += CONVEYOR_SPEED / 16.0
	}

	if g.CurrentSublevel().HasFlowingWater && g.t%4 == 0 {
		xSpeed := 0.0
		if g.player.position.Y >= 3 && g.player.position.Y <= 5 {
			xSpeed = +1.0 / 16
		}

		if g.player.position.Y >= 8 && g.player.position.Y <= 10 {
			xSpeed = -1.0 / 16
		}
		xOffset := 0.0
		if xSpeed > 0 {
			xOffset = 4.0 / 16
		} else if xSpeed < 0 {
			xOffset = -4.0 / 16
		}

		playerCenter := VectorAdd(g.player.position, Vector{.5, .5})
		adjacentVectorX := VectorFloor(VectorAdd(playerCenter, Vector{X: xSpeed + xOffset, Y: 0}))
		if adjacentVectorX.X >= 0 && adjacentVectorX.X <= 19 {
			tileIWillBeOn := g.CurrentSublevel().tileMap[int(adjacentVectorX.Y)][int(adjacentVectorX.X)]
			if !TileIsSolid(g, tileIWillBeOn) {
				g.player.position.X += xSpeed
			}
		}
	}

	g.handlePlayerMovement()

	// g.setTileToWall()

	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.selectedSlot = 0
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.selectedSlot = 1
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.selectedSlot = 2
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.selectedSlot = 3
	} else if ebiten.IsKeyPressed(ebiten.Key5) {
		g.selectedSlot = 4
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
	// 	fmt.Println(int(gridCursorX), ", ", int(gridCursorY))
	// }

	heldItem := g.inventory[g.selectedSlot]

	clickUseType := ""

	if cursorSelectionTile == nil && heldItem != nil && heldItem.id == "box" {
		clickUseType = "place_box"
	}

	var selectedConveyorItem *ConveyorItem
	var selectedConveyorItemIndex int
	for index, conveyorItem := range g.CurrentSublevel().conveyorItems {
		itemCenterX, itemCenterY := conveyorItem.X+0.5, conveyorItem.Y+0.5

		distance := math.Sqrt(math.Pow(itemCenterX-gridCursorX, 2) + math.Pow(itemCenterY-gridCursorY, 2))

		if distance < 0.5 {
			clickUseType = "conveyor_item"
			selectedConveyorItem = conveyorItem
			selectedConveyorItemIndex = index
			break
		}
	}

	var selectedItem *InGameItem
	var selectedItemIndex int
	for index, item := range g.CurrentSublevel().inGameItems {

		itemCenterX, itemCenterY := item.X+0.5, item.Y+0.5

		distance := math.Sqrt(math.Pow(itemCenterX-gridCursorX, 2) + math.Pow(itemCenterY-gridCursorY, 2))

		if distance < 0.5 {
			clickUseType = "item"
			selectedItem = item
			selectedItemIndex = index
			break
		}
	}

	if cursorSelectionTile != nil && cursorSelectionTile.Type == "box" {
		clickUseType = "box"
	}

	if cursorSelectionTile != nil && cursorSelectionTile.Type == "wire" {
		clickUseType = "wire"
	}

	var selectedEnemy *Enemy
	var selectedEnemyIndex int
	for index, enemy := range g.CurrentSublevel().Enemies {
		enemyCenterX, enemyCenterY := enemy.position.X+0.5, enemy.position.Y+0.5

		distance := math.Sqrt(math.Pow(enemyCenterX-gridCursorX, 2) + math.Pow(enemyCenterY-gridCursorY, 2))

		if distance < 1.25 {
			clickUseType = "enemy"
			selectedEnemy = enemy
			selectedEnemyIndex = index
			break
		}
	}

	if cursorSelectionTile != nil && cursorSelectionTile.Type == "vent_down" {
		clickUseType = "vent_down"
	}
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "vent_down_open" {
		clickUseType = "vent_down_open"
	}
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "electrical_switch" {
		clickUseType = "electrical_switch"
	}
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "escape_ladder" {
		clickUseType = "escape_ladder"
	}
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "central_console" {
		clickUseType = "central_console"
	}
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "template_holder" {
		clickUseType = "template_holder"
	}
	if tileImOn != nil && tileImOn.Type == "box" && g.CurrentSublevel().adjacentSpaces.Up != "" {
		aboveSpaceId := g.CurrentSublevel().adjacentSpaces.Up
		aboveSpace := g.spaces[aboveSpaceId]

		tileAbove := GetTileFromTileMap(aboveSpace.tileMap, targetX, targetY)
		if tileAbove != nil {
			switch tileAbove.Type {
			case "vent_down":
				clickUseType = "vent_up"
			case "vent_down_open":
				clickUseType = "vent_up_open"
			}
		}
	}

	// Selection name
	switch clickUseType {
	case "conveyor_item":
		g.selectionName = itemData[selectedConveyorItem.itemType.id].name
		if selectedConveyorItem.itemType.resultData != "" {
			g.selectionName += " (" + selectedConveyorItem.itemType.resultData + ")"
		}
	case "item":
		g.selectionName = itemData[selectedItem.itemType.id].name
	case "vent_up", "vent_down":
		g.selectionName = "Vent"
	case "vent_up_open", "vent_down_open":
		g.selectionName = "Vent (Unscrewed)"
	case "box":
		g.selectionName = "Box"
	case "wire":
		g.selectionName = "Wire Mesh"
	case "electrical_switch":
		g.selectionName = "Electrical Switch"
	case "central_console":
		g.selectionName = "Central Console (Insert USB)"
	case "escape_ladder":
		g.selectionName = "Escape Ladder"
	case "enemy":
		g.selectionName = "Robot"
	case "place_box":
		g.selectionName = "Place Box"
	case "template_holder":
		if g.templateItem == nil {
			g.selectionName = "Insert Template"
		} else {
			g.selectionName = "Remove Template"
		}
	default:
		g.selectionName = ""
	}

	// Selection Box Position
	switch clickUseType {
	case "conveyor_item":
		g.targetX, g.targetY = selectedConveyorItem.X, selectedConveyorItem.Y
	case "item":
		g.targetX, g.targetY = selectedItem.X, selectedItem.Y
	case "enemy":
		g.targetX, g.targetY = math.Floor(selectedEnemy.position.X*16)/16, math.Floor(selectedEnemy.position.Y*16)/16
	case "electrical_switch":
		g.targetX, g.targetY = float64(cursorX/16), float64(cursorY/16)-0.25
	default:
		g.targetX, g.targetY = float64(cursorX/16), float64(cursorY/16)
	}

	distanceBetweenPlayerAndTarget := math.Sqrt(math.Pow(g.targetX-g.player.position.X, 2) + math.Pow(g.targetY-g.player.position.Y, 2))

	switch clickUseType {
	case "vent_down", "vent_down_open", "vent_up", "vent_up_open":
		g.selectionInRange = distanceBetweenPlayerAndTarget < 0.8
	case "escape_ladder":
		g.selectionInRange = distanceBetweenPlayerAndTarget < 2
	default:
		g.selectionInRange = distanceBetweenPlayerAndTarget < 3
	}

	switch clickUseType {
	case "wire":
		if heldItem == nil || heldItem.id != "wire_cutters" {
			g.selectionInRange = false
		}
	case "vent_up", "vent_down":
		if heldItem == nil || heldItem.id != "screwdriver" {
			g.selectionInRange = false
		}
	case "enemy":
		if !g.HasItem("hammer") {
			g.selectionInRange = false
		}
		if heldItem == nil || heldItem.id != "hammer" {
		}
	case "central_console":
		if heldItem == nil || heldItem.id != "hacking_usb" {
			g.selectionInRange = false
		}
	case "electrical_switch":
		if g.isElectricityDown || g.isLockDown {
			g.selectionInRange = false
		}
	}

	hasSelection := g.selectionName != ""

	if hasSelection && g.selectionInRange && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		switch clickUseType {
		case "conveyor_item":
			replaceSlot := GetOrFreeSlotForItemInHotbar(g)
			g.inventory[replaceSlot] = selectedConveyorItem.itemType
			g.CurrentSublevel().conveyorItems = slices.Delete(g.CurrentSublevel().conveyorItems, selectedConveyorItemIndex, selectedConveyorItemIndex+1)

			CompleteTask(g, 3)

			switch selectedConveyorItem.itemType.id {
			case "hacking_chip":
				CompleteTask(g, 14)
			}
			if !g.isElectricityDown {
				StartLockdown(g)
			}
		case "item":
			replaceSlot := GetOrFreeSlotForItemInHotbar(g)
			g.inventory[replaceSlot] = selectedItem.itemType
			g.CurrentSublevel().inGameItems = slices.Delete(g.CurrentSublevel().inGameItems, selectedItemIndex, selectedItemIndex+1)

			switch selectedItem.itemType.id {
			case "auth_chip":
				CompleteTask(g, 9)
			}
		case "vent_down_open":
			if g.CurrentSublevel().adjacentSpaces.Down != "" {
				g.currentSpaceId = g.CurrentSublevel().adjacentSpaces.Down
				g.prevLevelDirection = "Down"
				g.prevLevelMoveTime = g.t
			}
		case "vent_down":
			g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
				Type: "vent_down_open",
			}
		case "vent_up":
			aboveSpaceId := g.CurrentSublevel().adjacentSpaces.Up
			aboveSpace := g.spaces[aboveSpaceId]

			aboveSpace.tileMap[targetY][targetX] = &Tile{
				Type: "vent_down_open",
			}

			if g.currentSpaceId == "sewer_entrance" {
				CompleteTask(g, 2)
			}
		case "vent_up_open":
			switch g.currentSpaceId {
			case "sewer_0":
				CompleteTask(g, 8)
			case "sewer_4":
				CompleteTask(g, 17)
			}
			g.currentSpaceId = g.CurrentSublevel().adjacentSpaces.Up
			g.prevLevelDirection = "Up"
			g.prevLevelMoveTime = g.t
		case "box":
			cursorSelectionTile.Damage += 0.2

			if cursorSelectionTile.Damage >= 1 {
				replaceSlot := GetOrFreeSlotForItemInHotbar(g)
				g.inventory[replaceSlot] = &Item{
					id: "box",
				}
				g.CurrentSublevel().tileMap[targetY][targetX] = nil
			}
		case "wire":
			cursorSelectionTile.Damage += 0.2

			if cursorSelectionTile.Damage >= 1 {
				g.CurrentSublevel().tileMap[targetY][targetX] = nil
				CompleteTask(g, 7)
			}
		case "enemy":
			kbScale := 0.05
			if g.HasItem("hammer") {
				selectedEnemy.Health -= 3
				kbScale = 0.24

				PlaySound(
					g.audioContext,
					RandomSound(impactPlankSound),
					1,
				)
			}

			enemyToPlayerVector := VectorNormalise(VectorSubtract(selectedEnemy.position, g.player.position))
			knockbackVector := VectorScale(enemyToPlayerVector, kbScale)
			selectedEnemy.Acceleration = knockbackVector

			// kill enemy
			if selectedEnemy.Health <= 0 {
				g.CurrentSublevel().Enemies = slices.Delete(g.CurrentSublevel().Enemies, selectedEnemyIndex, selectedEnemyIndex+1)

				hasAuth := false
				for _, item := range g.inventory {
					if item == nil {
						continue
					}
					if item.id == "auth_card" {
						hasAuth = true
						break
					}
				}
				if !hasAuth {
					g.CurrentSublevel().inGameItems = append(g.CurrentSublevel().inGameItems,
						&InGameItem{
							itemType: &Item{
								id: "auth_chip",
							},
							X: selectedEnemy.position.X,
							Y: selectedEnemy.position.Y,
						},
					)
				}

				PlaySound(
					g.audioContext,
					RandomSound(robotDeathSound),
					1,
				)

			}
		case "escape_ladder":
			if g.tasks[16].complete && g.tasks[18].complete {
				g.inWinScreen = true
				g.inExclusiveUIMode = true
				g.isLockDown = false
			} else {
				g.EndDay()
				CompleteTask(g, 4)
			}
		case "place_box":
			g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
				Type:   "box",
				Damage: 0,
			}

			if g.currentSpaceId == "sewer_entrance" {
				CompleteTask(g, 1)
			}

			g.inventory[g.selectedSlot] = nil
		case "template_holder":
			if g.templateItem == nil && g.inventory[g.selectedSlot] != nil && g.inventory[g.selectedSlot].id == "template" {
				g.templateItem = g.inventory[g.selectedSlot]
				g.inventory[g.selectedSlot] = nil

				switch g.templateItem.resultData {
				case "Reprogrammer":
					CompleteTask(g, 16)
				case "Hacking USB":
					CompleteTask(g, 13)
				case "Mind Control":
				}

			} else {
				replaceSlot := GetOrFreeSlotForItemInHotbar(g)
				g.inventory[replaceSlot] = g.templateItem
				g.templateItem = nil
			}

			if !g.isElectricityDown {
				StartLockdown(g)
			}
		case "central_console":
			CompleteTask(g, 18)
			g.inventory[g.selectedSlot] = nil
			StartLockdown(g)
		case "electrical_switch":
			g.isElectricityDown = true
			g.electricityDownRemainingTime = 60 * 61
			CompleteTask(g, 11)
		}
	}

	switch clickUseType {
	case "box", "wire":
		g.progressBar = 1 - cursorSelectionTile.Damage
	case "enemy":
		g.progressBar = float64(selectedEnemy.Health) / float64(selectedEnemy.MaxHealth)
	default:
		g.progressBar = -1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		dropItemSlot(g, g.selectedSlot)
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
	// 	PrintSublevel(g)
	// }

	// UpdateConveyors(g, g.CurrentSublevel())

	for _, e := range g.CurrentSublevel().Enemies {
		dx, dy := (e.position.X - g.player.position.X), (e.position.Y - g.player.position.Y)
		m := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

		e.Move(g.CurrentSublevel(), (g.player.position.X), (g.player.position.Y))
		isAttacking := e.AttackStartT+e.AttackLength > g.t
		if isAttacking {
			if m < 2 && g.t%10 == 0 {
				g.Health -= 1

				if g.Health <= 0 {
					g.dayEndedInDeath = true
					g.isLockDown = false
					g.inExclusiveUIMode = true
				}
			}
		} else if m < 2.5 {
			if g.t%30 == 0 {
				e.AttackStartT = g.t
				e.AttackLength = 31

				PlaySound(
					g.audioContext,
					RandomSound(shootSound),
					0.5,
				)
			}
		}
	}

	return nil
}

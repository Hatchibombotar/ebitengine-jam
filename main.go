package main

import (
	"fmt"
	"hatchi/disconnect/superui"
	"image"
	"image/color"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const CONVEYOR_SPEED = 2.0

type Game struct {
	t      int
	player *Character

	uiContext  *superui.UIContext
	hudUI      *superui.UIContainer
	craftingUI *superui.UIContainer
	todoUI     *superui.UIContainer

	inventory    [8]*Item
	selectedSlot int

	inCraftingUi bool
	inTodoUI     bool
	inEndScreen  bool

	selectedRecipe int

	sublevel string

	sublevels map[string]*Sublevel

	// Relative to grid
	targetX, targetY float64

	selectionName    string
	selectionInRange bool

	day int

	progressBar float64

	timeRemaining int

	endScreenUI *superui.UIContainer

	tasks []*Task

	Health    int
	MaxHealth int

	templateItem *Item
}

func (g *Game) CurrentSublevel() *Sublevel {
	sublevel, ok := g.sublevels[g.sublevel]
	if !ok {
		panic(fmt.Sprint("sublevel doesn't exist.", g.sublevel))
	}
	return sublevel
}

const ITEM_REACH_RANGE = 1.0

func (g *Game) Update() error {
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
	g.hudUI.Update()
	g.uiContext.Update()

	// Update called each tick (1/60 s)
	g.t += 1
	if g.timeRemaining > 0 {
		g.timeRemaining -= 1
	}

	if g.CurrentSublevel().Spawners != nil {
		for _, spawner := range g.CurrentSublevel().Spawners {
			if g.t%(CONVEYOR_SPEED*16*4) != 0 {
				continue
			}

			g.CurrentSublevel().conveyorItems = append(g.CurrentSublevel().conveyorItems,
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

	g.handlePlayerMovement()

	g.setTileToWall()

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

	cursorX, cursorY := ebiten.CursorPosition()
	targetX, targetY := cursorX/16, cursorY/16
	gridCursorX, gridCursorY := float64(cursorX)/16, float64(cursorY)/16

	playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	cursorSelectionTile := GetTileFromTileMap(g.CurrentSublevel().tileMap, targetX, targetY)
	tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

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
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "template_holder" {
		clickUseType = "template_holder"
	}
	if tileImOn != nil && tileImOn.Type == "box" && g.CurrentSublevel().adjacentSpaces.Up != "" {
		aboveSpaceId := g.CurrentSublevel().adjacentSpaces.Up
		aboveSpace := g.sublevels[aboveSpaceId]

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
		g.selectionInRange = distanceBetweenPlayerAndTarget < 1
	case "escape_ladder":
		g.selectionInRange = distanceBetweenPlayerAndTarget < 2
	default:
		g.selectionInRange = distanceBetweenPlayerAndTarget < 3
	}

	hasSelection := g.selectionName != ""

	if hasSelection && g.selectionInRange && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		switch clickUseType {
		case "conveyor_item":
			replaceSlot := GetOrFreeSlotForItemInHotbar(g)
			g.inventory[replaceSlot] = selectedConveyorItem.itemType
			g.CurrentSublevel().conveyorItems = slices.Delete(g.CurrentSublevel().conveyorItems, selectedConveyorItemIndex, selectedConveyorItemIndex+1)
		case "item":
			replaceSlot := GetOrFreeSlotForItemInHotbar(g)
			g.inventory[replaceSlot] = selectedItem.itemType
			g.CurrentSublevel().inGameItems = slices.Delete(g.CurrentSublevel().inGameItems, selectedItemIndex, selectedItemIndex+1)
		case "vent_down_open":
			if g.CurrentSublevel().adjacentSpaces.Down != "" {
				g.sublevel = g.CurrentSublevel().adjacentSpaces.Down
			}
		case "vent_down":
			g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
				Type: "vent_down_open",
			}
		case "vent_up":
			aboveSpaceId := g.CurrentSublevel().adjacentSpaces.Up
			aboveSpace := g.sublevels[aboveSpaceId]

			aboveSpace.tileMap[targetY][targetX] = &Tile{
				Type: "vent_down_open",
			}
		case "vent_up_open":
			g.sublevel = g.CurrentSublevel().adjacentSpaces.Up
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
			}
		case "enemy":
			selectedEnemy.Health -= 2
			enemyToPlayerVector := VectorNormalise(VectorSubtract(selectedEnemy.position, g.player.position))
			knockbackVector := VectorScale(enemyToPlayerVector, 0.15)
			selectedEnemy.Acceleration = knockbackVector

			// kill enemy
			if selectedEnemy.Health <= 0 {
				g.CurrentSublevel().Enemies = slices.Delete(g.CurrentSublevel().Enemies, selectedEnemyIndex, selectedEnemyIndex+1)
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
		case "escape_ladder":
			g.EndDay()
		case "place_box":
			g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
				Type:   "box",
				Damage: 0,
			}

			g.inventory[g.selectedSlot] = nil
		case "template_holder":
			if g.templateItem == nil && g.inventory[g.selectedSlot] != nil && g.inventory[g.selectedSlot].id == "template" {
				g.templateItem = g.inventory[g.selectedSlot]
				g.inventory[g.selectedSlot] = nil
			} else {
				replaceSlot := GetOrFreeSlotForItemInHotbar(g)
				g.inventory[replaceSlot] = g.templateItem
				g.templateItem = nil
			}
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

	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		PrintSublevel(g)
	}

	UpdateConveyors(g)

	for _, e := range g.CurrentSublevel().Enemies {
		e.Move(g.CurrentSublevel(), (g.player.position.X), (g.player.position.Y))
	}

	return nil
}

func (g *Game) handlePlayerMovement() {
	g.player.startLerpT = 1

	playerCenter := VectorAdd(g.player.position, Vector{.5, .5})

	xSpeed, ySpeed := 0.0, 0.0

	speed := 0.1

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		xSpeed += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		xSpeed -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ySpeed -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ySpeed += speed
	}

	xOffset := 0.0
	if xSpeed > 0 {
		xOffset = 4.0 / 16
	} else if xSpeed < 0 {
		xOffset = -4.0 / 16
	}

	yOffset := 0.0
	if ySpeed > 0 {
		yOffset = 0.5
	} else if ySpeed < 0 {
		yOffset = -5.0 / 16
	}

	adjacentVectorX := VectorFloor(VectorAdd(playerCenter, Vector{X: xSpeed + xOffset, Y: 0}))
	if adjacentVectorX.X < 0 {
		if g.CurrentSublevel().adjacentSpaces.West != "" {
			g.player.position.X = 19
			g.sublevel = g.CurrentSublevel().adjacentSpaces.West
		}
		return
	} else if adjacentVectorX.X > 19 {
		if g.CurrentSublevel().adjacentSpaces.East != "" {
			g.player.position.X = 0
			g.sublevel = g.CurrentSublevel().adjacentSpaces.East
		}
		return
	}

	adjacentVectorY := VectorFloor(VectorAdd(playerCenter, Vector{X: 0, Y: ySpeed + yOffset}))
	if adjacentVectorY.Y < 0 {
		if g.CurrentSublevel().adjacentSpaces.North != "" {
			g.player.position.Y = 14
			g.sublevel = g.CurrentSublevel().adjacentSpaces.North
		}
		return
	} else if adjacentVectorY.Y > 14 {
		if g.CurrentSublevel().adjacentSpaces.South != "" {
			g.player.position.Y = 0
			g.sublevel = g.CurrentSublevel().adjacentSpaces.South
		}
		return
	}

	if TileIsSolid(g.CurrentSublevel().tileMap[int(adjacentVectorX.Y)][int(adjacentVectorX.X)]) {
		xSpeed = 0
	}
	if TileIsSolid(g.CurrentSublevel().tileMap[int(adjacentVectorY.Y)][int(adjacentVectorY.X)]) {
		ySpeed = 0
	}

	if xSpeed != 0 || ySpeed != 0 {
		g.player.startLerpT = g.t
	}

	if xSpeed != 0 && ySpeed != 0 {
		xSpeed *= 0.7071
		ySpeed *= 0.7071
	}

	if xSpeed > 0 {
		g.player.facingDirection = Vector{X: 1, Y: 0}
	} else if xSpeed < 0 {
		g.player.facingDirection = Vector{X: -1, Y: 0}
	} else if ySpeed > 0 {
		g.player.facingDirection = Vector{X: 0, Y: 1}
	} else if ySpeed < 0 {
		g.player.facingDirection = Vector{X: 0, Y: -1}
	}

	g.player.position.X += xSpeed
	g.player.position.Y += ySpeed
}

func (g *Game) setTileToWall() {
	cursorX, cursorY := ebiten.CursorPosition()
	targetX, targetY := (cursorX / 16), (cursorY / 16)

	if ebiten.IsKeyPressed(ebiten.KeyL) {
		g.CurrentSublevel().tileMap[targetY][targetX] = nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "wall",
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key0) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "conveyor_left",
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key9) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "conveyor_down",
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key8) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "machine",
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "vent_down",
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "wire",
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.CurrentSublevel().conveyorItems = append(g.CurrentSublevel().conveyorItems,
			&ConveyorItem{
				X: float64(cursorX / 16),
				Y: float64(cursorY / 16),
				itemType: &Item{
					id: "copper_sheet",
				},
			},
		)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.CurrentSublevel().conveyorItems = append(g.CurrentSublevel().conveyorItems,
			&ConveyorItem{
				X: float64(cursorX / 16),
				Y: float64(cursorY / 16),
				itemType: &Item{
					id: "resin_board",
				},
			},
		)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.inEndScreen {
		screen.Fill(color.Black)
		g.endScreenUI.Draw(screen)
		return
	}

	screen.DrawImage(
		g.CurrentSublevel().Background, nil,
	)

	for y, row := range g.CurrentSublevel().tileMap {
		for x, tile := range row {
			if tile == nil {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))

			switch tile.Type {
			case "box":
				screen.DrawImage(item_box, op)
			case "conveyor_left":
				// screen.DrawImage(conveyor_left, op)
				t := CONVEYOR_SPEED * (g.t / 10)
				screen.DrawImage(
					conveyor_left_flipbook.SubImage(image.Rect(0, (t%4)*16, 16, 16+((t%4)*16))).(*ebiten.Image),
					op,
				)
			case "conveyor_down":
				// screen.DrawImage(conveyor_down, op)
				t := CONVEYOR_SPEED * (g.t / 10)
				screen.DrawImage(
					conveyor_down_flipbook.SubImage(image.Rect(0, (t%4)*16, 16, 16+((t%4)*16))).(*ebiten.Image),
					op,
				)
			case "machine":
				op.GeoM.Translate(0, -6)
				screen.DrawImage(machine, op)
			case "vent_down":
				screen.DrawImage(vent, op)
			case "vent_down_open":
				screen.DrawImage(vent_open, op)
			case "wire":
				screen.DrawImage(wire, op)
			case "template_holder":
				if g.templateItem != nil {
					op.GeoM.Translate(0, -3)
					screen.DrawImage(item_template, op)
				}
			}
		}
	}

	playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

	playerIsAccended := TileIsAccended(tileImOn)

	if tileImOn != nil && tileImOn.Type == "conveyor_left" && g.t%10 == 0 {
		g.player.position.X -= CONVEYOR_SPEED / 16.0
	}

	if tileImOn != nil && tileImOn.Type == "conveyor_down" && g.t%10 == 0 {
		g.player.position.Y += CONVEYOR_SPEED / 16.0
	}

	if playerIsAccended {
		g.player.Draw(screen, g, 0, -6, false)
	} else {
		g.player.Draw(screen, g, 0, 0, false)
	}

	for _, inGameItem := range g.CurrentSublevel().inGameItems {
		itemData := itemData[inGameItem.itemType.id]

		playerDistance := VectorMagnitude(
			VectorSubtract(
				g.player.position,
				Vector{X: float64(inGameItem.X), Y: float64(inGameItem.Y)},
			),
		)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(inGameItem.X)*16, float64(inGameItem.Y)*16)

		if playerDistance < 1.1 {
			op.GeoM.Translate(0, math.Sin(float64(g.t)*0.5)*1)
		}

		screen.DrawImage(itemData.image, op)
	}

	for _, inGameItem := range g.CurrentSublevel().conveyorItems {
		itemData := itemData[inGameItem.itemType.id]
		if itemData == nil {
			continue
		}

		playerDistance := VectorMagnitude(
			VectorSubtract(
				g.player.position,
				Vector{X: float64(inGameItem.X), Y: float64(inGameItem.Y)},
			),
		)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(inGameItem.X)*16, float64(inGameItem.Y)*16)

		if playerDistance < 1.1 {
			op.GeoM.Translate(0, math.Sin(float64(g.t)*0.5)*1)
		}

		screen.DrawImage(itemData.image, op)
	}

	for y, row := range g.CurrentSublevel().tileMap {
		for x, tile := range row {
			if tile == nil {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))

			switch tile.Type {
			case "machine":
				op.GeoM.Translate(0, -6)
				screen.DrawImage(machine, op)
			}
		}
	}

	if g.CurrentSublevel().Overlay != nil {
		screen.DrawImage(g.CurrentSublevel().Overlay, nil)
	}

	// cursorX, cursorY := ebiten.CursorPosition()

	// debug draw
	if ebiten.IsKeyPressed(ebiten.KeyF8) {
		for y, row := range g.CurrentSublevel().tileMap {
			for x, tile := range row {
				if tile == nil {
					continue
				}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*16), float64(y*16))

				switch tile.Type {
				case "wall":
					screen.DrawImage(debug_wall, op)
				}
			}
		}
	}

	if g.CurrentSublevel().adjacentSpaces.Up != "" {
		aboveSpaceId := g.CurrentSublevel().adjacentSpaces.Up
		aboveSpace := g.sublevels[aboveSpaceId]

		standingOnBox := tileImOn != nil && tileImOn.Type == "box"

		for y, row := range aboveSpace.tileMap {
			for x, tile := range row {
				if tile == nil {
					continue
				}
				currentTile := g.CurrentSublevel().tileMap[y][x]
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*16), float64(y*16))
				op.ColorScale.ScaleAlpha(0.65)

				if standingOnBox {
					switch tile.Type {
					case "vent_down":
						screen.DrawImage(vent, op)
					case "vent_down_open":
						screen.DrawImage(vent_open, op)
					}
				} else {
					op.ColorScale.ScaleAlpha(0.6)
					if currentTile != nil && currentTile.Type == "box" {
						op.ColorScale.ScaleAlpha(0.3)
					}
					switch tile.Type {
					case "vent_down":
						screen.DrawImage(vent_shadow, op)
					case "vent_down_open":
						screen.DrawImage(vent_open_shadow, op)
					}
				}
			}
		}
	}

	for _, e := range g.CurrentSublevel().Enemies {
		e.Draw(screen, g, 0, 0, false)
	}

	ops := &ebiten.DrawRectShaderOptions{}
	ops.Images[0] = screen

	shaded := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	shaded.DrawRectShader(320, 240, crtshader, ops)
	screen.DrawImage(shaded, nil)

	if !g.uiContext.IsHovered() {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate((g.targetX)*16, (g.targetY)*16)

		if g.selectionName != "" {
			if g.selectionInRange {
				screen.DrawImage(target_green, op)
			} else {
				screen.DrawImage(target_red, op)
			}
		} else {
			// screen.DrawImage(target, op)
		}
	}

	screen.DrawImage(g.CurrentSublevel().Vignette, nil)

	if g.inCraftingUi {
		g.craftingUI.Draw(screen)
	}
	if g.inTodoUI {
		g.todoUI.Draw(screen)
	}

	g.hudUI.Draw(screen)

	op := &text.DrawOptions{}
	op.GeoM.Translate(320-4, 0)
	op.PrimaryAlign = text.AlignEnd

	text.Draw(screen, fmt.Sprint("Day ", g.day), &text.GoTextFace{
		Source: fontFaceSource,
		Size:   16,
	}, op)

	smallFontFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   8,
	}

	op = &text.DrawOptions{}
	op.GeoM.Translate(320-4, 2+16)
	op.PrimaryAlign = text.AlignEnd
	op.ColorScale.ScaleWithColor(color.Gray{180})

	text.Draw(screen, g.CurrentSublevel().Title, smallFontFace, op)

	totalSeconds := g.timeRemaining / 60
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	formattedTimeRemaining := fmt.Sprintf("%02d:%02d", minutes, seconds)

	timeWidth, _ := text.Measure(formattedTimeRemaining, smallFontFace, 1)
	op = &text.DrawOptions{}
	op.GeoM.Translate(320-4-timeWidth, 2+16+16)
	op.PrimaryAlign = text.AlignEnd
	op.ColorScale.ScaleWithColor(color.Gray{180})

	text.Draw(screen, "Time Remaining: ", smallFontFace, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(320-4, 2+16+16)
	op.PrimaryAlign = text.AlignEnd
	op.ColorScale.ScaleWithColor(color.Gray{180})

	text.Draw(screen, formattedTimeRemaining, smallFontFace, op)

	if g.selectionName != "" {
		textX, textY := g.targetX*16+2, (g.targetY*16)-8-4
		textWidth, textHeight := text.Measure(g.selectionName, smallFontFace, 1)

		op = &text.DrawOptions{}
		op.GeoM.Translate(textX, textY)

		if textX+textWidth > 320 {
			op.PrimaryAlign = text.AlignEnd
			textX -= textWidth
		}

		vector.FillRect(screen, float32(textX)-2, float32(textY)-1, float32(textWidth)+4, float32(textHeight)+2, color.RGBA{32, 32, 32, 255}, true)

		text.Draw(screen, g.selectionName, smallFontFace, op)
	}

	if g.progressBar != -1 {
		textX, textY := g.targetX*16+2, (g.targetY*16)+16+4
		width, height := 20, 3
		vector.FillRect(screen, float32(textX)-2, float32(textY)-2, float32(width)+4, float32(height)+4, color.RGBA{32, 32, 32, 255}, true)

		colour := color.RGBA{110, 206, 84, 255} // green
		if g.progressBar < 0.33 {
			colour = color.RGBA{206, 89, 84, 255} // red
		} else if g.progressBar < 0.66 {
			colour = color.RGBA{206, 149, 84, 255} // orange
		}
		vector.FillRect(screen, float32(textX), float32(textY), float32(width)*float32(g.progressBar), float32(height), colour, true)
	}

	healthBarX, healthBarY := 8, 200
	width, height := 65, 4
	percentHealth := float64(g.Health) / float64(g.MaxHealth)
	vector.FillRect(screen, float32(healthBarX)-2, float32(healthBarY)-2, float32(width)+4, float32(height)+4, color.RGBA{32, 32, 32, 255}, true)

	colour := color.RGBA{110, 206, 84, 255} // green
	if percentHealth < 0.33 {
		colour = color.RGBA{206, 89, 84, 255} // red
	} else if percentHealth < 0.66 {
		colour = color.RGBA{206, 149, 84, 255} // orange
	}
	vector.FillRect(screen, float32(healthBarX), float32(healthBarY), float32(width)*float32(percentHealth), float32(height), colour, true)

	// ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.ActualFPS()))
}

func (g *Game) StartDay() {
	g.inEndScreen = false
	g.timeRemaining = 3 * 60 * 60
	g.day += 1
	g.Health = g.MaxHealth
}
func (g *Game) EndDay() {
	g.inEndScreen = true
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Disconnect.")
	g := &Game{
		t:         0,
		uiContext: superui.NewUIContext(),

		inventory: [8]*Item{
			{id: "string"},
			{id: "string"},
			{id: "rod"},
			{id: "box"},
		},
		sublevel:  "sewer_entrance",
		sublevels: createSublevels(),

		tasks: getTaskList(),

		Health:    20,
		MaxHealth: 20,

		templateItem: &Item{
			id:         "template",
			resultData: "Mind Control",
		},
	}

	g.StartDay()

	g.hudUI = createHudUi(g.uiContext, g)
	g.craftingUI = createCraftingUi(g.uiContext, g)
	g.endScreenUI = CreateEndScreen(g.uiContext, g)
	g.todoUI = createTodoUI(g.uiContext, g)

	p := &Character{
		position:        Vector{14, 6.5},
		startLerpT:      -1000,
		facingDirection: Vector{-1, 0},
		walkSpeed:       .11,
		speedMultiplier: 1,
	}
	g.player = p

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

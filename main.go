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

	inventory    [5]*Item
	selectedSlot int

	inCraftingUi   bool
	selectedRecipe int

	sublevel string

	sublevels map[string]*Sublevel

	// Relative to grid
	targetX, targetY float64

	selectionName    string
	selectionInRange bool

	day int
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
	g.t += 1

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
			openCraftingUI(g)
		}
	}

	g.handlePlayerMovement()

	g.setTileToWall()

	g.uiContext.PreUpdate()
	if g.inCraftingUi {
		g.craftingUI.Update()
	}
	g.hudUI.Update()

	g.uiContext.Update()

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

	clickUseType := ""

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

	if cursorSelectionTile != nil && cursorSelectionTile.Type == "vent_down" {
		clickUseType = "vent_down"
	}
	if cursorSelectionTile != nil && cursorSelectionTile.Type == "box" {
		clickUseType = "box"
	}

	switch clickUseType {
	case "conveyor_item":
		g.selectionName = itemData[selectedConveyorItem.itemType.id].name
		g.targetX, g.targetY = selectedConveyorItem.X, selectedConveyorItem.Y
	case "item":
		g.selectionName = itemData[selectedItem.itemType.id].name
		g.targetX, g.targetY = selectedItem.X, selectedItem.Y
	case "vent_down":
		g.selectionName = "Vent"
		g.targetX, g.targetY = float64(cursorX/16), float64(cursorY/16)
	case "box":
		g.selectionName = "Box"
		g.targetX, g.targetY = float64(cursorX/16), float64(cursorY/16)
	default:
		g.selectionName = ""
		g.targetX, g.targetY = float64(cursorX/16), float64(cursorY/16)
	}

	distanceBetweenPlayerAndTarget := math.Sqrt(math.Pow(g.targetX-g.player.position.X, 2) + math.Pow(g.targetY-g.player.position.Y, 2))

	g.selectionInRange = distanceBetweenPlayerAndTarget < 3
	if clickUseType == "vent_down" {
		g.selectionInRange = distanceBetweenPlayerAndTarget < 1
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
		case "vent_down":
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
				if g.CurrentSublevel().adjacentSpaces.Down != "" {
					g.sublevel = g.CurrentSublevel().adjacentSpaces.Down
				}
			}
		case "box":
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		dropItemSlot(g, g.selectedSlot)
	}

	if tileImOn != nil && tileImOn.Type == "box" {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			if g.CurrentSublevel().adjacentSpaces.Up != "" {
				g.sublevel = g.CurrentSublevel().adjacentSpaces.Up
			}
		}
	}

	if hasSelection == false {
		g.ItemUseEvents()
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

func (g *Game) ItemUseEvents() {
	heldItem := g.inventory[g.selectedSlot]
	if heldItem == nil {
		return
	}

	if heldItem.id == "box" && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		cursorX, cursorY := ebiten.CursorPosition()
		targetX, targetY := (cursorX / 16), (cursorY / 16)

		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "box",
		}

		g.inventory[g.selectedSlot] = nil
	}
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

	if tileImOn != nil && tileImOn.Type == "box" {
		screen.DrawImage(vents, nil)
	}

	for _, e := range g.CurrentSublevel().Enemies {
		e.Draw(screen, g, 0, 0, false)
	}

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
			screen.DrawImage(target, op)
		}
	}

	screen.DrawImage(g.CurrentSublevel().Vignette, nil)

	if g.inCraftingUi {
		g.craftingUI.Draw(screen)
	}

	g.hudUI.Draw(screen)

	op := &text.DrawOptions{}
	op.GeoM.Translate(320-4, 0)
	op.PrimaryAlign = text.AlignEnd

	text.Draw(screen, fmt.Sprint("Day ", g.day), &text.GoTextFace{
		Source: fontFaceSource,
		Size:   16,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(320-4, 2+16)
	op.PrimaryAlign = text.AlignEnd

	text.Draw(screen, g.CurrentSublevel().Title, &text.GoTextFace{
		Source: fontFaceSource,
		Size:   8,
	}, op)

	smallFontFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   8,
	}

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

	// ops := &ebiten.DrawRectShaderOptions{}
	// ops.Images[0] = screen

	// shaded := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	// shaded.DrawRectShader(320, 240, crtshader, ops)
	// screen.DrawImage(shaded, nil)
	// ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	g := &Game{
		t:         0,
		uiContext: superui.NewUIContext(),

		inventory: [5]*Item{
			{id: "string"},
			{id: "string"},
			{id: "rod"},
			{id: "box"},
		},
		sublevel:  "sewer_entrance",
		sublevels: createSublevels(),
	}

	g.hudUI = createHudUi(g.uiContext, g)
	g.craftingUI = createCraftingUi(g.uiContext, g)

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

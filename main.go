package main

import (
	"hatchi/disconnect/superui"
	"image"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
}

func (g *Game) CurrentSublevel() *Sublevel {
	return g.sublevels[g.sublevel]
}

func (g *Game) Update() error {
	g.t += 1

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

	// handle item pickup
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		for itemIndex, item := range g.CurrentSublevel().inGameItems {
			distance := VectorMagnitude(VectorSubtract(Vector{X: float64(item.X), Y: float64(item.Y)}, Vector{X: float64(targetX), Y: float64(targetY)}))
			if distance < 1 {
				replaceSlot := g.selectedSlot
				if g.inventory[g.selectedSlot] == nil {
					replaceSlot = g.selectedSlot
				} else {
					emptySlot := false
					for i, slot := range g.inventory {
						if slot == nil {
							replaceSlot = i
							emptySlot = true
							break
						}
					}
					if !emptySlot {
						dropItemSlot(g, g.selectedSlot)
						replaceSlot = g.selectedSlot
					}
				}

				g.inventory[replaceSlot] = item.itemType

				g.CurrentSublevel().inGameItems = slices.Delete(g.CurrentSublevel().inGameItems, itemIndex, itemIndex+1)

				break
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		dropItemSlot(g, g.selectedSlot)
	}

	playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

	if tileImOn != nil && tileImOn.Type == "box" {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.sublevel = "factory"
		}
	} else {
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
	if TileIsSolid(g.CurrentSublevel().tileMap[int(adjacentVectorX.Y)][int(adjacentVectorX.X)]) {
		xSpeed = 0
	}

	adjacentVectorY := VectorFloor(VectorAdd(playerCenter, Vector{X: 0, Y: ySpeed + yOffset}))
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
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "wall",
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "conveyor_left",
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key9) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "conveyor_down",
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key8) {
		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "machine",
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.CurrentSublevel().conveyorItems = append(g.CurrentSublevel().conveyorItems,
			&ConveyorItem{
				X: float64(cursorX / 16),
				Y: float64(cursorY / 16),
				itemType: &Item{
					id: "circuit_board_finished",
				},
			},
		)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")

	switch g.sublevel {
	case "sewer":
		screen.DrawImage(test_bg, nil)
	case "factory":
		screen.DrawImage(factory_base, nil)
	}

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
			}
		}
	}

	playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

	playerIsAccended := TileIsAccended(tileImOn)

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

	cursorX, cursorY := ebiten.CursorPosition()

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

	if g.sublevel == "factory" {
		screen.DrawImage(right_wall_conveyor_overlay, nil)
		screen.DrawImage(top_wall_conveyor_overlay, nil)
		screen.DrawImage(left_wall_conveyor_overlay, nil)
	}

	if !g.uiContext.IsHovered() {
		targetX, targetY := cursorX/16, cursorY/16

		targetScreenX, targetScreenY := (targetX)*16, (targetY)*16

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(targetScreenX), float64(targetScreenY))

		overItem := false
		for _, item := range g.CurrentSublevel().inGameItems {
			distance := VectorMagnitude(VectorSubtract(Vector{X: float64(item.X), Y: float64(item.Y)}, Vector{X: float64(targetX), Y: float64(targetY)}))
			if distance < 1 {
				overItem = true
			}
		}

		if overItem {
			screen.DrawImage(target_green, op)
		} else {
			screen.DrawImage(target, op)
		}
	}

	if g.sublevel == "sewer" {
		screen.DrawImage(vignette, nil)
	} else {
		screen.DrawImage(vignette_mild, nil)
	}

	if g.inCraftingUi {
		g.craftingUI.Draw(screen)
	}

	for _, e := range g.CurrentSublevel().Enemies {
		e.Draw(screen, g, 0, 0, false)
	}

	g.hudUI.Draw(screen)
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
			{id: "tape"},
			{id: "tape"},
			{id: "tape"},
			{id: "box"},
		},
		sublevel:  "sewer",
		sublevels: createSublevels(),
	}

	g.hudUI = createHudUi(g.uiContext, g)
	g.craftingUI = createCraftingUi(g.uiContext, g)

	p := &Character{
		position:        Vector{1, 7},
		startLerpT:      -1000,
		facingDirection: Vector{1, 0},
		walkSpeed:       .11,
		speedMultiplier: 1,
	}
	g.player = p

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

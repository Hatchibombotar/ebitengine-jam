package main

import (
	"hatchi/disconnect/superui"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

	g.ItemUseEvents()

	return nil
}

func (g *Game) ItemUseEvents() {
	heldItemId := g.inventory[g.selectedSlot].id

	if heldItemId == "box" && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		cursorX, cursorY := ebiten.CursorPosition()
		targetX, targetY := (cursorX / 16), (cursorY / 16)

		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "box",
		}
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
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		cursorX, cursorY := ebiten.CursorPosition()
		targetX, targetY := (cursorX / 16), (cursorY / 16)

		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
			Type: "wall",
		}
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
			}
		}
	}

	playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

	if tileImOn != nil && tileImOn.Type == "box" {
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

	screen.DrawImage(vignette, nil)

	cursorX, cursorY := ebiten.CursorPosition()

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

	// debug draw
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

	if g.inCraftingUi {
		g.craftingUI.Draw(screen)
	}

	g.hudUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func createHudUi(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
	ui := superui.NewUI(uiContext)

	craftingButton := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			Padding:      superui.Spacing{Top: 4, Left: 4, Right: 4, Bottom: 4},
			PositionMode: superui.PositionFixed,
			X:            8,
			Y:            8,
			CursorShape:  ebiten.CursorShapePointer,
			IsFocusable:  true,

			OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
				if g.inCraftingUi {
					superui.FillNineSlice(screen, widget, button_nine_slice_inverted, 3)
				} else {
					superui.FillNineSlice(screen, widget, button_nine_slice, 3)
				}
			},

			OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
					g.inCraftingUi = !g.inCraftingUi
				}
			},
		},
		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				Width:      16,
				Height:     16,
				WidthMode:  superui.SizeFixed,
				HeightMode: superui.SizeFixed,

				OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(widget.GetResultX()), float64(widget.GetResultY()))
					screen.DrawImage(hammer, op)
				},
			},
		),
	)

	hotbar := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			PositionMode: superui.PositionFixed,
			X:            8,
			Y:            240 - 8 - 24,
			CursorShape:  ebiten.CursorShapePointer,
			IsFocusable:  true,

			LayoutDirection: superui.LayoutRow,
			Gap:             2,
		},
	)

	for slotIndex := range 5 {
		slot := superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				Width:      24,
				Height:     24,
				WidthMode:  superui.SizeFixed,
				HeightMode: superui.SizeFixed,

				IsFocusable: true,

				OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(widget.GetResultX()), float64(widget.GetResultY()))
					if g.selectedSlot == slotIndex {
						screen.DrawImage(hotbar_slot, op)
					} else {
						screen.DrawImage(hotbar_slot_unselected, op)
					}

					inventorySlot := g.inventory[slotIndex]
					if inventorySlot != nil {
						img := itemData[inventorySlot.id].image
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(float64(widget.GetResultX()+4), float64(widget.GetResultY()+4))
						screen.DrawImage(img, op)
					}
				},

				OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
					if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.HasFocusOn(w) {
						g.selectedSlot = slotIndex
					} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) && root.HasFocusOn(w) {
						dropItemSlot(g, g.selectedSlot)
					}
				},
			},
		)
		hotbar.AddChild(slot)
	}

	ui.AddChild(craftingButton)
	ui.AddChild(hotbar)

	return ui
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
	}
	g.player = p

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

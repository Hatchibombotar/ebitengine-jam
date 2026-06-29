package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const TRANISITION_TIME_HORIZ = 20
const TRANISITION_TIME_VERT = 20

func (g *Game) DrawSpace(screen *ebiten.Image, currentSublevel *Sublevel, drawPlayer bool) {
	waterT := g.t / 4
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((waterT)%16)-16, 0)
	screen.DrawImage(water_top, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-(float64((waterT)%16) + 16), 0)
	screen.DrawImage(water_bottom, op)

	screen.DrawImage(
		currentSublevel.Background, nil,
	)

	for y, row := range currentSublevel.tileMap {
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

	tileImOn := currentSublevel.tileMap[playerY][playerX]

	playerIsAccended := TileIsAccended(tileImOn)
	if drawPlayer {

		if !g.inExclusiveUIMode {
			if playerIsAccended {
				g.player.Draw(screen, g, 0, -6, false)
			} else {
				g.player.Draw(screen, g, 0, 0, false)
			}
		}
	}

	for _, inGameItem := range currentSublevel.inGameItems {
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
			// op.GeoM.Translate(0, math.Sin(float64(g.t)*0.5)*1)
		}

		screen.DrawImage(itemData.image, op)
	}

	for _, inGameItem := range currentSublevel.conveyorItems {
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
			// op.GeoM.Translate(0, math.Sin(float64(g.t)*0.5)*1)
		}

		screen.DrawImage(itemData.image, op)
	}

	for y, row := range currentSublevel.tileMap {
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

	if currentSublevel.Overlay != nil {
		screen.DrawImage(currentSublevel.Overlay, nil)
	}

	if currentSublevel.adjacentSpaces.Up != "" {
		aboveSpaceId := currentSublevel.adjacentSpaces.Up
		aboveSpace := g.spaces[aboveSpaceId]

		standingOnBox := tileImOn != nil && tileImOn.Type == "box"

		for y, row := range aboveSpace.tileMap {
			for x, tile := range row {
				if tile == nil {
					continue
				}
				currentTile := currentSublevel.tileMap[y][x]
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

	if !g.inExclusiveUIMode {
		for _, e := range currentSublevel.Enemies {
			e.Draw(screen, g, 0, math.Sin(float64(g.t)*0.1)*1, false)

			isAttacking := e.AttackStartT+e.AttackLength > g.t
			if !isAttacking {
				continue
			}

			t := g.t / 2
			op := &ebiten.DrawImageOptions{}

			dx, dy := (e.position.X - g.player.position.X), (e.position.Y - g.player.position.Y)
			m := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
			ndx, ndy := (e.position.X-g.player.position.X)/m, (e.position.Y-g.player.position.Y)/m
			angle := math.Atan2(dy, dx)
			op.GeoM.Translate(-8, -8)
			op.GeoM.Scale(-1, 1)
			op.GeoM.Rotate(angle)
			op.GeoM.Translate(8, 8)
			op.GeoM.Translate(e.position.X*16, e.position.Y*16)
			op.GeoM.Translate(-ndx*15, -ndy*15)
			screen.DrawImage(
				zap_flipbook.SubImage(image.Rect(0, (t%4)*16, 16, 16+((t%4)*16))).(*ebiten.Image),
				op,
			)
		}

	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	isMoving := g.prevLevelMoveTime+TRANISITION_TIME_HORIZ > g.t

	isMovingHorizonally := false
	isMovingVertically := false
	switch g.prevLevelDirection {
	case "North", "South", "East", "West":
		isMovingHorizonally = true
	case "Up", "Down":
		isMovingVertically = true
	}

	if isMoving && isMovingHorizonally {
		p := float64(g.t-g.prevLevelMoveTime) / TRANISITION_TIME_HORIZ
		oldScreen := ebiten.NewImage(320, 240)
		newScreen := ebiten.NewImage(320, 240)
		var oldSpaceId string
		translationX, translationY := 0.0, 0.0
		switch g.prevLevelDirection {
		case "East":
			oldSpaceId = g.CurrentSublevel().adjacentSpaces.West
			translationX = -320
		case "West":
			oldSpaceId = g.CurrentSublevel().adjacentSpaces.East
			translationX = 320
		case "North":
			oldSpaceId = g.CurrentSublevel().adjacentSpaces.South
			translationY = 240
		case "South":
			oldSpaceId = g.CurrentSublevel().adjacentSpaces.North
			translationY = -240
		}

		if oldSpaceId != "" {
			oldSpace := g.spaces[oldSpaceId]

			g.DrawSpace(oldScreen, oldSpace, false)
			g.DrawSpace(newScreen, g.CurrentSublevel(), true)

			op1 := &ebiten.DrawImageOptions{}
			op1.GeoM.Translate(translationX*p, translationY*p)
			screen.DrawImage(oldScreen, op1)
		}
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate((translationX*p)-translationX, (translationY*p)-translationY)
		screen.DrawImage(newScreen, op2)
	} else if isMoving && isMovingVertically {

		g.DrawSpace(screen, g.CurrentSublevel(), true)
		p := float64(g.t-g.prevLevelMoveTime) / TRANISITION_TIME_VERT

		var opacity uint8 = 0
		opacity = uint8(255.0 * (1 - p))

		if p < 0.3 {
			opacity = uint8(255.0 * (1 - p))
		} else {
			opacity = 0
		}

		vector.FillRect(screen, 0, 0, 320, 240, color.RGBA{0, 0, 0, opacity}, false)

	} else {
		g.DrawSpace(screen, g.CurrentSublevel(), true)

	}

	// playerX, playerY := int(g.player.position.X+0.5), int(g.player.position.Y+0.5)

	// tileImOn := g.CurrentSublevel().tileMap[playerY][playerX]

	// cursorX, cursorY := ebiten.CursorPosition()

	// debug draw
	// if ebiten.IsKeyPressed(ebiten.KeyF8) {
	// 	for y, row := range g.CurrentSublevel().tileMap {
	// 		for x, tile := range row {
	// 			if tile == nil {
	// 				continue
	// 			}
	// 			op := &ebiten.DrawImageOptions{}
	// 			op.GeoM.Translate(float64(x*16), float64(y*16))

	// 			switch tile.Type {
	// 			case "wall":
	// 				screen.DrawImage(debug_wall, op)
	// 			}
	// 		}
	// 	}
	// }

	ops := &ebiten.DrawRectShaderOptions{}
	ops.Images[0] = screen

	shaded := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	shaded.DrawRectShader(320, 240, crtshader, ops)
	screen.DrawImage(shaded, nil)

	if !g.uiContext.IsHovered() && !g.inExclusiveUIMode {
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

	if g.isElectricityDown && g.CurrentSublevel().Vignette == vignette_mild {
		screen.DrawImage(vignette, nil)
	}

	if g.isLockDown {
		screen.DrawImage(vignette_red, nil)
	}

	if g.inExclusiveUIMode {
		vector.FillRect(screen, 0, 0, 320, 240, color.RGBA{0, 0, 0, 150}, false)

		if g.inMainMenu {
			DrawMainMenu(screen, g)
		}
		if g.dayEndedInDeath {
			g.lossScreen.Draw(screen)
		}
		if g.inIntroScreen {
			g.introScreen.Draw(screen)
		}
		if g.inEndScreen {
			g.endScreenUI.Draw(screen)
		}
		if g.inWinScreen {
			g.winScreen.Draw(screen)
		}
		return
	}

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

	if g.isElectricityDown {
		totalSeconds := g.electricityDownRemainingTime / 60
		minutes := totalSeconds / 60
		seconds := totalSeconds % 60

		formattedTimeRemaining := fmt.Sprintf("%02d:%02d", minutes, seconds)

		timeWidth, _ := text.Measure(formattedTimeRemaining, smallFontFace, 1)
		op = &text.DrawOptions{}
		op.GeoM.Translate(320-4-timeWidth, 2+16+16)
		op.PrimaryAlign = text.AlignEnd
		op.ColorScale.ScaleWithColor(color.Gray{180})

		text.Draw(screen, "Electricity Down For: ", smallFontFace, op)

		op = &text.DrawOptions{}
		op.GeoM.Translate(320-4, 2+16+16)
		op.PrimaryAlign = text.AlignEnd
		op.ColorScale.ScaleWithColor(color.Gray{180})

		text.Draw(screen, formattedTimeRemaining, smallFontFace, op)
	} else if g.isLockDown {
		op = &text.DrawOptions{}
		op.GeoM.Translate(320-4, 2+16+16)
		op.PrimaryAlign = text.AlignEnd
		op.ColorScale.ScaleWithColor(color.Gray{180})

		text.Draw(screen, "EMERGENCY LOCKDOWN", smallFontFace, op)
	}

	if g.selectionName != "" && !g.uiContext.IsHovered() {
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

	if g.currentTask != nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(52, 4)
		op.PrimaryAlign = text.AlignStart
		op.ColorScale.ScaleWithColor(color.White)

		text.Draw(screen, "Current Task:", smallFontFace, op)

		op = &text.DrawOptions{}
		op.GeoM.Translate(52, 4+8+2)
		op.PrimaryAlign = text.AlignStart
		op.ColorScale.ScaleWithColor(color.Gray{180})

		text.Draw(screen, g.currentTask.description, smallFontFace, op)
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

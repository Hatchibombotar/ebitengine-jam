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

	if !g.inExclusiveUIMode {
		if playerIsAccended {
			g.player.Draw(screen, g, 0, -6, false)
		} else {
			g.player.Draw(screen, g, 0, 0, false)
		}
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
			// op.GeoM.Translate(0, math.Sin(float64(g.t)*0.5)*1)
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
			// op.GeoM.Translate(0, math.Sin(float64(g.t)*0.5)*1)
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

	if !g.inExclusiveUIMode {
		for _, e := range g.CurrentSublevel().Enemies {
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

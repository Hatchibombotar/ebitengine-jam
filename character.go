package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
	startLerpT int

	position Vector

	facingDirection Vector
	spriteIndex     int
	walkSpeed       float64
	deathPhase      int

	speedMultiplier int

	isMoving bool
}

func (c *Character) GetLerpProgress(g *Game) float64 {
	return min((float64(g.t)-float64(c.startLerpT))*c.walkSpeed, 1)
}

func (c *Character) Draw(screen *ebiten.Image, g *Game, offsetX, offsetY float64, isTransition bool) {
	op := &ebiten.DrawImageOptions{}

	// c.visiblePosition = LerpVectors(
	// 	c.startPositon,
	// 	c.endPosition,
	// 	c.GetLerpProgress(g),
	// )

	spriteDirection := 0

	frame := 1
	if c.GetLerpProgress(g) < 1 || c.isMoving {
		frame = (c.speedMultiplier * (g.t / 6)) % 4
	}

	if VectorIs(c.facingDirection, Vector{0, 1}) {
		spriteDirection = 0 + frame
	} else if VectorIs(c.facingDirection, Vector{0, -1}) {
		spriteDirection = 8 + frame
	} else if VectorIs(c.facingDirection, Vector{-1, 0}) {
		spriteDirection = 4 + frame
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(16, 0)
	} else if VectorIs(c.facingDirection, Vector{1, 0}) {
		spriteDirection = 4 + frame
	}

	if c.deathPhase > 0 {
		spriteDirection = 12 + c.deathPhase - 1
	}

	op.GeoM.Translate(float64(int(c.position.X*16)), float64(int(c.position.Y*16)))
	op.GeoM.Translate(offsetX, offsetY)

	if isTransition {
		op.ColorScale.ScaleWithColor(color.Gray{50})
	}

	img := spritesheet.SubImage(image.Rect(spriteDirection*16, c.spriteIndex*16, (1+spriteDirection)*16, (1+c.spriteIndex)*16)).(*ebiten.Image)
	screen.DrawImage(img, op)
}

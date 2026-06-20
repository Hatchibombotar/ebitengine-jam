package superui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func FillSolid(screen *ebiten.Image, widget GenericWidget, fillColor color.Color) {
	vector.DrawFilledRect(
		screen,
		float32(widget.GetResultX()), float32(widget.GetResultY()),
		float32(widget.GetResultWidth()), float32(widget.GetResultHeight()),
		fillColor,
		false,
	)
}

// written by ai
func FillNineSlice(screen *ebiten.Image, widget GenericWidget, nineSliceImage *ebiten.Image, nineSliceWidth int) {
	destX := widget.GetResultX()
	destY := widget.GetResultY()
	destW := widget.GetResultWidth()
	destH := widget.GetResultHeight()

	// Source slice sizes
	slice := nineSliceWidth
	imgW, imgH := nineSliceImage.Bounds().Dx(), nineSliceImage.Bounds().Dy()

	// Destination rectangles
	midW := destW - 2*slice
	midH := destH - 2*slice

	if midW < 0 || midH < 0 {
		// Not enough space to draw nine-slice
		return
	}

	// Helper to draw a region
	drawRegion := func(sx, sy, sw, sh int, dx, dy, dw, dh float64) {
		src := nineSliceImage.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(dw/float64(sw), dh/float64(sh))
		op.GeoM.Translate(dx, dy)
		screen.DrawImage(src, op)
	}

	// Top-left
	drawRegion(0, 0, slice, slice, float64(destX), float64(destY), float64(slice), float64(slice))
	// Top
	drawRegion(slice, 0, imgW-2*slice, slice, float64(destX+slice), float64(destY), float64(midW), float64(slice))
	// Top-right
	drawRegion(imgW-slice, 0, slice, slice, float64(destX+slice+midW), float64(destY), float64(slice), float64(slice))

	// Left
	drawRegion(0, slice, slice, imgH-2*slice, float64(destX), float64(destY+slice), float64(slice), float64(midH))
	// Center
	drawRegion(slice, slice, imgW-2*slice, imgH-2*slice, float64(destX+slice), float64(destY+slice), float64(midW), float64(midH))
	// Right
	drawRegion(imgW-slice, slice, slice, imgH-2*slice, float64(destX+slice+midW), float64(destY+slice), float64(slice), float64(midH))

	// Bottom-left
	drawRegion(0, imgH-slice, slice, slice, float64(destX), float64(destY+slice+midH), float64(slice), float64(slice))
	// Bottom
	drawRegion(slice, imgH-slice, imgW-2*slice, slice, float64(destX+slice), float64(destY+slice+midH), float64(midW), float64(slice))
	// Bottom-right
	drawRegion(imgW-slice, imgH-slice, slice, slice, float64(destX+slice+midW), float64(destY+slice+midH), float64(slice), float64(slice))
}

func CreateFillSolid(fillColor color.Color) func(screen *ebiten.Image, widget GenericWidget, _ *UIContainer) {
	return func(screen *ebiten.Image, widget GenericWidget, _ *UIContainer) {
		FillSolid(screen, widget, fillColor)
	}
}

// TODO: add image helpers

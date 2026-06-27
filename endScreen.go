package main

import (
	"hatchi/disconnect/superui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DrawEndScreen(screen *ebiten.Image) {
	screen.Fill(color.Black)

}

func UpdateEndScreen() {

}

func CreateEndScreen(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
	ui := superui.NewUI(uiContext)

	endScreenUI := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			AlignHorizontal: superui.AlignCenter,
			WidthMode:       superui.SizeFixed,
			HeightMode:      superui.SizeFixed,
			Width:           320,
			Height:          240,
			Gap:             4,
			Padding:         superui.Spacing{Top: 24},
		},
		superui.NewTextWidget(
			&superui.TextWidgetOps{
				Face: &text.GoTextFace{
					Source: fontFaceSource,
					Size:   16,
				},
				Color:         color.White,
				WrapBehaviour: superui.NoWrap,
			},
			"Escaped!",
		),

		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				HeightMode: superui.SizeFixed,
				Height:     64,
			},
		),

		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
					if canCraftRecipe(g, g.selectedRecipe) {
						superui.FillNineSlice(screen, widget, button_nine_slice, 3)
					} else {
						superui.FillNineSlice(screen, widget, button_nine_slice_disabled, 3)
					}
				},
				CursorShape: ebiten.CursorShapePointer,
				Padding:     superui.Spacing{Top: 3, Right: 5, Bottom: 5, Left: 5},
				IsFocusable: true,
				OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
					if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
						g.StartDay()
					}
				},
			},
			superui.NewTextWidget(
				&superui.TextWidgetOps{
					Face: &text.GoTextFace{
						Source: fontFaceSource,
						Size:   8,
					},
					Color:         color.White,
					WrapBehaviour: superui.NoWrap,
				},
				"Next Day",
			),
		),
	)

	ui.AddChild(endScreenUI)

	return ui
}

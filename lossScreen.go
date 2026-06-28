package main

import (
	"hatchi/disconnect/superui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func CreateLossScreen(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
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
			"You died...",
		),

		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				HeightMode: superui.SizeFixed,
				Height:     28,
			},
		),

		superui.NewTextWidget(
			&superui.TextWidgetOps{
				Face: &text.GoTextFace{
					Source: fontFaceSource,
					Size:   8,
				},
				Color:         color.White,
				WrapBehaviour: superui.NoWrap,
			},
			"Message from HQ:",
		),
		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
					superui.FillNineSlice(screen, widget, hud_button_nine_slice, 3)
				},
				WidthMode: superui.SizeFixed,
				Width:     148,
				Padding:   superui.Spacing{Top: 2, Right: 6, Bottom: 0, Left: 6},
			},
			superui.NewTextWidget(
				&superui.TextWidgetOps{
					Face: &text.GoTextFace{
						Source: fontFaceSource,
						Size:   8,
					},
					Color:         color.White,
					WrapBehaviour: superui.WrapText,
				},
				"Sorry to hear of your death but the world still needs saving. Have another go?\n",
			),
		),

		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				HeightMode: superui.SizeFixed,
				Height:     100,
			},
		),

		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				PositionMode:    superui.PositionFixed,
				AlignHorizontal: superui.AlignCenter,
				WidthMode:       superui.SizeFixed,
				HeightMode:      superui.SizeFixed,
				Width:           320,
				Height:          240,
				Gap:             4,
				Padding:         superui.Spacing{Top: 182},
			},
			superui.NewBoxWidget(
				&superui.BoxWidgetOps{
					OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
						if root.IsHovered(widget) {
							superui.FillNineSlice(screen, widget, hud_button_nine_slice_inverted, 3)
						} else {
							superui.FillNineSlice(screen, widget, hud_button_nine_slice, 3)
						}
					},
					CursorShape: ebiten.CursorShapePointer,
					Padding:     superui.Spacing{Top: 3, Right: 5, Bottom: 5, Left: 5},
					IsFocusable: true,
					OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
						if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
							g.Init()
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
					"Start Over",
				),
			),
		),
	)

	ui.AddChild(endScreenUI)

	return ui
}

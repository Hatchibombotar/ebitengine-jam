package main

import (
	"hatchi/disconnect/superui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DrawMainMenu(screen *ebiten.Image, g *Game) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(320/2, 16)
	op.PrimaryAlign = text.AlignCenter

	text.Draw(screen, `"Join Us"`, &text.GoTextFace{
		Source: fontFaceSource,
		Size:   48,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(320/2, 16+48+8)
	op.PrimaryAlign = text.AlignCenter

	text.Draw(screen, "A game by Hatchibombotar", &text.GoTextFace{
		Source: fontFaceSource,
		Size:   8,
	}, op)

	g.mainMenuUi.Draw(screen)
}

func CreateMainMenuUi(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
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

		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				HeightMode: superui.SizeFixed,
				Height:     150,
			},
		),

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
				Padding:     superui.Spacing{Top: 1, Right: 9, Bottom: 6, Left: 9},
				IsFocusable: true,
				OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
					if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
						g.inMainMenu = false
						g.inIntroScreen = true
					}
				},
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
				"Start",
			),
		),
	)

	ui.AddChild(endScreenUI)

	return ui
}

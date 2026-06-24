package main

import (
	"hatchi/disconnect/superui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

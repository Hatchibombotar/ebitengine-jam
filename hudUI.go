package main

import (
	"hatchi/disconnect/superui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func createHudUi(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
	ui := superui.NewUI(uiContext)

	topButtonBar := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			PositionMode:    superui.PositionFixed,
			LayoutDirection: superui.LayoutRow,
			X:               4,
			Y:               4,
			Gap:             4,
		},
	)

	craftingButton := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			Padding: superui.Spacing{Top: 2, Left: 2, Right: 2, Bottom: 2},

			CursorShape: ebiten.CursorShapePointer,
			IsFocusable: true,

			OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
				if root.IsHovered(widget) || g.inCraftingUi {
					superui.FillNineSlice(screen, widget, hud_button_nine_slice_inverted, 2)
				} else {
					superui.FillNineSlice(screen, widget, hud_button_nine_slice, 2)
				}
			},

			OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
					g.inCraftingUi = !g.inCraftingUi
					if g.inCraftingUi {
						g.inTodoUI = false
						openCraftingUI(g)
					}
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

	todoListButton := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			Padding: superui.Spacing{Top: 2, Left: 2, Right: 2, Bottom: 2},

			CursorShape: ebiten.CursorShapePointer,
			IsFocusable: true,

			OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
				if root.IsHovered(widget) || g.inTodoUI {
					superui.FillNineSlice(screen, widget, hud_button_nine_slice_inverted, 2)
				} else {
					superui.FillNineSlice(screen, widget, hud_button_nine_slice, 2)
				}
			},

			OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
					g.inTodoUI = !g.inTodoUI
					if g.inTodoUI {
						g.inCraftingUi = false
					}
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
					screen.DrawImage(todo_list, op)
				},
			},
		),
	)

	topButtonBar.AddChild(craftingButton)
	topButtonBar.AddChild(todoListButton)

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
						img := GetHeldItemSprite(itemData[inventorySlot.id])
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

	ui.AddChild(topButtonBar)
	ui.AddChild(hotbar)

	return ui
}

func GetOrFreeSlotForItemInHotbar(g *Game) int {
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
	return replaceSlot
}

package main

import (
	"hatchi/disconnect/superui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func openCraftingUI(g *Game) {
	g.selectedRecipe = -1
}

func createCraftingUi(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
	ui := superui.NewUI(uiContext)

	craftingUi := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			PositionMode: superui.PositionFixed,
			X:            8,
			Y:            4 + 32 + 8 + 4,

			LayoutDirection: superui.LayoutRow,
			Gap:             2,

			Padding: superui.Spacing{Top: 4, Left: 4, Right: 4, Bottom: 4},

			OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
				superui.FillNineSlice(screen, widget, box_nine_slice, 3)
			},
		},
	)

	recipeList := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			LayoutType: superui.LayoutGrid,
			Columns:    5,
			Gap:        3,
		},
	)
	craftingUi.AddChild(recipeList)

	// Add divider
	craftingUi.AddChild(superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			LayoutDirection: superui.LayoutRow,
			WidthMode:       superui.SizeFixed,
			HeightMode:      superui.SizeFixed,
			Width:           3,
			Height:          128,
			OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(widget.GetResultX()), float64(widget.GetResultY()))

				screen.DrawImage(crafting_divider, op)
			},
		},
	))

	for recipeIndex, recipe := range recipeData {
		item := itemData[recipe.result]

		recipeButton := superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				WidthMode:   superui.SizeFixed,
				HeightMode:  superui.SizeFixed,
				Width:       16,
				Height:      16,
				CursorShape: ebiten.CursorShapePointer,
				IsFocusable: true,

				OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(widget.GetResultX()), float64(widget.GetResultY()))
					if canCraftRecipe(g, recipeIndex) {
						screen.DrawImage(recipe_slot_active, op)
					} else {
						screen.DrawImage(recipe_slot, op)
					}
					screen.DrawImage(GetHeldItemSprite(item), op)
				},

				OnInputUpdate: func(w superui.GenericWidget, root *superui.UIContainer) {
					if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && root.IsHovered(w) {
						g.selectedRecipe = recipeIndex
					}
				},
			},
		)
		recipeList.AddChild(recipeButton)
	}

	// Recipe Example
	craftingUi.AddChild(
		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				LayoutDirection: superui.LayoutRow,
				WidthMode:       superui.SizeFixed,
				HeightMode:      superui.SizeFixed,
				Width:           64,
			},
			superui.NewConditionalDisplay(
				&superui.ConditionalDisplayOps{
					ShouldShow: func() bool {
						return g.selectedRecipe != -1
					},
				},
				superui.NewBoxWidget(
					&superui.BoxWidgetOps{
						Gap: 4,
					},
					superui.NewTextWidget(
						&superui.TextWidgetOps{
							Face: &text.GoTextFace{
								Source: fontFaceSource,
								Size:   8,
							},
							Color:         color.Black,
							WrapBehaviour: superui.NoWrap,
							GetDynamicText: func() string {
								if g.selectedRecipe == -1 {
									return "<item>"
								}
								selectedRecipeResultItem := recipeData[g.selectedRecipe].result
								return itemData[selectedRecipeResultItem].name
							},
						},
						"<item>",
					),
					superui.NewBoxWidget(
						&superui.BoxWidgetOps{
							HeightMode: superui.SizeFixed,
							Height:     1,
						},
					),
					superui.NewTextWidget(
						&superui.TextWidgetOps{
							Face: &text.GoTextFace{
								Source: fontFaceSource,
								Size:   8,
							},
							Color:         color.Gray{60},
							WrapBehaviour: superui.NoWrap,
						},
						"Ingredients",
					),
					superui.NewBoxWidget(
						&superui.BoxWidgetOps{
							LayoutType: superui.LayoutGrid,
							Columns:    3,
							Gap:        3,
						},
						func() []superui.GenericWidget {
							result := []superui.GenericWidget{}

							for ingredientIndex := range 3 {
								ingredientButton := superui.NewConditionalDisplay(
									&superui.ConditionalDisplayOps{
										ShouldShow: func() bool {
											ingredients := recipeData[g.selectedRecipe].ingredients
											return len(ingredients) > ingredientIndex
										},
									},
									superui.NewBoxWidget(
										&superui.BoxWidgetOps{
											WidthMode:   superui.SizeFixed,
											HeightMode:  superui.SizeFixed,
											Width:       16,
											Height:      16,
											CursorShape: ebiten.CursorShapePointer,
											IsFocusable: true,

											OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
												ingredients := recipeData[g.selectedRecipe].ingredients
												ingredient := ingredients[ingredientIndex]

												op := &ebiten.DrawImageOptions{}
												op.GeoM.Translate(float64(widget.GetResultX()), float64(widget.GetResultY()))
												if canCraftRecipe(g, g.selectedRecipe) {
													screen.DrawImage(recipe_slot_active, op)
												} else {
													screen.DrawImage(recipe_slot, op)
												}
												screen.DrawImage(GetHeldItemSprite(itemData[ingredient]), op)
											},
										},
									),
								)

								result = append(result, ingredientButton)

							}

							return result
						}()...,
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
									craftRecipe(g, g.selectedRecipe)
								}
							},
						},
						superui.NewTextWidget(
							&superui.TextWidgetOps{
								Face: &text.GoTextFace{
									Source: fontFaceSource,
									Size:   8,
								},
								Color:         color.Black,
								WrapBehaviour: superui.NoWrap,
							},
							"Craft",
						),
					),
				),
			),
		),
	)

	ui.AddChild(craftingUi)
	return ui
}

func canCraftRecipe(g *Game, recipeIndex int) bool {
	if recipeIndex < 0 || recipeIndex >= len(recipeData) {
		return false
	}

	recipe := recipeData[recipeIndex]

	// Count what we need
	needed := make(map[string]int)
	for _, ingredient := range recipe.ingredients {
		needed[ingredient]++
	}

	// Count what we have in inventory
	have := make(map[string]int)
	for _, item := range g.inventory {
		if item != nil {
			have[item.id]++
		}
	}

	// Check if we have enough of each ingredient
	for ingredient, count := range needed {
		if have[ingredient] < count {
			return false
		}
	}

	return true
}

func craftRecipe(g *Game, recipeIndex int) {
	if !canCraftRecipe(g, recipeIndex) {
		return
	}

	recipe := recipeData[recipeIndex]

	// Count ingredients to remove
	toRemove := make(map[string]int)
	for _, ingredient := range recipe.ingredients {
		toRemove[ingredient]++
	}

	// Remove ingredients from inventory
	for i := 0; i < len(g.inventory); i++ {
		if g.inventory[i] != nil && toRemove[g.inventory[i].id] > 0 {
			toRemove[g.inventory[i].id]--
			g.inventory[i] = nil
		}
	}

	// Add result to inventory
	resultKey := recipe.result
	for i := 0; i < len(g.inventory); i++ {
		if g.inventory[i] == nil {
			g.inventory[i] = &Item{id: resultKey}
			break
		}
	}
}

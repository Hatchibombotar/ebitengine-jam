package main

import (
	"fmt"
	"hatchi/disconnect/superui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func createTodoUI(uiContext *superui.UIContext, g *Game) *superui.UIContainer {
	ui := superui.NewUI(uiContext)

	ui.AddChild(
		superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				PositionMode: superui.PositionFixed,
				X:            8,
				Y:            4 + 16 + 8,
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
				"Tasks",
			),
		),
	)

	craftingContainerRoot := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			PositionMode: superui.PositionFixed,
			X:            8,
			Y:            4 + 32 + 8 + 8,

			Gap: 2,

			Padding: superui.Spacing{Top: 0, Left: 4, Right: 4, Bottom: 4},

			OnDraw: func(screen *ebiten.Image, widget superui.GenericWidget, root *superui.UIContainer) {
				superui.FillNineSlice(screen, widget, box_nine_slice, 3)
			},
		},
	)

	craftingUi := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			Padding:         superui.Spacing{Top: 4, Bottom: 4},
			LayoutDirection: superui.LayoutRow,
			Gap:             2,
		},
	)

	craftingContainerRoot.AddChild(craftingUi)

	taskList := superui.NewBoxWidget(
		&superui.BoxWidgetOps{
			Gap: 2,
		},
	)
	craftingUi.AddChild(taskList)

	for _, task := range g.tasks {
		taskItem := superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				LayoutDirection: superui.LayoutRow,
			},
		)

		taskText := superui.NewTextWidget(
			&superui.TextWidgetOps{
				Face: &text.GoTextFace{
					Source: fontFaceSource,
					Size:   8,
				},
				Color:         color.Black,
				WrapBehaviour: superui.NoWrap,
			},
			fmt.Sprint("• ", task.description),
		)
		if task.complete {
			taskText.Op.Color = color.RGBA{15, 100, 15, 255}
		}

		taskItem.AddChild(taskText)
		taskList.AddChild(taskItem)
	}

	ui.AddChild(craftingContainerRoot)
	return ui
}

type Task struct {
	description string
	complete    bool
}

func getTaskList() []*Task {
	return []*Task{
		{
			description: "Craft a screwdriver",
			complete:    true,
		},
		{
			description: "Place a box down in Sewer (East 1)",
			complete:    true,
		},
		{
			description: "Unscrew a vent",
			complete:    false,
		},
	}
}

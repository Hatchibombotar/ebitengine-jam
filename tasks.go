package main

import (
	"fmt"
	"hatchi/disconnect/superui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func CompleteTask(g *Game, task int) {
	if g.tasks[task].complete {
		return
	}
	g.tasks[task].complete = true
	g.completedTasks = append(g.completedTasks, g.tasks[task])
	UpdateTaskRange(g)

	if g.tasks[task].resultFunc != nil {
		g.tasks[task].resultFunc(g)
	}

	switch task {
	case 3:
		e1 := createEnemy()
		e1.position = Vector{19, 14}
		e2 := createEnemy()
		e2.position = Vector{1, 14}
		g.CurrentSublevel().Enemies = append(g.CurrentSublevel().Enemies, e1, e2)
	}
}
func UpdateTaskRange(g *Game) {
	startI, endI := 0, 0
	for i, task := range g.tasks {
		if task.StartingTask {
			startI = i
		}
		g.currentTask = task
		if !task.complete {
			break
		}
	}
	for i, task := range g.tasks {
		if i <= startI {
			continue
		}
		if task.StartingTask {
			endI = i - 1
			break
		}
	}
	if endI == 0 {
		endI = len(g.tasks)
	}
	g.startTask = startI
	g.endTask = endI
}

type Task struct {
	description  string
	comment      string
	complete     bool
	StartingTask bool

	resultFunc func(g *Game)
}

func getTaskList() []*Task {
	return []*Task{
		// Day 1
		{ // [0]
			description:  "Craft a screwdriver",
			complete:     false,
			StartingTask: true,
		},
		{
			description: "Place a box down in Sewer (Entrance)",
			complete:    false,
		},
		{
			description: "Unscrew a vent",
			complete:    false,
		},
		{
			description: "Steal an item",
			complete:    false,
		},
		{
			description: "Escape!",
			complete:    false,
			comment:     "That was a close call! Craft a hammer to defend yourself from the robots tomorrow!",
			resultFunc: func(g *Game) {
				g.spaces["sewer_entrance"].inGameItems = append(
					g.spaces["sewer_entrance"].inGameItems,
					&InGameItem{itemType: &Item{id: "string"}, X: 13, Y: 1},
					&InGameItem{itemType: &Item{id: "rod"}, X: 9, Y: 11},
					&InGameItem{itemType: &Item{id: "string"}, X: 13, Y: 11},
					&InGameItem{itemType: &Item{id: "string"}, X: 6, Y: 2},
				)
				g.spaces["sewer_5"].inGameItems = append(
					g.spaces["sewer_5"].inGameItems,
					&InGameItem{itemType: &Item{id: "rod"}, X: 12, Y: 7},
					&InGameItem{itemType: &Item{id: "rod"}, X: 13, Y: 12},
				)
			},
		},
		// Day 2
		{
			description:  "Craft a Hammer",
			complete:     false,
			StartingTask: true,
		},
		{ // [6]
			description: "Craft Wire Cutters",
			complete:    false,
		},
		{
			description: "Cut through wire",
			complete:    false,
		},
		{
			description: "Climb up vent in Sewer (West 1)",
			complete:    false,
		},
		{
			description: "Kill robot to get Authentication Chip",
			complete:    false,
		},
		{
			description: "Craft an Access Card",
			complete:    false,
		},
		{
			description: "Turn power off",
			complete:    false,
			comment:     "Yeah, they really don't like you turning the power off. It does stop them from seeing you touch the production line though!",
		},
		// Day 3
		{ // [12]
			description:  "Craft Template Rewriter",
			complete:     false,
			StartingTask: true,
		},
		{
			description: "Put a Hacking USB Template in the Template Machine",
			complete:    false,
		},
		{
			description: "Collect Hacking USB chip from production line",
			complete:    false,
		},
		{
			description: "Craft Hacking USB",
			complete:    false,
		},
		// Day 4
		{ // [16]
			description:  "Put Reprogrammer Template in the Template Machine",
			complete:     false,
			StartingTask: true,
			comment:      "Almost there! We just need to take control of the central console!",
		},
		{
			description: "Climb up vent in Sewer (North 2)",
			complete:    false,
		},
		{
			description: "Take control of the central console",
			complete:    false,
		},
		{
			description: "RUN!",
			complete:    false,
		},
	}
}

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
		&superui.BoxWidgetOps{},
	)
	craftingUi.AddChild(taskList)

	for i, task := range g.tasks {
		taskItem := superui.NewBoxWidget(
			&superui.BoxWidgetOps{
				LayoutDirection: superui.LayoutRow,
				Padding:         superui.Spacing{Bottom: 2},
			},
		)

		taskItem.AddChild(
			// incomplete
			superui.NewConditionalDisplay(
				&superui.ConditionalDisplayOps{
					ShouldShow: func() bool {
						return !task.complete
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
					fmt.Sprint("• ", task.description),
				),
			),
		)
		taskItem.AddChild(
			// complete
			superui.NewConditionalDisplay(
				&superui.ConditionalDisplayOps{
					ShouldShow: func() bool {
						return task.complete
					},
				},
				superui.NewTextWidget(
					&superui.TextWidgetOps{
						Face: &text.GoTextFace{
							Source: fontFaceSource,
							Size:   8,
						},
						Color:         color.RGBA{15, 100, 15, 255},
						WrapBehaviour: superui.NoWrap,
					},
					fmt.Sprint("• ", task.description),
				),
			),
		)
		taskList.AddChild(
			superui.NewConditionalDisplay(
				&superui.ConditionalDisplayOps{
					ShouldShow: func() bool {
						return g.startTask <= i && i <= g.endTask
					},
				},
				taskItem,
			),
		)
	}

	ui.AddChild(craftingContainerRoot)
	return ui
}

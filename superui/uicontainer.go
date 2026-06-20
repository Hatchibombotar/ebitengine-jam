package superui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type UIContainer struct {
	focusedWidget GenericWidget
	children      []GenericWidget

	// cursorShape ebiten.CursorShapeType

	hoveredWidget  GenericWidget
	hoveredWidgets []GenericWidget

	root *UIContext
}

func (e *UIContainer) SetFocusOn(w GenericWidget) {
	e.focusedWidget = w
}

func (e *UIContainer) HasFocusOn(widget GenericWidget) bool {
	return e.focusedWidget == widget
}

func (e *UIContainer) Update() {
	// e.cursorShape = ebiten.CursorShapeDefault
	e.hoveredWidgets = []GenericWidget{}
	e.hoveredWidget = nil

	for _, child := range e.children {
		child.Update()
		child.UpdateInput(e)
	}

	// e.root.cursorShape = e.cursorShape

	removeFocusIfClick := true
	for _, w := range e.hoveredWidgets {
		if w.IsFocusable() {
			removeFocusIfClick = false
		}
	}
	if removeFocusIfClick && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		e.focusedWidget = nil
	}

	if e.hoveredWidget != nil {
		e.root.isHovered = true
	}
}

func (e *UIContainer) IsHovered(widget GenericWidget) bool {
	mouseX, mouseY := e.root.CursorPosition()
	r := image.Rect(widget.GetResultX(), widget.GetResultY(), widget.GetResultX()+widget.GetResultWidth(), widget.GetResultY()+widget.GetResultHeight())
	p := image.Point{mouseX, mouseY}

	return p.In(r)
}

func (e *UIContainer) Draw(screen *ebiten.Image) {
	for _, child := range e.children {
		child.Draw(screen, e)
	}

	if ebiten.IsKeyPressed(ebiten.KeyF12) {
		debugWidget := e.hoveredWidget

		if debugWidget != nil {
			vector.DrawFilledRect(
				screen,
				float32(debugWidget.GetResultX()),
				float32(debugWidget.GetResultY()),
				float32(debugWidget.GetResultWidth()),
				float32(debugWidget.GetResultHeight()),
				color.RGBA{50, 50, 255, 150},
				false,
			)
			ebitenutil.DebugPrint(screen, debugWidget.String())
		}

	}

}

func (e *UIContainer) String() string {
	str := "<Root>\n"

	for _, child := range e.children {
		switch child := child.(type) {
		case *Box:
			str += child.string("  ") + "\n"
		default:
			str += "  " + child.String() + "\n"
		}
	}
	str += "</Root>"
	return str
}

func NewUI(root *UIContext, children ...GenericWidget) *UIContainer {
	eventHandler := &UIContainer{
		children: children,
		root:     root,
	}
	return eventHandler
}

func (r *UIContainer) SetCursorShape(shape ebiten.CursorShapeType) {
	r.root.cursorShape = shape
}

func (r *UIContainer) AddChild(w GenericWidget) {
	r.children = append(r.children, w)
}

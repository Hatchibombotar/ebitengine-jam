package superui

import "github.com/hajimehoshi/ebiten/v2"

type UIContext struct {
	GetCursorPositionFunc func() (int, int)
	cursorShape           ebiten.CursorShapeType
	isHovered             bool
}

func (r *UIContext) IsHovered() bool {
	return r.isHovered
}

func NewUIContext() *UIContext {
	root := &UIContext{}

	return root
}

func (r *UIContext) CursorPosition() (int, int) {
	if r.GetCursorPositionFunc == nil {
		return ebiten.CursorPosition()
	} else {
		return r.GetCursorPositionFunc()
	}
}

func (r *UIContext) PreUpdate() {
	r.cursorShape = ebiten.CursorShapeDefault
	r.isHovered = false
}

func (r *UIContext) Update() {
	ebiten.SetCursorShape(r.cursorShape)
}

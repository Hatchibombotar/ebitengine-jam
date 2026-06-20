package superui

import "github.com/hajimehoshi/ebiten/v2"

type GenericWidget interface {
	Draw(screen *ebiten.Image, root *UIContainer)
	Update()
	UpdateInput(*UIContainer)

	UpdatePreStage()
	UpdateSizeWidthFitPass()
	UpdateSizeHeightFitPass()
	UpdateSizeWrapWidth()
	UpdateSizeWidthGrowPass()
	UpdateSizeHeightGrowPass()
	UpdatePosition(x, y int)

	SetParent(parent GenericWidget)
	GetResultWidth() int
	GetResultHeight() int
	GetResultX() int
	GetResultY() int
	CanGrowWidth() bool
	CanGrowHeight() bool
	GetMinWidth() int
	GetMinHeight() int

	SetResultWidth(int)
	SetResultHeight(int)

	String() string

	IsFocusable() bool

	GetPositionMode() PositionMode
}

type FocusableWidget interface {
	SetFocused(bool)
}

package superui

import "github.com/hajimehoshi/ebiten/v2"

type ConditionalDisplay struct {
	ops   *ConditionalDisplayOps
	child GenericWidget
}

type ConditionalDisplayOps struct {
	ShouldShow func() bool
}

func NewConditionalDisplay(op *ConditionalDisplayOps, element GenericWidget) *ConditionalDisplay {
	w := &ConditionalDisplay{
		ops:   op,
		child: element,
	}

	return w
}

func (c *ConditionalDisplay) ShouldShow() bool {
	if ebiten.IsKeyPressed(ebiten.KeyF9) {
		return true
	}
	return c.ops.ShouldShow()
}

func (c *ConditionalDisplay) Draw(screen *ebiten.Image, root *UIContainer) {
	if !c.ShouldShow() {
		return
	}
	c.child.Draw(screen, root)
}

// TODO: instead of returning empty values,
// make a new function under GenericWidget that checks if the item is part of the layout.
// this will be also useful for other things.

func (c *ConditionalDisplay) CanGrowHeight() bool {
	if !c.ShouldShow() {
		return false
	}
	return c.child.CanGrowHeight()
}
func (c *ConditionalDisplay) CanGrowWidth() bool {
	if !c.ShouldShow() {
		return false
	}
	return c.child.CanGrowWidth()
}
func (c *ConditionalDisplay) GetMinHeight() int {
	if !c.ShouldShow() {
		return 0
	}
	return c.child.GetMinHeight()
}
func (c *ConditionalDisplay) GetMinWidth() int {
	if !c.ShouldShow() {
		return 0
	}
	return c.child.GetMinWidth()
}
func (c *ConditionalDisplay) GetResultHeight() int {
	if !c.ShouldShow() {
		return 0
	}
	return c.child.GetResultHeight()
}
func (c *ConditionalDisplay) GetResultWidth() int {
	if !c.ShouldShow() {
		return 0
	}
	return c.child.GetResultWidth()
}
func (c *ConditionalDisplay) GetResultX() int {
	if !c.ShouldShow() {
		return 0
	}
	return c.child.GetResultX()
}
func (c *ConditionalDisplay) GetResultY() int {
	if !c.ShouldShow() {
		return 0
	}
	return c.child.GetResultY()
}
func (c *ConditionalDisplay) SetParent(parent GenericWidget) {
	if !c.ShouldShow() {
		return
	}
	c.child.SetParent(parent)
}
func (c *ConditionalDisplay) SetResultHeight(h int) {
	if !c.ShouldShow() {
		return
	}
	c.child.SetResultHeight(h)
}
func (c *ConditionalDisplay) SetResultWidth(w int) {
	if !c.ShouldShow() {
		return
	}
	c.child.SetResultHeight(w)
}
func (c *ConditionalDisplay) String() string {
	if !c.ShouldShow() {
		return "<ConditionalDisplay/>"
	}
	return c.child.String()
}
func (c *ConditionalDisplay) Update() {
	if !c.ShouldShow() {
		return
	}
	c.child.Update()
}
func (c *ConditionalDisplay) UpdatePreStage() {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdatePreStage()
}
func (c *ConditionalDisplay) UpdateInput(r *UIContainer) {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdateInput(r)
}
func (c *ConditionalDisplay) UpdatePosition(x int, y int) {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdatePosition(x, y)
}
func (c *ConditionalDisplay) UpdateSizeHeightFitPass() {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdateSizeHeightFitPass()
}
func (c *ConditionalDisplay) UpdateSizeHeightGrowPass() {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdateSizeHeightGrowPass()
}
func (c *ConditionalDisplay) UpdateSizeWidthFitPass() {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdateSizeWidthFitPass()
}
func (c *ConditionalDisplay) UpdateSizeWidthGrowPass() {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdateSizeWidthGrowPass()
}
func (c *ConditionalDisplay) UpdateSizeWrapWidth() {
	if !c.ShouldShow() {
		return
	}
	c.child.UpdateSizeWrapWidth()
}

func (w *ConditionalDisplay) IsFocusable() bool {
	return false
}

func (w *ConditionalDisplay) GetPositionMode() PositionMode {
	return PositionAuto
}

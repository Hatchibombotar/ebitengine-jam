package superui

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SizeMode = int

const (
	SizeFit SizeMode = iota
	SizeFixed
	SizeGrow
)

type LayoutDirection int

const (
	LayoutColumn LayoutDirection = iota
	LayoutRow
)

type LayoutType int

const (
	LayoutFlex LayoutType = iota
	LayoutGrid
)

type AlignMode int

const (
	AlignStart AlignMode = iota
	AlignCenter
	AlignEnd
)

type PositionMode = int

const (
	PositionAuto PositionMode = iota
	PositionFixed
	PositionRelative
)

type Box struct {
	Op *BoxWidgetOps

	Children []GenericWidget
	parent   GenericWidget

	eventHandler *UIContainer

	resultX, resultY          int
	resultWidth, resultHeight int
	minWidth, minHeight       int
}

type BoxWidgetOps struct {
	WidthMode  SizeMode
	Width      int
	HeightMode SizeMode
	Height     int

	Padding Spacing

	// DrawBackgroundColour color.Color
	Draw func(w Box, screen *ebiten.Image)

	LayoutType      LayoutType
	LayoutDirection LayoutDirection
	Gap             int

	// Rows or column count for LayoutType == Grid
	Rows, Columns int

	// align along the horizonal and vertical axies
	AlignHorizontal, AlignVertical AlignMode

	OnDraw func(screen *ebiten.Image, widget GenericWidget, root *UIContainer)

	OnInputUpdate func(widget GenericWidget, eventRoot *UIContainer)

	IsFocusable bool
	CursorShape ebiten.CursorShapeType

	PositionMode PositionMode

	// X and Y position of widget
	//
	// If Op.PositionMode is set to PositionFixed, relative to 0, 0
	//
	// If Op.PositionMode is set to PositionRelative, relative to Top Left & Right of Parent
	//
	// Otherwise, no other effects
	X, Y int
}

func NewBoxWidget(ops *BoxWidgetOps, children ...GenericWidget) *Box {
	w := &Box{
		Op:       ops,
		Children: children,
	}
	for _, child := range children {
		child.SetParent(w)
	}
	return w
}

func (w *Box) AddChild(child GenericWidget) {
	w.Children = append(w.Children, child)
}

func (w *Box) Draw(screen *ebiten.Image, root *UIContainer) {
	if w.Op.OnDraw != nil {
		w.Op.OnDraw(screen, w, root)
	}

	for _, child := range w.Children {
		child.Draw(screen, root)
	}
}

func (w *Box) Update() {
	w.UpdatePreStage()
	w.UpdateSizeWidthFitPass()
	w.UpdateSizeWidthGrowPass()
	w.UpdateSizeWrapWidth()
	w.UpdateSizeHeightFitPass()
	w.UpdateSizeHeightGrowPass()
	w.UpdatePosition(0, 0)
	// w.PropagateUpdates()
}

func (w *Box) UpdatePreStage() {
	for _, child := range w.Children {
		child.UpdatePreStage()
	}
}

func (w *Box) UpdateSizeWidthFitPass() {
	if w.Op.HeightMode == SizeFixed {
		w.resultWidth = w.Op.Width
		w.minWidth = w.Op.Width
		for _, child := range w.Children {
			child.UpdateSizeWidthFitPass()
		}
		return
	}
	if w.Op.LayoutType == LayoutGrid {
		childrenPlaced := 0
		width := 0
		minWidth := 0
		rows := int(math.Ceil(float64(len(w.Children)) / float64(w.Op.Columns)))
		for range rows {
			rowTotalWidth := 0
			rowTotalMinWidth := 0

			for range w.Op.Columns {
				if childrenPlaced >= len(w.Children) {
					break
				}

				child := w.Children[childrenPlaced]
				child.UpdateSizeWidthFitPass()
				// TODO: Account for fixed & relative position mode for children

				rowTotalWidth += child.GetResultWidth()
				rowTotalMinWidth += child.GetMinWidth()
				childrenPlaced += 1
			}
			if rowTotalWidth > width {
				width = rowTotalWidth
			}
			if rowTotalMinWidth > minWidth {
				minWidth = rowTotalMinWidth
			}
		}
		width += w.Op.Padding.Left + w.Op.Padding.Right + (w.Op.Columns-1)*w.Op.Gap
		minWidth += w.Op.Padding.Left + w.Op.Padding.Right + (w.Op.Columns-1)*w.Op.Gap

		switch w.Op.WidthMode {
		case SizeFit, SizeGrow:
			w.resultWidth = width
			w.minWidth = minWidth
		case SizeFixed:
			w.resultWidth = w.Op.Width
			w.minWidth = w.Op.Width
		default:
			panic("Unexpected sizing mode")
		}

	} else if w.Op.LayoutDirection == LayoutRow {
		minWidth := w.Op.Padding.Left + w.Op.Padding.Right + (len(w.Children)-1)*w.Op.Gap
		width := minWidth

		for _, child := range w.Children {
			child.UpdateSizeWidthFitPass()
			width += child.GetResultWidth()
			minWidth += child.GetMinWidth()
		}

		switch w.Op.WidthMode {
		case SizeFit, SizeGrow:
			w.resultWidth = width
			w.minWidth = minWidth
		case SizeFixed:
			w.resultWidth = w.Op.Width
			w.minWidth = w.Op.Width
		default:
			panic("Unexpected sizing mode")
		}

	} else if w.Op.LayoutDirection == LayoutColumn {
		width := 0
		minWidth := 0

		for _, child := range w.Children {
			child.UpdateSizeWidthFitPass()

			width = max(width, child.GetResultWidth())
			minWidth = max(minWidth, child.GetMinWidth())
		}
		width += w.Op.Padding.Left + w.Op.Padding.Right
		minWidth += w.Op.Padding.Left + w.Op.Padding.Right

		switch w.Op.WidthMode {
		case SizeFit, SizeGrow:
			w.resultWidth = width
			w.minWidth = minWidth
		case SizeFixed:
			w.resultWidth = w.Op.Width
			w.minWidth = w.Op.Width
		default:
			panic("Unexpected sizing mode")
		}
	}
}

func (w *Box) UpdateSizeHeightFitPass() {
	if w.Op.HeightMode == SizeFixed {
		w.resultHeight = w.Op.Height
		w.minHeight = w.Op.Height
		for _, child := range w.Children {
			child.UpdateSizeHeightFitPass()
		}
		return
	}
	if w.Op.LayoutType == LayoutGrid {
		childrenPlaced := 0
		height := 0
		minHeight := 0
		rows := int(math.Ceil(float64(len(w.Children)) / float64(w.Op.Columns)))
		for range rows {
			rowMaxHeight := 0
			rowMaxMinHeight := 0

			for range w.Op.Columns {
				if childrenPlaced >= len(w.Children) {
					break
				}

				child := w.Children[childrenPlaced]
				child.UpdateSizeHeightFitPass()
				// TODO: Account for fixed & relative position mode for children

				rowMaxHeight = max(rowMaxHeight, child.GetResultHeight())
				rowMaxMinHeight = max(rowMaxMinHeight, child.GetMinHeight())
				childrenPlaced += 1
			}
			height += rowMaxHeight
			minHeight += rowMaxMinHeight
		}
		height += w.Op.Padding.Top + w.Op.Padding.Bottom + (rows-1)*w.Op.Gap
		minHeight += w.Op.Padding.Top + w.Op.Padding.Bottom + (rows-1)*w.Op.Gap

		w.resultHeight = height
		w.minHeight = minHeight

	} else if w.Op.LayoutDirection == LayoutRow {
		minHeight := 0
		height := 0

		for _, child := range w.Children {
			child.UpdateSizeHeightFitPass()

			minHeight = max(minHeight, child.GetMinHeight())
			height = max(height, child.GetResultHeight())
		}
		height += w.Op.Padding.Top + w.Op.Padding.Bottom
		minHeight += w.Op.Padding.Top + w.Op.Padding.Bottom

		w.resultHeight = height
		w.minHeight = minHeight

	} else if w.Op.LayoutDirection == LayoutColumn {
		minHeight := w.Op.Padding.Top + w.Op.Padding.Bottom + (len(w.Children)-1)*w.Op.Gap
		height := minHeight

		for _, child := range w.Children {
			child.UpdateSizeHeightFitPass()
			height += child.GetResultHeight()
			minHeight += child.GetMinHeight()
		}

		w.resultHeight = height
		w.minHeight = minHeight
	}
}

func (parent *Box) UpdateSizeWidthGrowPass() {
	remainingWidth := parent.GetResultWidth()
	if parent.Op.LayoutType == LayoutGrid {
		if parent.Op.LayoutDirection == LayoutColumn {
			if parent.Op.Columns == 0 {
				panic("Columns set to 0.")
			}
			remainingWidth /= parent.Op.Columns
			remainingWidth -= parent.Op.Gap * (parent.Op.Columns - 1)
		} else {
			panic("Unexpected Layout Direction")
		}
	}
	remainingWidth -= parent.Op.Padding.Left + parent.Op.Padding.Right

	growable := make([]GenericWidget, 0)
	shrinkable := make([]GenericWidget, 0)
	for _, child := range parent.Children {
		remainingWidth -= child.GetResultWidth()
		if child.CanGrowWidth() {
			growable = append(growable, child)
		}
		if child.GetMinWidth() != child.GetResultWidth() {
			shrinkable = append(shrinkable, child)
		}
	}
	remainingWidth -= (len(parent.Children) - 1) * parent.Op.Gap

	// grow elements
	for remainingWidth > 0 {
		if len(growable) == 0 {
			break
		}
		smallest := growable[0].GetResultWidth()
		secondSmallest := 99999999
		widthToAdd := remainingWidth
		for _, child := range growable {
			if child.GetResultWidth() < smallest {
				secondSmallest = smallest
				smallest = child.GetResultWidth()
			}
			if child.GetResultWidth() > smallest {
				secondSmallest = min(secondSmallest, child.GetResultWidth())
				widthToAdd = secondSmallest - smallest
			}
		}

		widthToAdd = min(widthToAdd, remainingWidth/len(growable))

		// temp
		if widthToAdd == 0 {
			widthToAdd = remainingWidth
		}

		for _, child := range growable {
			if child.GetResultWidth() == smallest {
				child.SetResultWidth(child.GetResultWidth() + widthToAdd)
				remainingWidth -= widthToAdd
			}
		}
	}

	// shrink elements
	for remainingWidth < 0 {
		if len(shrinkable) == 0 {
			break
		}
		largest := shrinkable[0].GetResultWidth()
		secondLargest := 99999999
		widthToAdd := remainingWidth
		for _, child := range shrinkable {
			if child.GetResultWidth() > largest {
				secondLargest = largest
				largest = child.GetResultWidth()
			}
			if child.GetResultWidth() < largest {
				secondLargest = max(secondLargest, child.GetResultWidth())
				widthToAdd = secondLargest - largest
			}
		}

		widthToAdd = min(widthToAdd, remainingWidth/len(shrinkable))

		// temp
		if widthToAdd == 0 {
			widthToAdd = remainingWidth
		}

		newShrinkable := make([]GenericWidget, 0)
		for _, child := range shrinkable {
			previousWidth := child.GetResultWidth()
			if child.GetResultWidth() == largest {
				child.SetResultWidth(child.GetResultWidth() + widthToAdd)
				if child.GetResultWidth() < child.GetMinWidth() {
					child.SetResultWidth(child.GetMinWidth())
					continue
				}
				remainingWidth -= child.GetResultWidth() - previousWidth
				if child.GetResultWidth() == child.GetMinWidth() {
					continue
				}
			}
			newShrinkable = append(newShrinkable, child)
		}
		shrinkable = newShrinkable
	}

	for _, child := range parent.Children {
		child.UpdateSizeWidthGrowPass()
	}
}

// TODO: this doesn't work, make it like width
func (parent *Box) UpdateSizeHeightGrowPass() {
	remainingHeight := parent.GetResultHeight()
	remainingHeight -= parent.Op.Padding.Top + parent.Op.Padding.Bottom

	for _, child := range parent.Children {
		if !child.CanGrowHeight() {
			remainingHeight -= child.GetResultHeight()
		}
	}

	for _, child := range parent.Children {
		if child.CanGrowHeight() {
			child.SetResultHeight(remainingHeight)
		}
	}

	for _, child := range parent.Children {
		child.UpdateSizeHeightGrowPass()
	}
}

func (w *Box) UpdateSizeWrapWidth() {
	for _, child := range w.Children {
		child.UpdateSizeWrapWidth()
	}
}

func (w *Box) UpdatePosition(x, y int) {
	switch w.Op.PositionMode {
	case PositionFixed:
		w.resultX = w.Op.X
		w.resultY = w.Op.Y
	case PositionRelative:
		w.resultX = x + w.Op.X
		w.resultY = y + w.Op.Y
	default:
		w.resultX = x
		w.resultY = y
	}

	if w.Op.LayoutType == LayoutGrid {
		startX := w.resultX + w.Op.Padding.Left
		currentY := w.resultY + w.Op.Padding.Top
		childrenPlaced := 0
		rows := int(math.Ceil(float64(len(w.Children)) / float64(w.Op.Columns)))
		for range rows {
			rowMaxHeight := 0

			currentX := startX
			for range w.Op.Columns {
				if childrenPlaced >= len(w.Children) {
					break
				}

				child := w.Children[childrenPlaced]
				child.UpdatePosition(currentX, currentY)
				// TODO: Account for fixed & relative position mode for children

				rowMaxHeight = max(rowMaxHeight, child.GetResultHeight())
				childrenPlaced += 1

				currentX += child.GetResultWidth() + w.Op.Gap

				// TODO: amount shifted should not be based on item width.
				// currentX += (w.resultWidth / w.Op.Columns) + w.Op.Gap
			}

			currentY += rowMaxHeight + w.Op.Gap
		}
	} else if w.Op.LayoutDirection == LayoutRow {
		currentX := 0
		if w.Op.AlignHorizontal == AlignStart {
			currentX = w.resultX + w.Op.Padding.Left
		} else {
			panic("not implemented")
		}

		for _, child := range w.Children {
			if child.GetPositionMode() == PositionFixed || child.GetPositionMode() == PositionRelative {
				child.UpdatePosition(w.resultX, w.resultY)
				continue
			}
			currentY := 0
			if w.Op.AlignVertical == AlignStart {
				currentY = w.resultY + w.Op.Padding.Top
			} else {
				panic("not implemented")
			}

			child.UpdatePosition(currentX, currentY)

			currentX += child.GetResultWidth()
			currentX += w.Op.Gap
		}
	} else if w.Op.LayoutDirection == LayoutColumn {
		currentY := 0
		if w.Op.AlignVertical == AlignStart {
			currentY = w.resultY + w.Op.Padding.Top
		} else {
			panic("")
		}

		for _, child := range w.Children {
			if child.GetPositionMode() == PositionFixed || child.GetPositionMode() == PositionRelative {
				child.UpdatePosition(w.resultX, w.resultY)
				continue
			}
			currentX := 0
			switch w.Op.AlignHorizontal {
			case AlignStart:
				currentX = w.resultX + w.Op.Padding.Left
			case AlignCenter:
				currentX = w.resultX + (w.GetResultWidth()-child.GetResultWidth())/2
			default:
				panic("")
			}
			child.UpdatePosition(currentX, currentY)

			currentY += child.GetResultHeight()
			currentY += w.Op.Gap
		}
	}
}

type Spacing struct {
	Top, Right, Bottom, Left int
}

func (s Spacing) WithAll(value int) Spacing {
	s.Top = value
	s.Right = value
	s.Bottom = value
	s.Left = value

	return s
}

// generic things
func (w *Box) SetParent(parent GenericWidget) {
	w.parent = parent
}

func (w *Box) GetResultWidth() int {
	return w.resultWidth
}

func (w *Box) GetResultHeight() int {
	return w.resultHeight
}

func (w *Box) GetResultX() int {
	return w.resultX
}

func (w *Box) GetResultY() int {
	return w.resultY
}

func (w *Box) CanGrowWidth() bool {
	return w.Op.WidthMode == SizeGrow
}
func (w *Box) CanGrowHeight() bool {
	return w.Op.HeightMode == SizeGrow
}

func (w *Box) GetMinWidth() int {
	return w.minWidth
}
func (w *Box) GetMinHeight() int {
	return w.minHeight
}

func (w *Box) SetResultWidth(value int) {
	w.resultWidth = value
}
func (w *Box) SetResultHeight(value int) {
	w.resultHeight = value
}

func (w *Box) String() string {
	return w.string("")
}

func (w *Box) string(padding string) string {
	str := padding + fmt.Sprint(
		"<Box",
		" width(", w.Op.WidthMode, ")=", w.resultWidth,
		", height(", w.Op.HeightMode, ")=", w.resultHeight,
		", padding=(", w.Op.Padding.Left, ",", w.Op.Padding.Top, ",", w.Op.Padding.Right, ",", w.Op.Padding.Bottom, ")",
		", x=", w.resultX,
		", y=", w.resultY,
		", hasEventHandler=", w.eventHandler != nil,
		">") + "\n"
	for _, child := range w.Children {
		switch child := child.(type) {
		case *Box:
			str += child.string(padding+"  ") + "\n"
		default:
			str += padding + "  " + child.String() + "\n"
		}
	}
	str += padding + "</Box>"
	return str
}

func (w *Box) UpdateInput(root *UIContainer) {
	if w.Op.IsFocusable && root.IsHovered(w) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		root.SetFocusOn(w)
	}

	if w.Op.CursorShape != 0 && root.IsHovered(w) {
		root.SetCursorShape(w.Op.CursorShape)
	}

	if w.Op.OnInputUpdate != nil {
		w.Op.OnInputUpdate(w, root)
	}

	if root.IsHovered(w) {
		root.hoveredWidget = w
		root.hoveredWidgets = append(root.hoveredWidgets, w)
	}

	for _, child := range w.Children {
		child.UpdateInput(root)
	}
}

func (w *Box) IsFocusable() bool {
	return w.Op.IsFocusable
}

func (w *Box) GetPositionMode() PositionMode {
	return PositionAuto
}

package superui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	text   string
	Op     *TextWidgetOps
	parent GenericWidget

	resultX, resultY          int
	resultWidth, resultHeight int
	minWidth, minHeight       int

	// resultFontFace *text.GoTextFace

	lines []string
}

type WrapBehaviour int

const (
	WrapText WrapBehaviour = iota
	NoWrap
)

type TextWidgetOps struct {
	Color         color.Color
	Face          text.Face
	WrapBehaviour WrapBehaviour
	// If set, will be called on every update call. If the text has changed, it will update.
	GetDynamicText func() string

	TextAlign AlignMode
}

func NewTextWidget(ops *TextWidgetOps, text string) *Text {
	w := &Text{
		Op:   ops,
		text: text,
	}
	return w
}

func (w *Text) Draw(screen *ebiten.Image, root *UIContainer) {
	face := w.Op.Face

	for i, line := range w.lines {
		op := &text.DrawOptions{}

		fontSize := face.Metrics().HAscent

		if w.Op.TextAlign == AlignCenter {
			lineWidth := textWidth(face, line)
			op.GeoM.Translate(float64(w.resultWidth)/2-float64(lineWidth)/2, 0)
		}

		op.GeoM.Translate(float64(w.resultX), float64(w.resultY+(i*int(fontSize))))

		if w.Op.Color == nil {
			panic("Colour not defined for text")
		}
		op.ColorScale.ScaleWithColor(w.Op.Color)

		text.Draw(screen, line, face, op)
	}
}

func (w *Text) Update() {
	panic("Text Update() function should never be called.")
}
func (w *Text) UpdatePreStage() {
	if w.Op.GetDynamicText != nil {
		newText := w.Op.GetDynamicText()
		if w.text != newText {
			w.text = newText
			// recompute stuff
		}
	}
}
func (w *Text) UpdateSizeWidthFitPass() {
	minWidth, maxWidth := w.calculateTextBounds()
	w.minWidth = minWidth
	w.resultWidth = maxWidth

	w.minHeight = 0
	w.resultHeight = 0
}
func (w *Text) UpdateSizeHeightFitPass()  {}
func (w *Text) UpdateSizeWidthGrowPass()  {}
func (w *Text) UpdateSizeHeightGrowPass() {}

func (w *Text) UpdateSizeWrapWidth() {
	w.calculateLineWrappings()
}

func (w *Text) UpdatePosition(x, y int) {
	w.resultX = x
	w.resultY = y
}

func (w *Text) calculateTextBounds() (minWidth int, maxWidth int) {
	if w.Op.Face == nil {
		panic("font face not definied")
	}
	if w.Op.WrapBehaviour == NoWrap {
		width := textWidth(w.Op.Face, w.text)
		return width, width
	}
	maxTextWidth := 0
	totalWidth := 0

	startIndex := 0
	for endIndex, char := range w.text {
		if char == ' ' || char == '\n' || char == '\t' {
			str := string(w.text[startIndex:endIndex])
			wordWidth := textWidth(w.Op.Face, str)

			totalWidth += wordWidth
			if char == ' ' || char == '\t' {
				totalWidth += textWidth(w.Op.Face, string(char))
			}

			maxTextWidth = max(maxTextWidth, wordWidth)
			startIndex = endIndex
		}
	}
	str := string(w.text[startIndex:len(w.text)])
	wordWidth := textWidth(w.Op.Face, str)

	totalWidth += wordWidth

	maxTextWidth = max(maxTextWidth, wordWidth)

	return maxTextWidth, totalWidth
}

func textWidth(fontFace text.Face, str string) int {
	w, _ := text.Measure(str, fontFace, 0)
	return int(w)
	// return int(text.Advance(str, fontFace))
}

func (w *Text) calculateLineWrappings() {
	face := w.Op.Face
	lines := make([]string, 0)

	if w.Op.WrapBehaviour == NoWrap {
		w.resultWidth = textWidth(face, w.text)
		fontSize := face.Metrics().HAscent
		w.resultHeight = int(fontSize)
		lines = append(lines, w.text)
		w.lines = lines
		return
	}
	var endIdx, p int
	maxWidth := 0

	for endIdx < len(w.text) {
		wi := 0
		endIdx = p
		startIdx := endIdx
		for endIdx < len(w.text) && w.text[endIdx] != '\n' {
			word := p
			for p < len(w.text) && w.text[p] != ' ' && w.text[p] != '\n' {
				p++
			}
			if wi > maxWidth {
				maxWidth = wi
			}
			wi += textWidth(face, w.text[word:p])
			if wi > w.resultWidth && endIdx != startIdx {
				break
			}
			if p < len(w.text) {
				wi += textWidth(face, string(w.text[p]))
			}
			endIdx = p
			p++
		}

		lines = append(lines, w.text[startIdx:endIdx])
		p = endIdx + 1
	}

	w.lines = lines
	// if maxWidth < w.resultWidth {
	// 	w.InternalWidget.width = maxWidth
	// } else {
	// 	w.InternalWidget.width = w.resultWidth
	// }
	fontSize := face.Metrics().HAscent
	w.resultHeight = len(lines) * int(fontSize)
	if len(lines) > 1 {
		// previously removed for some reason
		// reason: when added it breaks shit
		w.resultWidth = maxWidth
	}
}

// generic stuff
func (w *Text) SetParent(parent GenericWidget) {
	w.parent = parent
}

func (w *Text) GetResultWidth() int {
	return w.resultWidth
}

func (w *Text) GetResultHeight() int {
	return w.resultHeight
}

func (w *Text) GetResultX() int {
	return w.resultX
}

func (w *Text) GetResultY() int {
	return w.resultY
}

func (w *Text) CanGrowWidth() bool {
	return false
}
func (w *Text) CanGrowHeight() bool {
	return false
}

func (w *Text) GetMinHeight() int {
	return 0
}
func (w *Text) GetMinWidth() int {
	// if w.Op.WrapBehaviour == NoWrap {
	width, _ := w.calculateTextBounds()
	return width
	// }
	// return 0
}

func (w *Text) SetResultWidth(value int) {
	w.resultWidth = value
}
func (w *Text) SetResultHeight(value int) {
	w.resultHeight = value
}

func (w *Text) String() string {
	return fmt.Sprint(
		"<Text width=", w.resultWidth, " height=", w.resultHeight, " lines=", len(w.lines), " size=", w.Op.Face.Metrics().HAscent, "/>",
	)
}

func (w *Text) UpdateInput(root *UIContainer) {
	if root.IsHovered(w) {
		root.hoveredWidget = w
	}
}

func (w *Text) SetText(newText string) {
	w.text = newText
}

func (w *Text) IsFocusable() bool {
	return false
}

func (w *Text) GetPositionMode() PositionMode {
	return PositionAuto
}

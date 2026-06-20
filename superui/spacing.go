package superui

func NewSpacingWidget(x, y int) *Box {
	return NewBoxWidget(
		&BoxWidgetOps{
			WidthMode:  SizeFixed,
			HeightMode: SizeFixed,
			Width:      x,
			Height:     y,
		},
	)
}

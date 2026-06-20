package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var fontData = ReadFileBytes("assets/fonts/kenney-mini.ttf")
var fontFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}
	fontFaceSource = s
}

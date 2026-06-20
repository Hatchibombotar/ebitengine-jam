package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed assets/**
var emb embed.FS

func LoadImageFromPath(path string) *ebiten.Image {
	file, err := emb.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	sheet := ebiten.NewImageFromImage(img)

	return sheet
}

func ReadFileBytes(path string) []byte {
	bytes, err := emb.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return bytes
}

const SAMPLE_RATE = 48000

func ReadOggBytesFromPath(path string) []byte {
	data, err := emb.ReadFile(path)

	if err != nil {
		panic(err)
	}

	s, err := vorbis.DecodeWithSampleRate(SAMPLE_RATE, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(s)
	if err != nil {
		panic(err)
	}

	return b
}

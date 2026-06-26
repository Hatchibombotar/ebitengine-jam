package main

import "github.com/hajimehoshi/ebiten/v2"

var crtshadersrc = ReadFileBytes("assets/crtshader.kage")

var crtshader *ebiten.Shader

func init() {

	s, err := ebiten.NewShader(crtshadersrc)
	if err != nil {
		panic(err)
	}
	crtshader = s
}

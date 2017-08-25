package ui

import "github.com/veandco/go-sdl2/sdl"

func NewColor(rgba uint32) sdl.Color {
	c := sdl.Color{}
	c.R = uint8((rgba & 0xff000000) >> 24)
	c.G = uint8((rgba & 0x00ff0000) >> 16)
	c.B = uint8((rgba & 0x0000ff00) >> 8)
	c.A = uint8((rgba & 0x000000ff) >> 0)
	return c
}

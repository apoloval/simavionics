package ui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type ValueRenderer struct {
	value    interface{}
	renderer *sdl.Renderer
	font     *ttf.Font
	color    sdl.Color
	texture  *sdl.Texture
	w        int32
	h        int32
}

func NewValueRenderer(renderer *sdl.Renderer, font *ttf.Font, color sdl.Color) *ValueRenderer {
	value := &ValueRenderer{renderer: renderer, font: font, color: color}
	return value
}

func (vr *ValueRenderer) SetValue(v interface{}) {
	if v != vr.value {
		vr.value = v
		if vr.texture != nil {
			vr.texture.Destroy()
		}
		text := fmt.Sprintf("%v", vr.value)
		if len(text) == 0 {
			text = " "
		}
		texture, w, h, err := vr.renderText(text)
		if err != nil {
			panic(err)
		}
		vr.texture = texture
		vr.w = w
		vr.h = h
	}
}

func (vr *ValueRenderer) Render(x int32, y int32) {
	if vr.texture != nil {
		vr.renderer.Copy(vr.texture, nil, &sdl.Rect{x, y, vr.w, vr.h})
	}
}

func (vr *ValueRenderer) renderText(text string) (*sdl.Texture, int32, int32, error) {
	surface, err := vr.font.RenderUTF8_Blended(text, vr.color)
	if err != nil {
		return nil, 0, 0, err
	}
	defer surface.Free()

	texture, err := vr.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, 0, 0, err
	}
	return texture, surface.W, surface.H, nil
}

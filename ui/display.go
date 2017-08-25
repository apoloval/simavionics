package ui

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Display struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func NewDisplay(title string) (*Display, error) {
	var err error
	d := &Display{}

	ttf.Init()
	if !sdl.SetHintWithPriority(sdl.HINT_RENDER_SCALE_QUALITY, "best", sdl.HINT_OVERRIDE) {
		log.Printf("[ui.display] Cannot set HINT_RENDER_SCALE_QUALITY value")
	}

	d.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}

	d.renderer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, err
	}

	return d, nil
}
func (d *Display) Renderer() *sdl.Renderer {
	return d.renderer
}

func (d *Display) Destroy() {
	d.renderer.Destroy()
	d.window.Destroy()
}

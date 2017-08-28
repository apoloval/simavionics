package ui

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	DisplayWidth  = 1024
	DisplayHeight = 768
)

type Display struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	positioner Positioner
}

func NewDisplay(title string) (*Display, error) {
	var err error
	d := &Display{positioner: NewPositioner(DisplayWidth, DisplayHeight, 1.0, 1.0)}

	ttf.Init()
	if !sdl.SetHintWithPriority(sdl.HINT_RENDER_SCALE_QUALITY, "best", sdl.HINT_OVERRIDE) {
		log.Printf("[ui.display] Cannot set HINT_RENDER_SCALE_QUALITY value")
	}

	d.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		DisplayWidth, DisplayHeight,
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

func (d *Display) Positioner() *Positioner {
	return &d.positioner
}

func (d *Display) Destroy() {
	d.renderer.Destroy()
	d.window.Destroy()
}

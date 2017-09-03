package ui

import (
	"github.com/op/go-logging"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var log = logging.MustGetLogger("display")

type Display struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	positioner Positioner
}

func NewDisplay(title string, width, height uint) (*Display, error) {
	var err error
	d := &Display{
		positioner: NewPositioner(int32(width), int32(height), 1.0, 1.0),
	}

	ttf.Init()
	if !sdl.SetHintWithPriority(sdl.HINT_RENDER_SCALE_QUALITY, "best", sdl.HINT_OVERRIDE) {
		log.Error("Cannot set HINT_RENDER_SCALE_QUALITY value")
	}

	d.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(width), int(height),
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

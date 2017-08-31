package main

import (
	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/ui"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const greenColor = 0x00C000FF

type apuPage struct {
	n1  float64
	egt float64

	n1Chan       <-chan simavionics.EventValue
	egtChan      <-chan simavionics.EventValue
	flapOpenChan <-chan simavionics.EventValue

	backgroundTexture *sdl.Texture
	pointerTexture    *sdl.Texture

	n1Text       *ui.ValueRenderer
	egtText      *ui.ValueRenderer
	flapOpenText *ui.ValueRenderer
}

func newAPUPage(bus simavionics.EventBus, display *ui.Display) (*apuPage, error) {
	renderer := display.Renderer()
	positioner := display.Positioner()

	fontSize := positioner.Map(ui.RectF{H: 0.0375})
	font, err := ttf.OpenFont("assets/fonts/Carlito-Regular.ttf", int(fontSize.H))
	if err != nil {
		return nil, err
	}

	page := &apuPage{
		n1Chan:       bus.Subscribe(apu.EventEngineN1),
		egtChan:      bus.Subscribe(apu.EventEngineEGT),
		flapOpenChan: bus.Subscribe(apu.EventFlap),
		n1Text:       ui.NewValueRenderer(renderer, font, ui.NewColor(greenColor)),
		egtText:      ui.NewValueRenderer(renderer, font, ui.NewColor(greenColor)),
		flapOpenText: ui.NewValueRenderer(renderer, font, ui.NewColor(greenColor)),
	}

	page.backgroundTexture, err = img.LoadTexture(renderer, "assets/ecam-apu-background.png")
	if err != nil {
		return nil, err
	}

	page.pointerTexture, err = img.LoadTexture(renderer, "assets/apu-gauge-pointer.png")
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (p *apuPage) processEvents() {
	select {
	case v := <-p.n1Chan:
		p.n1 = v.Float64()
		p.n1Text.SetValue(int(p.n1))
	case v := <-p.egtChan:
		p.egt = v.Float64()
		p.egtText.SetValue(int(p.egt))
	case v := <-p.flapOpenChan:
		if v.Bool() {
			p.flapOpenText.SetValue("FLAP OPEN")
		} else {
			p.flapOpenText.SetValue("")
		}
	default:
	}
}

func (p *apuPage) render(display *ui.Display) {
	renderer := display.Renderer()
	positioner := display.Positioner()

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Copy(p.backgroundTexture, nil, nil)

	n1PointerRect := positioner.Map(ui.RectF{X: 0.365625, Y: 0.485416, W: 0.007812, H: 0.097917})
	egtPointerRect := positioner.Map(ui.RectF{X: 0.365625, Y: 0.672916, W: 0.007812, H: 0.097917})

	renderer.CopyEx(p.pointerTexture, nil, &n1PointerRect, (p.n1*170.0/100.0)+65.0, &sdl.Point{2, 1}, sdl.FLIP_NONE)
	renderer.CopyEx(p.pointerTexture, nil, &egtPointerRect, (p.egt*125.0/1000.0)+65.0, &sdl.Point{2, 1}, sdl.FLIP_NONE)

	n1TextRect := positioner.Map(ui.RectF{X: 0.375, Y: 0.489583})
	egtTextRect := positioner.Map(ui.RectF{X: 0.375, Y: 0.677083})
	flapOpenTextRect := positioner.Map(ui.RectF{X: 0.578125, Y: 0.5625})

	p.n1Text.Render(n1TextRect.X, n1TextRect.Y)
	p.egtText.Render(egtTextRect.X, egtTextRect.Y)
	p.flapOpenText.Render(flapOpenTextRect.X, flapOpenTextRect.Y)

	renderer.Present()
}

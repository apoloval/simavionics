package main

import (
	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/ui"
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

func newAPUPage(bus simavionics.EventBus, renderer *sdl.Renderer) (*apuPage, error) {
	font, err := ttf.OpenFont("assets/fonts/Carlito-Regular.ttf", 18)
	if err != nil {
		return nil, err
	}

	page := &apuPage{
		n1Chan:       bus.Subscribe(apu.EngineStateN1),
		egtChan:      bus.Subscribe(apu.EngineStateEGT),
		flapOpenChan: bus.Subscribe(apu.StatusFlapOpen),
		n1Text:       ui.NewValueRenderer(renderer, font, ui.NewColor(greenColor)),
		egtText:      ui.NewValueRenderer(renderer, font, ui.NewColor(greenColor)),
		flapOpenText: ui.NewValueRenderer(renderer, font, ui.NewColor(greenColor)),
	}

	page.backgroundTexture, err = ui.LoadTextureFromBMP(renderer, "assets/apu-background.bmp")
	if err != nil {
		return nil, err
	}

	page.pointerTexture, err = ui.LoadTextureFromBMP(renderer, "assets/apu-gauge-pointer.bmp")
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

func (p *apuPage) render(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Copy(p.backgroundTexture, nil, nil)

	renderer.CopyEx(p.pointerTexture, nil, &sdl.Rect{237, 232, 2, 45}, (p.n1*175.0/100.0)+60.0, &sdl.Point{0, 0}, sdl.FLIP_NONE)
	renderer.CopyEx(p.pointerTexture, nil, &sdl.Rect{237, 322, 2, 45}, (p.egt*135.0/1000.0)+60.0, &sdl.Point{0, 0}, sdl.FLIP_NONE)

	p.n1Text.Render(240, 235)
	p.egtText.Render(240, 325)
	p.flapOpenText.Render(370, 270)

	renderer.Present()
}

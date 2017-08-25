package main

import (
	"log"

	"fmt"

	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/event/remote"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {

	log.Printf("[main] Initializing screen")
	ttf.Init()

	sdl.SetHintWithPriority(sdl.HINT_RENDER_SCALE_QUALITY, "best", sdl.HINT_OVERRIDE)

	win, _ := sdl.CreateWindow(
		"SimAvionics A320 APU",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480,
		sdl.WINDOW_SHOWN,
	)
	defer win.Destroy()

	render, _ := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	defer render.Destroy()

	background, _ := loadTexture(render, "assets/apu-background.bmp")
	defer background.Destroy()

	pointer, _ := loadTexture(render, "assets/apu-gauge-pointer.bmp")
	defer pointer.Destroy()

	font, err := ttf.OpenFont("assets/fonts/Carlito-Regular.ttf", 18)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	log.Printf("[main] Initializing SimAvionics remote bus")
	bus, _ := remote.NewSlaveEventBus("tcp://localhost:7001")

	var n1 = 0.0
	var egt = 0.0

	n1Text := newRenderedValue(render, font, sdl.Color{0, 200, 0, 255})
	egtText := newRenderedValue(render, font, sdl.Color{0, 200, 0, 255})
	flapOpenText := newRenderedValue(render, font, sdl.Color{0, 200, 0, 255})

	n1Chan := bus.Subscribe(apu.EngineStateN1)
	egtChan := bus.Subscribe(apu.EngineStateEGT)
	flapOpenChan := bus.Subscribe(apu.StatusFlapOpen)
	for {
		select {
		case v := <-n1Chan:
			n1 = v.Float64()
			n1Text.set(int(n1))
		case v := <-egtChan:
			egt = v.Float64()
			egtText.set(int(egt))
		case v := <-flapOpenChan:
			if v.Bool() {
				flapOpenText.set("FLAP OPEN")
			} else {
				flapOpenText.set("")
			}
		default:
		}

		render.SetDrawColor(0, 0, 0, 255)
		render.Clear()
		render.Copy(background, nil, nil)

		render.CopyEx(pointer, nil, &sdl.Rect{237, 232, 2, 45}, (n1*175.0/100.0)+60.0, &sdl.Point{0, 0}, sdl.FLIP_NONE)
		render.CopyEx(pointer, nil, &sdl.Rect{237, 322, 2, 45}, (egt*135.0/1000.0)+60.0, &sdl.Point{0, 0}, sdl.FLIP_NONE)

		n1Text.render(240, 235)
		egtText.render(240, 325)
		flapOpenText.render(370, 270)

		render.Present()

		sdl.PollEvent()
	}
}

type renderedValue struct {
	value    interface{}
	renderer *sdl.Renderer
	font     *ttf.Font
	color    sdl.Color
	texture  *sdl.Texture
	w        int32
	h        int32
}

func newRenderedValue(renderer *sdl.Renderer, font *ttf.Font, color sdl.Color) *renderedValue {
	value := &renderedValue{renderer: renderer, font: font, color: color}
	return value
}

func (rv *renderedValue) set(v interface{}) {
	if v != rv.value {
		rv.value = v
		if rv.texture != nil {
			rv.texture.Destroy()
		}
		texture, w, h, err := renderTextTexture(rv.renderer, rv.font, fmt.Sprintf("%v", rv.value), rv.color)
		if err != nil {
			panic(err)
		}
		rv.texture = texture
		rv.w = w
		rv.h = h
	}
}

func (rv *renderedValue) render(x int32, y int32) {
	if rv.texture != nil {
		rv.renderer.Copy(rv.texture, nil, &sdl.Rect{x, y, rv.w, rv.h})
	}
}

func loadTexture(render *sdl.Renderer, file string) (*sdl.Texture, error) {
	surface, err := sdl.LoadBMP(file)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := render.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	return texture, nil
}

func renderTextTexture(render *sdl.Renderer, font *ttf.Font, text string, color sdl.Color) (*sdl.Texture, int32, int32, error) {
	surface, err := font.RenderUTF8_Blended(text, color)
	if err != nil {
		return nil, 0, 0, err
	}
	defer surface.Free()

	texture, err := render.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, 0, 0, err
	}
	return texture, surface.W, surface.H, nil
}

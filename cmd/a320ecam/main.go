package main

import (
	"log"

	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/event/remote"
	"github.com/apoloval/simavionics/ui"
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

	background, _ := ui.LoadTextureFromBMP(render, "assets/apu-background.bmp")
	defer background.Destroy()

	pointer, _ := ui.LoadTextureFromBMP(render, "assets/apu-gauge-pointer.bmp")
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

	n1Text := ui.NewValueRenderer(render, font, sdl.Color{0, 200, 0, 255})
	egtText := ui.NewValueRenderer(render, font, sdl.Color{0, 200, 0, 255})
	flapOpenText := ui.NewValueRenderer(render, font, sdl.Color{0, 200, 0, 255})

	n1Chan := bus.Subscribe(apu.EngineStateN1)
	egtChan := bus.Subscribe(apu.EngineStateEGT)
	flapOpenChan := bus.Subscribe(apu.StatusFlapOpen)
	for {
		select {
		case v := <-n1Chan:
			n1 = v.Float64()
			n1Text.SetValue(int(n1))
		case v := <-egtChan:
			egt = v.Float64()
			egtText.SetValue(int(egt))
		case v := <-flapOpenChan:
			if v.Bool() {
				flapOpenText.SetValue("FLAP OPEN")
			} else {
				flapOpenText.SetValue("")
			}
		default:
		}

		render.SetDrawColor(0, 0, 0, 255)
		render.Clear()
		render.Copy(background, nil, nil)

		render.CopyEx(pointer, nil, &sdl.Rect{237, 232, 2, 45}, (n1*175.0/100.0)+60.0, &sdl.Point{0, 0}, sdl.FLIP_NONE)
		render.CopyEx(pointer, nil, &sdl.Rect{237, 322, 2, 45}, (egt*135.0/1000.0)+60.0, &sdl.Point{0, 0}, sdl.FLIP_NONE)

		n1Text.Render(240, 235)
		egtText.Render(240, 325)
		flapOpenText.Render(370, 270)

		render.Present()

		sdl.PollEvent()
	}
}

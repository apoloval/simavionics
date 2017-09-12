package main

import (
	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
	"github.com/apoloval/simavionics/ui"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	greenColor = 0x00C000FF
	amberColor = 0xCC7832FF
)

type apuPage struct {
	bleedValve bool
	bleedPSI   float64
	n1         float64
	egt        float64

	isEnergized bool

	eventChanAvailable  <-chan simavionics.EventValue
	eventChanBleed      <-chan simavionics.EventValue
	eventChanBleedValve <-chan simavionics.EventValue
	eventChanEnergized  <-chan simavionics.EventValue
	eventChanEngineN1   <-chan simavionics.EventValue
	eventChanEngineEGT  <-chan simavionics.EventValue
	eventChanFlap       <-chan simavionics.EventValue
	eventChanGenPerc    <-chan simavionics.EventValue
	eventChanGenVolt    <-chan simavionics.EventValue
	eventChanGenFreq    <-chan simavionics.EventValue

	backgroundTexture *sdl.Texture
	pointerTexture    *sdl.Texture

	textBleedPSI  *ui.ValueRenderer
	textEngineN1  *ui.ValueRenderer
	textEngineEGT *ui.ValueRenderer
	textEngineXX  *ui.ValueRenderer
	textGenPerc   *ui.ValueRenderer
	textGenVolt   *ui.ValueRenderer
	textGenFreq   *ui.ValueRenderer
	textFlapOpen  *ui.ValueRenderer
	textAvailable *ui.ValueRenderer
}

func newAPUPage(bus simavionics.EventBus, display *ui.Display) (*apuPage, error) {
	renderer := display.Renderer()
	positioner := display.Positioner()

	smallFontSize := positioner.Map(ui.RectF{H: 0.0375})
	smallFont, err := ttf.OpenFont("assets/fonts/Carlito-Regular.ttf", int(smallFontSize.H))
	if err != nil {
		return nil, err
	}

	bigFontSize := positioner.Map(ui.RectF{H: 0.045})
	bigFont, err := ttf.OpenFont("assets/fonts/Carlito-Regular.ttf", int(bigFontSize.H))
	if err != nil {
		return nil, err
	}

	page := &apuPage{
		eventChanAvailable:  bus.Subscribe(apu.EventAvailable),
		eventChanBleed:      bus.Subscribe(apu.EventBleed),
		eventChanBleedValve: bus.Subscribe(apu.EventBleedValve),
		eventChanEnergized:  bus.Subscribe(apu.EventMaster),
		eventChanEngineN1:   bus.Subscribe(apu.EventEngineN1),
		eventChanEngineEGT:  bus.Subscribe(apu.EventEngineEGT),
		eventChanFlap:       bus.Subscribe(apu.EventFlap),
		eventChanGenPerc:    bus.Subscribe(apu.EventGenPercentage),
		eventChanGenVolt:    bus.Subscribe(apu.EventGenVoltage),
		eventChanGenFreq:    bus.Subscribe(apu.EventGenFrequency),
		textBleedPSI:        ui.NewValueRenderer(renderer, smallFont, ui.NewColor(greenColor)),
		textEngineN1:        ui.NewValueRenderer(renderer, smallFont, ui.NewColor(greenColor)),
		textEngineEGT:       ui.NewValueRenderer(renderer, smallFont, ui.NewColor(greenColor)),
		textEngineXX:        ui.NewValueRenderer(renderer, smallFont, ui.NewColor(amberColor)),
		textGenPerc:         ui.NewValueRenderer(renderer, smallFont, ui.NewColor(greenColor)),
		textGenVolt:         ui.NewValueRenderer(renderer, smallFont, ui.NewColor(amberColor)),
		textGenFreq:         ui.NewValueRenderer(renderer, smallFont, ui.NewColor(amberColor)),
		textFlapOpen:        ui.NewValueRenderer(renderer, smallFont, ui.NewColor(greenColor)),
		textAvailable:       ui.NewValueRenderer(renderer, bigFont, ui.NewColor(greenColor)),
	}

	// Set the initial values for texts
	page.textBleedPSI.SetValue(0)
	page.textEngineN1.SetValue(0)
	page.textEngineEGT.SetIntValue(0, 5, 10)
	page.textEngineXX.SetValue("XX")
	page.textGenPerc.SetValue(0)
	page.textGenVolt.SetValue("XX")
	page.textGenFreq.SetValue("XX")

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
	for {
		select {
		case v := <-p.eventChanAvailable:
			if v.Bool() {
				p.textAvailable.SetValue("AVAIL")
			} else {
				p.textAvailable.SetValue("")
			}
		case v := <-p.eventChanBleedValve:
			p.bleedValve = v.Bool()
		case v := <-p.eventChanBleed:
			p.textBleedPSI.SetValue(int(v.Float64()))
		case v := <-p.eventChanEngineN1:
			p.n1 = v.Float64()
			p.textEngineN1.SetValue(int(p.n1))
		case v := <-p.eventChanEngineEGT:
			p.egt = v.Float64()
			p.textEngineEGT.SetIntValue(int(p.egt), 5, 10)
		case v := <-p.eventChanGenPerc:
			p.textGenPerc.SetValue(int(v.Float64()))
		case v := <-p.eventChanGenVolt:
			if v.Float64() == 0.0 {
				p.textGenVolt.SetColor(ui.NewColor(amberColor))
				p.textGenVolt.SetValue("XX")
			} else {
				p.textGenVolt.SetColor(ui.NewColor(greenColor))
				p.textGenVolt.SetValue(int(v.Float64()))
			}
		case v := <-p.eventChanGenFreq:
			if v.Float64() == 0.0 {
				p.textGenFreq.SetColor(ui.NewColor(amberColor))
				p.textGenFreq.SetValue("XX")
			} else {
				p.textGenFreq.SetColor(ui.NewColor(greenColor))
				p.textGenFreq.SetValue(int(v.Float64()))
			}
		case v := <-p.eventChanFlap:
			if v.Bool() {
				p.textFlapOpen.SetValue("FLAP OPEN")
			} else {
				p.textFlapOpen.SetValue("")
			}
		case v := <-p.eventChanEnergized:
			p.isEnergized = v.Bool()
		default:
			return
		}
	}
}

func (p *apuPage) render(display *ui.Display) {
	renderer := display.Renderer()

	p.renderBackground(display)
	p.renderEngineParams(display)
	p.renderBleed(display)
	p.renderMessages(display)
	p.renderGen(display)

	renderer.Present()
}

func (p *apuPage) renderBackground(display *ui.Display) {
	renderer := display.Renderer()

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Copy(p.backgroundTexture, nil, nil)
}

func (p *apuPage) renderEngineParams(display *ui.Display) {
	renderer := display.Renderer()
	positioner := display.Positioner()

	n1PointerRect := positioner.Map(ui.RectF{X: 0.365625, Y: 0.485416, W: 0.007812, H: 0.097917})
	egtPointerRect := positioner.Map(ui.RectF{X: 0.365625, Y: 0.672916, W: 0.007812, H: 0.097917})
	textRectEngineN1 := positioner.Map(ui.RectF{X: 0.375, Y: 0.489583})
	textRectEngineEGT := positioner.Map(ui.RectF{X: 0.375, Y: 0.677083})

	if p.isEnergized {
		renderer.CopyEx(p.pointerTexture, nil, &n1PointerRect, (p.n1*170.0/100.0)+65.0, &sdl.Point{2, 1}, sdl.FLIP_NONE)
		renderer.CopyEx(p.pointerTexture, nil, &egtPointerRect, (p.egt*125.0/1000.0)+65.0, &sdl.Point{2, 1}, sdl.FLIP_NONE)
		p.textEngineN1.Render(textRectEngineN1.X, textRectEngineN1.Y)
		p.textEngineEGT.Render(textRectEngineEGT.X, textRectEngineEGT.Y)
	} else {
		p.textEngineXX.Render(textRectEngineN1.X, textRectEngineN1.Y)
		p.textEngineXX.Render(textRectEngineEGT.X, textRectEngineEGT.Y)
	}
}

func (p *apuPage) renderMessages(display *ui.Display) {
	positioner := display.Positioner()

	textRectFlapOpen := positioner.Map(ui.RectF{X: 0.578125, Y: 0.5625})
	textRectAvailable := positioner.Map(ui.RectF{X: 0.45, Y: 0.15})

	p.textFlapOpen.Render(textRectFlapOpen.X, textRectFlapOpen.Y)
	p.textAvailable.Render(textRectAvailable.X, textRectAvailable.Y)
}

func (p *apuPage) renderBleed(display *ui.Display) {
	renderer := display.Renderer()
	positioner := display.Positioner()

	renderer.SetDrawColor(0, 255, 0, 255)
	if p.bleedValve {
		x1, y1 := positioner.MapCoords(0.6630859375, 0.167317708333333)
		x2, y2 := positioner.MapCoords(0.6630859375, 0.223958333333333)
		renderer.DrawLine(x1, y1, x2, y2)
	} else {
		x1, y1 := positioner.MapCoords(0.641845703125, 0.1953125)
		x2, y2 := positioner.MapCoords(0.6845703125, 0.1953125)
		renderer.DrawLine(x1, y1, x2, y2)
	}
	textRectBleedPSI := positioner.Map(ui.RectF{X: 0.63, Y: 0.3})
	p.textBleedPSI.Render(textRectBleedPSI.X, textRectBleedPSI.Y)
}

func (p *apuPage) renderGen(display *ui.Display) {
	positioner := display.Positioner()

	textRectPerc := positioner.Map(ui.RectF{X: 0.300, Y: 0.218})
	textRectVolt := positioner.Map(ui.RectF{X: 0.300, Y: 0.257})
	textRectFreq := positioner.Map(ui.RectF{X: 0.300, Y: 0.296})

	p.textGenPerc.Render(textRectPerc.X, textRectPerc.Y)
	p.textGenVolt.Render(textRectVolt.X, textRectVolt.Y)
	p.textGenFreq.Render(textRectFreq.X, textRectFreq.Y)
}

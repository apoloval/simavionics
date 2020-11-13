package ui

import "github.com/veandco/go-sdl2/sdl"

type RectF struct {
	X float64
	Y float64
	W float64
	H float64
}

type Positioner struct {
	xFactor float64
	yFactor float64
}

func NewPositioner(w int32, h int32, refW float64, refH float64) Positioner {
	return Positioner{
		xFactor: float64(w) / refW,
		yFactor: float64(h) / refH,
	}
}

func (p Positioner) Map(r RectF) sdl.Rect {
	return sdl.Rect{
		X: int32(r.X * p.xFactor),
		Y: int32(r.Y * p.yFactor),
		W: int32(r.W * p.xFactor),
		H: int32(r.H * p.yFactor),
	}
}

func (p Positioner) MapCoords(x float64, y float64) (int32, int32) {
	return int32(x * p.xFactor), int32(y * p.yFactor)
}

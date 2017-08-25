package ui

import "github.com/veandco/go-sdl2/sdl"

func LoadTextureFromBMP(render *sdl.Renderer, file string) (*sdl.Texture, error) {
	surface, err := sdl.LoadBMP(file)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := render.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	return texture, nil
}

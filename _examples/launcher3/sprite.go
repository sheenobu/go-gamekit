package main

import "github.com/veandco/go-sdl2/sdl"

type sprite struct {
	position sdl.Rect
	state    int
	sheet    *sheet
}

func newSprite(pos sdl.Rect, sheet *sheet, state int) *sprite {
	return &sprite{
		position: pos,
		sheet:    sheet,
		state:    state,
	}
}

func (sp *sprite) Render(r *sdl.Renderer) {
	sp.sheet.Copy(r, &sp.position, sp.state)
}

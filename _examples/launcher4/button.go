package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/rxgen/rx"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

func newButton(r *sdl.Rect, sheet *sheet, textureID int) *button {
	return &button{
		sprite:        newSprite(*r, sheet, textureID),
		Clicked:       rx.NewBool(false),
		clickedState:  false,
		hoveringState: false,
	}
}

type button struct {
	sprite *sprite

	Clicked *rx.Bool

	clickedState  bool
	hoveringState bool
}

func (b *button) Run(ctx context.Context, m *gamekit.Mouse) {

	posS := m.Position.Subscribe()
	clickS := m.LeftButtonState.Subscribe()
	defer posS.Close()
	defer clickS.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case pos := <-posS.C:
			x := pos.L
			y := pos.R
			p := b.sprite.position
			b.hoveringState = x > p.X && y > p.Y && x < p.X+p.W && y < p.Y+p.H
		case leftClick := <-clickS.C:
			if b.hoveringState && leftClick {
				b.clickedState = true
				b.Clicked.Set(true)
			} else {
				b.clickedState = false
				b.Clicked.Set(false)
			}
		}
	}
}

func (b *button) Render(r *sdl.Renderer) {
	b.sprite.Render(r)
}

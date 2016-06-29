package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type button struct {
	X int32
	Y int32
	W int32
	H int32

	T *sdl.Texture

	Click func()
}

func (b *button) Run(ctx context.Context, m *gamekit.Mouse) {
	posS := m.Position.Subscribe()
	clickS := m.LeftButtonState.Subscribe()
	defer posS.Close()
	defer clickS.Close()

	hovering := false

	for {
		select {
		case <-ctx.Done():
			return
		case pos := <-posS.C:
			x := pos.L
			y := pos.R
			hovering = x > b.X && y > b.Y && x < b.X+b.W && y < b.Y+b.H
		case leftClick := <-clickS.C:
			if hovering && leftClick {
				b.Click()
			}
		}
	}
}

func (b *button) Render(r *sdl.Renderer) {
	if b.T == nil {
		return
	}

	r.Copy(b.T, nil, &sdl.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H})
}

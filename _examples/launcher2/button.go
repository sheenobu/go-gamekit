package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/rxgen/rx"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

func newButton(r *sdl.Rect, t *sdl.Texture) *button {
	return &button{
		X:            r.X,
		Y:            r.Y,
		W:            r.W,
		H:            r.H,
		T:            t,
		Clicked:      rx.NewBool(false),
		clickedState: false,
	}
}

type button struct {
	X int32
	Y int32
	W int32
	H int32

	T *sdl.Texture

	Clicked      *rx.Bool
	clickedState bool
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
	if b.T == nil {
		return
	}

	r.Copy(b.T, nil, &sdl.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H})
}

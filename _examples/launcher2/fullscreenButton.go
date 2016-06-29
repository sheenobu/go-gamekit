package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type fullscreenButton struct {
	X int32
	Y int32
	W int32
	H int32

	checkbox *sdl.Texture
	self     <-chan bool
	others   []chan<- bool
	selected bool
}

func (fb *fullscreenButton) Run(ctx context.Context, m *gamekit.Mouse) {
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
			hovering = x > fb.X && y > fb.Y && x < fb.X+fb.W && y < fb.Y+fb.H
		case s, more := <-fb.self:
			if !more {
				return
			}
			fb.selected = s
		case leftClick := <-clickS.C:
			if hovering && leftClick {

				// mark self as selected
				fb.selected = true

				// mark all others in group as unselected
				for _, o := range fb.others {
					o <- false
				}
			}
		}
	}
}

func (fb *fullscreenButton) Render(r *sdl.Renderer) {

	if fb.selected {
		r.Copy(fb.checkbox, nil, &sdl.Rect{
			X: 107*2 + fb.X, Y: 9*2 + fb.Y, W: 9 * 2, H: 9 * 2,
		})
	}

}

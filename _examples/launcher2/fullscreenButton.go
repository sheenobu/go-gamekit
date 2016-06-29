package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type fullscreenButton struct {
	button *button

	checkbox *sdl.Texture
	self     <-chan bool
	others   []chan<- bool
	selected bool
}

func newFullscreenButton(r *sdl.Rect, t *sdl.Texture) *fullscreenButton {
	return &fullscreenButton{
		button:   newButton(r, nil),
		checkbox: t,
	}
}

func (fb *fullscreenButton) Run(ctx context.Context, m *gamekit.Mouse) {
	clickSub := fb.button.Clicked.Subscribe()
	defer clickSub.Close()

	go fb.button.Run(ctx, m)

	for {
		select {
		case <-ctx.Done():
			return
		case s, more := <-fb.self:
			if !more {
				return
			}
			fb.selected = s
		case clicked := <-clickSub.C:
			if clicked {
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

	fb.button.Render(r)

	if fb.selected {
		p := fb.button
		r.Copy(fb.checkbox, nil, &sdl.Rect{
			X: 107*2 + p.X, Y: 9*2 + p.Y, W: 9 * 2, H: 9 * 2,
		})
	}
}

package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type fullscreenButton struct {
	button   *button
	checkbox *sprite

	self     <-chan bool
	others   []chan<- bool
	selected bool
}

func newFullscreenButton(r *sdl.Rect, sheet *sheet, textureID int, checkboxID int) *fullscreenButton {
	return &fullscreenButton{
		button:   newButton(r, sheet, textureID),
		checkbox: newSprite(sdl.Rect{X: 107*2 + r.X, Y: 9*2 + r.Y, W: 9 * 2, H: 9 * 2}, sheet, checkboxID),
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
		fb.checkbox.Render(r)
	}
}

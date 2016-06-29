package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type toggleButton struct {
	button   *button
	checkbox *sprite

	self     <-chan bool
	others   []chan<- bool
	selected bool
}

func newToggleButton(r *sdl.Rect, checkboxOffset *sdl.Rect, sheet *sheet, textureID int, checkboxID int) *toggleButton {
	return &toggleButton{
		button:   newButton(r, sheet, textureID),
		checkbox: newSprite(sdl.Rect{X: checkboxOffset.X + r.X, Y: checkboxOffset.Y + r.Y, W: checkboxOffset.W, H: checkboxOffset.H}, sheet, checkboxID),
	}
}

func (tb *toggleButton) Run(ctx context.Context, m *gamekit.Mouse) {
	clickSub := tb.button.Clicked.Subscribe()
	defer clickSub.Close()

	go tb.button.Run(ctx, m)

	for {
		select {
		case <-ctx.Done():
			return
		case s, more := <-tb.self:
			if !more {
				return
			}
			tb.selected = s
		case clicked := <-clickSub.C:
			if clicked {
				// mark self as selected
				tb.selected = true

				// mark all others in group as unselected
				for _, o := range tb.others {
					o <- false
				}
			}
		}
	}
}

func (tb *toggleButton) Render(r *sdl.Renderer) {

	tb.button.Render(r)

	if tb.selected {
		tb.checkbox.Render(r)
	}
}

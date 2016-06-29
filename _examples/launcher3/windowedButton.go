package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type windowedButton struct {
	button   *button
	checkbox *sprite

	self     <-chan bool
	others   []chan<- bool
	selected bool
}

func newWindowedButton(r *sdl.Rect, sheet *sheet, textureID int, checkboxID int) *windowedButton {
	return &windowedButton{
		button:   newButton(r, sheet, textureID),
		checkbox: newSprite(sdl.Rect{X: 107*2 + r.X, Y: 9*2 + r.Y, W: 9 * 2, H: 9 * 2}, sheet, checkboxID),
	}
}

func (wb *windowedButton) Run(ctx context.Context, m *gamekit.Mouse) {
	clickSub := wb.button.Clicked.Subscribe()
	defer clickSub.Close()

	go wb.button.Run(ctx, m)

	for {
		select {
		case <-ctx.Done():
			return
		case s, more := <-wb.self:
			if !more {
				return
			}
			wb.selected = s
		case clicked := <-clickSub.C:
			if clicked {
				// mark self as selected
				wb.selected = true

				// mark all others in group as unselected
				for _, o := range wb.others {
					o <- false
				}
			}
		}
	}
}

func (wb *windowedButton) Render(r *sdl.Renderer) {

	wb.button.Render(r)

	if wb.selected {
		wb.checkbox.Render(r)
	}
}

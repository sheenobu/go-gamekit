package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type windowedButton struct {
	X int32
	Y int32
	W int32
	H int32

	checkbox *sdl.Texture
	self     <-chan bool
	others   []chan<- bool
	selected bool
}

func (wb *windowedButton) Run(ctx context.Context, m *gamekit.Mouse) {
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
			hovering = x > wb.X && y > wb.Y && x < wb.X+wb.W && y < wb.Y+wb.H
		case s, more := <-wb.self:
			if !more {
				return
			}
			wb.selected = s
		case leftClick := <-clickS.C:
			if hovering && leftClick {

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

	if wb.selected {
		r.Copy(wb.checkbox, nil, &sdl.Rect{
			X: 107*2 + wb.X, Y: 9*2 + wb.Y, W: 9 * 2, H: 9 * 2,
		})
	}

}

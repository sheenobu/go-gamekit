package main

import (
	"github.com/sheenobu/go-gamekit"
	"golang.org/x/net/context"
)

type closeButton struct {
	X int32
	Y int32
	W int32
	H int32

	Close func()
}

func (cb *closeButton) Run(ctx context.Context, m *gamekit.Mouse) {
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
			hovering = x > cb.X && y > cb.Y && x < cb.X+cb.W && y < cb.Y+cb.H
		case leftClick := <-clickS.C:
			if hovering && leftClick {
				cb.Close()
			}
		}
	}
}

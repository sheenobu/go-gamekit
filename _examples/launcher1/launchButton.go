package main

import (
	"github.com/sheenobu/go-gamekit"
	"golang.org/x/net/context"
)

type launchButton struct {
	X int32
	Y int32
	W int32
	H int32

	Launch func()
}

func (lb *launchButton) Run(ctx context.Context, m *gamekit.Mouse) {
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
			hovering = x > lb.X && y > lb.Y && x < lb.X+lb.W && y < lb.Y+lb.H
		case leftClick := <-clickS.C:
			if hovering && leftClick {
				lb.Launch()
			}
		}
	}
}

package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/rxgen/rx"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// Button represents a button that can be pressed
type Button struct {
	ClickState *rx.Bool

	X int32
	Y int32

	clicked bool
	hovered bool

	renderer *sdl.Renderer
}

// NewButton creates a new button at the given position with the SDL renderer
func NewButton(x int32, y int32, r *sdl.Renderer) *Button {
	return &Button{
		X:        x,
		Y:        y,
		clicked:  false,
		hovered:  false,
		renderer: r,

		ClickState: rx.NewBool(false),
	}
}

// Position returns the position of the button
func (b *Button) Position() (int32, int32) {
	return b.X, b.Y
}

// Move moves the button
func (b *Button) Move(x int32, y int32) {
	b.X = x
	b.Y = y
}

// Render draws the button
func (b *Button) Render() {

	rect := sdl.Rect{
		W: 50,
		H: 20,
		X: b.X,
		Y: b.Y,
	}
	if b.clicked {
		b.renderer.SetDrawColor(0, 255, 255, 0)
	} else {
		b.renderer.SetDrawColor(255, 255, 255, 0)
	}
	b.renderer.FillRect(&rect)

	if b.hovered {
		b.renderer.SetDrawColor(0, 255, 255, 0)
	} else {
		b.renderer.SetDrawColor(0, 255, 0, 0)
	}

	b.renderer.DrawRect(&rect)
}

// Run runs the buttons event handling processes
func (b *Button) Run(ctx context.Context, win *gamekit.Window) {

	mousePosition := win.Mouse.Position.Subscribe()
	mouseLeftClick := win.Mouse.LeftButtonState.Subscribe()

	defer mousePosition.Close()
	defer mouseLeftClick.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case pos := <-mousePosition.C:
			x := pos.L
			y := pos.R

			if x > b.X && y > b.Y && x < b.X+50 && y < b.Y+20 {
				b.hovered = true
			} else {
				b.hovered = false
			}

		case leftClick := <-mouseLeftClick.C:
			if b.hovered && leftClick {
				if b.clicked {
					continue
				}
				b.clicked = true
			} else {
				if !b.clicked {
					continue
				}
				b.clicked = false
			}

			b.ClickState.Set(b.clicked)
		}
	}
}

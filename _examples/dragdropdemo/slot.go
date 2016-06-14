package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// Slot represents something that can hold a visual item
type Slot struct {
	Item *Item

	X int32
	Y int32
	W int32
	H int32

	Hovered bool

	renderer *sdl.Renderer
}

// NewSlot creates a new slot at the given position with the SDL renderer
func NewSlot(i *Item, x int32, y int32, w int32, h int32, r *sdl.Renderer) *Slot {

	s := &Slot{
		X:        x,
		Y:        y,
		W:        w,
		H:        h,
		renderer: r,
	}

	i.Reparent(s)

	return s
}

// Render draws the slot graphics to the SDL renderer
func (b *Slot) Render() {

	rect := sdl.Rect{
		W: b.W,
		H: b.H,
		X: b.X,
		Y: b.Y,
	}

	if b.Hovered {
		b.renderer.SetDrawColor(0, 255, 255, 0)
	} else {
		b.renderer.SetDrawColor(0, 255, 0, 0)
	}

	b.renderer.DrawRect(&rect)
}

// Run runs the slots event handlers
func (b *Slot) Run(ctx context.Context, win *gamekit.Window) {
	mousePosition := win.Mouse.Position.Subscribe()
	defer mousePosition.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case pos := <-mousePosition.C:
			x := pos.L
			y := pos.R

			if x > b.X && y > b.Y && x < b.X+b.W && y < b.Y+b.H {
				b.Hovered = true
			} else {
				b.Hovered = false
			}

		}
	}
}

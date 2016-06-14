package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/rxgen/rx"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// Item represents something that can be dragged
type Item struct {
	ClickState *rx.Bool
	Slot       *Slot

	X int32
	Y int32
	W int32
	H int32

	clicked bool
	hovered bool

	renderer *sdl.Renderer
}

// NewItem creates a new item at the given position with the SDL renderer
func NewItem(x int32, y int32, w int32, h int32, r *sdl.Renderer) *Item {
	return &Item{
		Slot:     nil,
		X:        x,
		Y:        y,
		W:        w,
		H:        h,
		clicked:  false,
		hovered:  false,
		renderer: r,

		ClickState: rx.NewBool(false),
	}
}

// Render draws the item to the SDL renderer
func (b *Item) Render() {

	rect := sdl.Rect{
		W: b.W,
		H: b.H,
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

// Run runs the item processes for tracking events
func (b *Item) Run(ctx context.Context, win *gamekit.Window) {

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

			if x > b.X && y > b.Y && x < b.X+b.W && y < b.Y+b.H {
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

// Reparent sets the parent of the new item, moving and resizing it to fit the slot
func (b *Item) Reparent(s *Slot) {

	if b == nil {
		return
	}

	// resize
	b.X = s.X + 2
	b.Y = s.Y + 2
	b.W = s.W - 4
	b.H = s.H - 4

	// if we are reparenting with ourselves
	if s.Item == b {
		return
	}

	// unparent the original parent
	if b.Slot != nil {
		b.Slot.Item = nil
	}

	// reparent
	b.Slot = s
	s.Item = b
}

//-- Draggable

// Position returns the position of the item
func (b *Item) Position() (int32, int32) {
	return b.X, b.Y
}

// Move moves the item
func (b *Item) Move(x int32, y int32) {
	b.X = x
	b.Y = y
}

//--

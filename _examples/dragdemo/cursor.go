package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// mouse cursor represented as a reactive go routine and a render function

// Cursor is the representation of the mouse as a graphical item that follows the mouses X and Y position
type Cursor struct {
	x int32
	y int32

	leftClicked bool

	renderer *sdl.Renderer
}

// NewCursor creates a new cursor that renders to the SDL renderer
func NewCursor(r *sdl.Renderer) *Cursor {
	return &Cursor{
		renderer: r,
	}
}

// Render renders the cursor
func (c *Cursor) Render() {

	rect := sdl.Rect{
		W: 20,
		H: 20,
		X: c.x - 10,
		Y: c.y - 10,
	}

	if c.leftClicked {
		c.renderer.SetDrawColor(0, 255, 0, 0)
	} else {
		c.renderer.SetDrawColor(255, 0, 0, 0)
	}

	c.renderer.FillRect(&rect)
}

// Run runs the cursors event handling routines
func (c *Cursor) Run(ctx context.Context, win *gamekit.Window) {

	mousePosition := win.Mouse.Position.Subscribe()
	mouseLeftClick := win.Mouse.LeftButtonState.Subscribe()

	defer mousePosition.Close()
	defer mouseLeftClick.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case pos := <-mousePosition.C:
			c.x = pos.L
			c.y = pos.R
		case leftClick := <-mouseLeftClick.C:
			c.leftClicked = leftClick
		}
	}
}

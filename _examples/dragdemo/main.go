package main

import (
	"runtime"

	"golang.org/x/net/context"

	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/loop"
)

func main() {

	// initialize the engine
	runtime.LockOSThread()
	gamekit.Init()

	// create the window
	wm := gamekit.NewWindowManager()
	win, err := wm.NewWindow("dragdemo", 800, 600)
	if err != nil {
		panic(err)
	}

	// create our context to handle lifecycle events
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create and run our graphical mouse cursor
	c := NewCursor(win.Renderer)
	go c.Run(ctx, win)

	// create and run our button
	b := NewButton(200, 100, win.Renderer)
	go b.Run(ctx, win)

	// make the button draggable by attaching the Mouse's X and Y to the
	// buttons X and Y, only if the button is pressed
	buttonPress := b.ClickState.Subscribe()
	defer buttonPress.Close()
	EnableDragging(ctx, b, win.Mouse, buttonPress)

	loop.Simple(wm, func() {

		win.Renderer.SetDrawColor(0, 0, 0, 0)
		win.Renderer.Clear()

		b.Render()
		c.Render()

		win.Renderer.Present()
	}).Run()

}

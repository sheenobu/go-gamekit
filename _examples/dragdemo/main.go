package main

import (
	"runtime"
	"time"

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
	win, err := wm.NewWindow("dragdemo", 800, 600, 0)
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

	// create a histogram for goroutine count
	hs := gamekit.CountHistogram{}

	go hs.Run(ctx, 300*time.Millisecond, func() int {
		return runtime.NumGoroutine()
	})

	loop.Simple(wm, func() {
		win.Renderer.SetDrawColor(128, 128, 128, 0)
		win.Renderer.Clear()

		b.Render()
		c.Render()

		hs.Render(win.Renderer)

		win.Renderer.Present()
	}).Run()

}

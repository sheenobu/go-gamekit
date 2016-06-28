package main

import (
	"runtime"

	"golang.org/x/net/context"

	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/loop"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize the engine
	runtime.LockOSThread()
	gamekit.Init()

	// build the main window
	wm := gamekit.NewWindowManager()
	win, err := wm.NewWindow("hello-world", 800, 600, 0)
	if err != nil {
		panic(err)
	}

	// build and run the simple game loop
	loop.Simple(wm, ctx, func() {
		win.Renderer.SetDrawColor(128, 128, 128, 255)
		win.Renderer.Clear()

		win.Renderer.Present()
	}).Run()

}

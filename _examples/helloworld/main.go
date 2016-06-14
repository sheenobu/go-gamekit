package main

import (
	"runtime"

	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/loop"
)

func main() {

	// initialize the engine
	runtime.LockOSThread()
	gamekit.Init()

	// build the main window
	wm := gamekit.NewWindowManager()
	win, err := wm.NewWindow("hello-world", 800, 600)
	if err != nil {
		panic(err)
	}

	// build and run the simple game loop
	loop.Simple(wm, func() {
		win.Renderer.Clear()

		win.Renderer.Present()
	}).Run()

}

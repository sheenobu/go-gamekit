package main

import (
	"fmt"
	"runtime"

	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {
	// initialize the engine
	runtime.LockOSThread()
	gamekit.Init()

	// Initialize TTF (TODO: move to gamekit.Init)
	ttf.Init()

	// run the launcher
	results := runLauncher()

	// if the launch flag is true, launch the game
	if results.Launch {
		game(results)
	}
}

func game(res launchResults) {
	dm := res.ChosenResolution
	fmt.Printf("Launching game\n")
	fmt.Printf("Fullscreen: %v\n", res.Fullscreen)
	fmt.Printf("Selected Resolution: %d/%d (%d)\n", dm.W, dm.H, dm.RefreshRate)
}

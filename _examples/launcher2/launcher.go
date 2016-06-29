package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/loop"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"golang.org/x/net/context"
)

type launchResults struct {
	Launch           bool
	ChosenResolution sdl.DisplayMode
	Fullscreen       bool
}

func runLauncher() (res launchResults) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// load the launch assets
	sfc, err := img.Load("./data/zd-launch.png")
	if err != nil {
		panic(err)
	}
	winW := int(sfc.W * 2)
	winH := int(sfc.H * 2)
	sfc.Free()

	// build the launcher window
	wm := gamekit.NewWindowManager()
	win, err := wm.NewWindow("zd - launch", winW, winH, 0)
	if err != nil {
		panic(err)
	}

	// load the launcher texture
	launchTexture, err := img.LoadTexture(win.Renderer, "./data/zd-launch.png")
	if err != nil {
		panic(err)
	}

	res.Launch = false
	launchOn := func() {
		res.Launch = true
		cancel()
	}

	// load the checkobx texture
	checkbox, err := img.LoadTexture(win.Renderer, "./data/checkbox.png")
	if err != nil {
		panic(err)
	}

	// run the interactive elements
	cb := &closeButton{133 * 2, 49 * 2, 53 * 2, 13 * 2, cancel}
	go cb.Run(ctx, win.Mouse)

	lb := &launchButton{192 * 2, 49 * 2, 53 * 2, 13 * 2, launchOn}
	go lb.Run(ctx, win.Mouse)

	tg := &toggleGroup{}
	go tg.Run(ctx, checkbox, win.Mouse)

	rp := newResolutionPicker(sdl.Rect{X: 134 * 2, Y: 4 * 2, W: 100 * 2, H: 41 * 2}, win.Renderer)
	go rp.Run(ctx, win.Mouse, &res)

	// build and run the simple game loop
	loop.Simple(wm, ctx, func() {
		win.Renderer.SetDrawColor(128, 128, 128, 255)
		win.Renderer.Clear()

		win.Renderer.Copy(launchTexture, nil, &sdl.Rect{X: 0, Y: 0, W: int32(winW), H: int32(winH)})

		tg.Render(win.Renderer)
		rp.Render(win.Renderer)

		win.Renderer.Present()
	}).Run()

	res.Fullscreen = tg.IsFullscreen()

	return
}

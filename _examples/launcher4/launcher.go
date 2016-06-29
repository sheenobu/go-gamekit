package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/gfx2"
	"github.com/sheenobu/go-gamekit/loop"
	"github.com/sheenobu/go-gamekit/ui"
	"github.com/veandco/go-sdl2/sdl"
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

	winW := 247 * 2
	winH := 65 * 2

	// build the launcher window
	wm := gamekit.NewWindowManager()
	win, err := wm.NewWindow("zd - launch", winW, winH, 0)
	if err != nil {
		panic(err)
	}

	// load the sprite sheet
	launcherSheet := gfx2.NewSheet(win.Renderer, "./data/launcher_sheet.png")

	windowedButtonID := launcherSheet.Add(&sdl.Rect{X: 0, Y: 3, W: 123, H: 27})
	fullscreenButtonID := launcherSheet.Add(&sdl.Rect{X: 0, Y: 35, W: 123, H: 27})
	scrollRegionID := launcherSheet.Add(&sdl.Rect{X: 130, Y: 3, W: 112, H: 43})
	closeButtonID := launcherSheet.Add(&sdl.Rect{X: 130, Y: 49, W: 53, H: 13})
	launchButtonID := launcherSheet.Add(&sdl.Rect{X: 189, Y: 49, W: 53, H: 13})
	checkboxID := launcherSheet.Add(&sdl.Rect{X: 247, Y: 0, W: 9, H: 9})
	arrowUpID := launcherSheet.Add(&sdl.Rect{X: 247, Y: 10, W: 9, H: 7})
	arrowDownID := launcherSheet.Add(&sdl.Rect{X: 247, Y: 18, W: 9, H: 7})

	// run the interactive elements
	cb := ui.NewButton(&sdl.Rect{X: 133 * 2, Y: 49 * 2, W: 53 * 2, H: 13 * 2}, launcherSheet, closeButtonID)
	go cb.Run(ctx, win.Mouse)

	lb := ui.NewButton(&sdl.Rect{X: 192 * 2, Y: 49 * 2, W: 53 * 2, H: 13 * 2}, launcherSheet, launchButtonID)
	go lb.Run(ctx, win.Mouse)

	tg := ui.NewToggleGroup("windowed")
	tg.Add(ui.NewToggleButton("windowed", &sdl.Rect{X: 3 * 2, Y: 3 * 2, W: 123 * 2, H: 27 * 2}, &sdl.Rect{X: 107 * 2, Y: 9 * 2, W: 9 * 2, H: 9 * 2}, launcherSheet, windowedButtonID, checkboxID))
	tg.Add(ui.NewToggleButton("fullscreen", &sdl.Rect{X: 3 * 2, Y: 35 * 2, W: 123 * 2, H: 27 * 2}, &sdl.Rect{X: 107 * 2, Y: 9 * 2, W: 9 * 2, H: 9 * 2}, launcherSheet, fullscreenButtonID, checkboxID))
	go tg.Run(ctx, win.Mouse)

	rp := newResolutionPicker(sdl.Rect{X: 133 * 2, Y: 3 * 2, W: 100 * 2, H: 41 * 2}, win.Renderer, launcherSheet, scrollRegionID, arrowUpID, arrowDownID)
	go rp.Run(ctx, win.Mouse, &res)

	go func() {

		launchSub := lb.Clicked.Subscribe()
		defer launchSub.Close()

		closeSub := cb.Clicked.Subscribe()
		defer closeSub.Close()

		selectedSub := tg.Selected.Subscribe()
		defer selectedSub.Close()

		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			case name := <-selectedSub.C:
				res.Fullscreen = name == "fullscreen"
			case clicked := <-closeSub.C:
				if clicked {
					return
				}
			case clicked := <-launchSub.C:
				if clicked {
					res.Launch = true
					return
				}
			}
		}
	}()

	// build and run the simple game loop
	loop.Simple(wm, ctx, func() {
		win.Renderer.SetDrawColor(0, 0, 0, 255)
		win.Renderer.Clear()

		cb.Render(win.Renderer)
		lb.Render(win.Renderer)
		tg.Render(win.Renderer)
		rp.Render(win.Renderer)

		win.Renderer.Present()
	}).Run()

	return
}

package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/rxgen/rx"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// multiple buttons that are related by state
type toggleGroup struct {
	wb *toggleButton
	fb *toggleButton

	Selected *rx.String
}

func newToggleGroup(sheet *sheet, windowedButtonID int, fullscreenButtonID int, checkboxID int) *toggleGroup {

	var tg toggleGroup
	tg.wb = newToggleButton("windowed", &sdl.Rect{X: 3 * 2, Y: 3 * 2, W: 123 * 2, H: 27 * 2}, &sdl.Rect{X: 107 * 2, Y: 9 * 2, W: 9 * 2, H: 9 * 2}, sheet, windowedButtonID, checkboxID)
	tg.fb = newToggleButton("fullscreen", &sdl.Rect{X: 3 * 2, Y: 35 * 2, W: 123 * 2, H: 27 * 2}, &sdl.Rect{X: 107 * 2, Y: 9 * 2, W: 9 * 2, H: 9 * 2}, sheet, fullscreenButtonID, checkboxID)

	tg.Selected = rx.NewString("windowed")

	return &tg
}

func (tg *toggleGroup) Run(ctx context.Context, m *gamekit.Mouse) {

	tg.wb.isSelected = true
	tg.fb.isSelected = false

	go tg.wb.Run(ctx, m, tg.Selected)
	go tg.fb.Run(ctx, m, tg.Selected)

	<-ctx.Done()
}

func (tg *toggleGroup) Render(r *sdl.Renderer) {
	tg.wb.Render(r)
	tg.fb.Render(r)
}

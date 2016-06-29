package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// multiple buttons that are related by state
type toggleGroup struct {
	wb *toggleButton
	fb *toggleButton
}

func newToggleGroup(sheet *sheet, windowedButtonID int, fullscreenButtonID int, checkboxID int) *toggleGroup {

	var tg toggleGroup
	tg.wb = newToggleButton(&sdl.Rect{X: 3 * 2, Y: 3 * 2, W: 123 * 2, H: 27 * 2}, &sdl.Rect{X: 107 * 2, Y: 9 * 2, W: 9 * 2, H: 9 * 2}, sheet, windowedButtonID, checkboxID)
	tg.fb = newToggleButton(&sdl.Rect{X: 3 * 2, Y: 35 * 2, W: 123 * 2, H: 27 * 2}, &sdl.Rect{X: 107 * 2, Y: 9 * 2, W: 9 * 2, H: 9 * 2}, sheet, fullscreenButtonID, checkboxID)

	return &tg
}

func (tg *toggleGroup) Run(ctx context.Context, m *gamekit.Mouse) {

	wbSelf := make(chan bool)
	fbSelf := make(chan bool)

	defer close(wbSelf)
	defer close(fbSelf)

	tg.wb.self = wbSelf
	tg.wb.others = []chan<- bool{fbSelf}
	tg.wb.selected = true

	tg.fb.self = fbSelf
	tg.fb.others = []chan<- bool{wbSelf}
	tg.fb.selected = false

	go tg.wb.Run(ctx, m)
	go tg.fb.Run(ctx, m)

	<-ctx.Done()
}

func (tg *toggleGroup) Render(r *sdl.Renderer) {
	tg.wb.Render(r)
	tg.fb.Render(r)
}

func (tg *toggleGroup) IsFullscreen() bool {
	return tg.fb.selected
}

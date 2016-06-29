package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

// multiple buttons that are related by state
type toggleGroup struct {
	wb *windowedButton
	fb *fullscreenButton
}

func (tg *toggleGroup) Run(ctx context.Context, checkbox *sdl.Texture, m *gamekit.Mouse) {

	wbSelf := make(chan bool)
	fbSelf := make(chan bool)

	defer close(wbSelf)
	defer close(fbSelf)

	tg.wb = newWindowedButton(&sdl.Rect{X: 3 * 2, Y: 3 * 2, W: 123 * 2, H: 27 * 2}, checkbox)
	tg.wb.self = wbSelf
	tg.wb.others = []chan<- bool{fbSelf}
	tg.wb.selected = true

	tg.fb = newFullscreenButton(&sdl.Rect{X: 3 * 2, Y: 35 * 2, W: 123 * 2, H: 27 * 2}, checkbox)
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

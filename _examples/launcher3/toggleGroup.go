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

	tg.wb = &windowedButton{3 * 2, 3 * 2, 123 * 2, 27 * 2, checkbox, wbSelf, []chan<- bool{fbSelf}, true}
	tg.fb = &fullscreenButton{3 * 2, 35 * 2, 123 * 2, 27 * 2, checkbox, fbSelf, []chan<- bool{wbSelf}, false}

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

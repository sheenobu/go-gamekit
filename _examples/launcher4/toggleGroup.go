package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/rxgen/rx"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/net/context"
)

type toggleGroup struct {
	buttons []*toggleButton
	initial string

	Selected *rx.String
}

func newToggleGroup(defaultSelected string) *toggleGroup {

	var tg toggleGroup
	tg.Selected = rx.NewString(defaultSelected)
	tg.initial = defaultSelected

	return &tg
}

func (tg *toggleGroup) Add(tb *toggleButton) {
	tg.buttons = append(tg.buttons, tb)
	if tg.initial == tb.name {
		tb.isSelected = true
	}
}

func (tg *toggleGroup) Run(ctx context.Context, m *gamekit.Mouse) {

	for _, b := range tg.buttons {
		go b.Run(ctx, m, tg.Selected)
	}

	<-ctx.Done()
}

func (tg *toggleGroup) Render(r *sdl.Renderer) {
	for _, b := range tg.buttons {
		b.Render(r)
	}
}

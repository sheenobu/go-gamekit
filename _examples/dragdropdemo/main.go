package main

import (
	"runtime"
	"time"

	"golang.org/x/net/context"

	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/drag"
	"github.com/sheenobu/go-gamekit/loop"
)

func main() {

	// start the engine
	runtime.LockOSThread()
	gamekit.Init()

	// create our window
	wm := gamekit.NewWindowManager()
	win, err := wm.NewWindow("dragdropdemo", 800, 600, 0)
	if err != nil {
		panic(err)
	}

	var slots []*Slot
	var dslots []*Slot
	var items []*Item

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dx := drag.NewRxDraggable(nil)

	// create our drag and drop items and slots by running

	for i := int32(1); i != 6; i++ {
		w := int32(50)
		h := int32(50)

		var item *Item

		// only add an item for every other slot
		if i%2 != 0 {
			item = NewItem(0, 0, 0, 0, win.Renderer)
			go item.Run(ctx, win)
			go EnableDragging(ctx, item, win.Mouse, dx, item.ClickState.Subscribe())
			items = append(items, item)
		}

		slot := NewSlot(item, 100, i*h+8*i, w, h, win.Renderer)
		dslot := NewSlot(nil, 300, i*h+8*i, w, h, win.Renderer)

		go dslot.Run(ctx, win)
		go slot.Run(ctx, win)
		go EnableDropping(ctx, dslot, win.Mouse, dx, true)
		go EnableDropping(ctx, slot, win.Mouse, dx, false)

		dslots = append(dslots, dslot)
		slots = append(slots, slot)
	}

	// create a histogram for goroutine count
	hs := gamekit.CountHistogram{}

	go hs.Run(ctx, 300*time.Millisecond, func() int {
		return runtime.NumGoroutine()
	})

	// Run the gameloop

	loop.Simple(wm, ctx, func() {

		win.Renderer.SetDrawColor(128, 128, 128, 0)
		win.Renderer.Clear()

		for _, i := range slots {
			i.Render()
		}

		for _, i := range dslots {
			i.Render()
		}

		for _, i := range items {
			i.Render()
		}

		hs.Render(win.Renderer)

		win.Renderer.Present()

	}).Run()
}

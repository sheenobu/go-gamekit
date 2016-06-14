package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/drag"
	"github.com/sheenobu/go-gamekit/pair"
	"github.com/sheenobu/rxgen/rx"
	"golang.org/x/net/context"
)

// Moveable is an interface for an entity which can be moved around the screen
type Moveable interface {

	// Move moves the entity to the given x and y position
	Move(x, y int32)
}

// EnableDragging runs dragging subprocesses for the given object when the boolean subscriber is true
func EnableDragging(ctx context.Context, o drag.Draggable, mouse *gamekit.Mouse, rxdrag *drag.RxDraggable, toDrag *rx.BoolSubscriber) {
	go func() {

		var mouseLocation <-chan pair.Int32Pair
		var mouseLocationClose = func() {}

		var offsetX int32
		var offsetY int32

		for {
			select {
			case <-ctx.Done():
				mouseLocationClose()
				return
			case coords := <-mouseLocation:
				// move the object when the mouse location comes in
				o.Move(coords.L-offsetX, coords.R-offsetY)
			case dragStatus := <-toDrag.C:

				if dragStatus {

					// if clicked, subscribe on the mouse cursor position and update the object
					// when the mouse cursor moves

					// subscribe to the mouse location
					s := mouse.Position.Subscribe()
					mouseLocation = s.C
					mouseLocationClose = s.Close

					// get the current position to calculate offset
					// offset ensures the relative location of the item compared to the mouse stays the mouse
					ml := mouse.Position.Get()
					curX, curY := o.Position()
					offsetX = ml.L - curX
					offsetY = ml.R - curY

					// signal that dragging has started
					rxdrag.Set(o)
				} else {

					// disable mouse location listener
					mouseLocationClose()
					mouseLocation = nil
					mouseLocationClose = func() {}

					// signal that dragging has stopped
					rxdrag.Set(nil)
				}
			}
		}
	}()
}

// EnableDropping runs dropping subproceses for the given slot.
func EnableDropping(ctx context.Context, s *Slot, mouse *gamekit.Mouse, rxdrag *drag.RxDraggable, isDestination bool) {
	go func() {

		sub := rxdrag.Subscribe()
		defer sub.Close()

		var curDragging drag.Draggable

		var mouseLocation <-chan pair.Int32Pair
		var mouseLocationClose = func() {}

		var dragX, dragY int32

		for {
			select {
			case <-ctx.Done():
				mouseLocationClose()
				return
			case c := <-mouseLocation:
				dragX = c.L
				dragY = c.R

			case o := <-sub.C:

				if o != nil {
					// dragging has started, track the mouse for collision
					s := mouse.Position.Subscribe()
					mouseLocation = s.C
					mouseLocationClose = s.Close

					// and save the current dragging object
					curDragging = o

					continue
				}

				// dragging has stopped...

				// stop tracking the mouse position
				mouseLocationClose()
				mouseLocation = nil
				mouseLocationClose = func() {}

				// ensure our current draggable is compatible with our slot object
				i, ok := curDragging.(*Item)
				if !ok {
					continue
				}

				// collision check
				collision := dragX > s.X && dragY > s.Y && dragX < s.X+s.W && dragY < s.Y+s.H

				// if we collided with an empty slot or our own slot, move the object to this new slot.
				if isDestination && collision && s.Item == nil || s.Item == i {
					i.Reparent(s)
					continue
				}

				// if we did not collide with any slot and the managed slot is our current slot, then reposition. This
				// provides the 'snap-back' you see when dragging fails.
				if s != nil && i != nil && i.Slot == s {
					i.Reparent(s)
				}
			}
		}
	}()
}

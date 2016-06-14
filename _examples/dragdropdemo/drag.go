package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/drag"
	"github.com/sheenobu/go-gamekit/pair"
	"github.com/sheenobu/rxgen/rx"
)

// Moveable is an interface for an entity which can be moved around the screen
type Moveable interface {

	// Move moves the entity to the given x and y position
	Move(x, y int32)
}

// ItemFollow moves the given Moveable everytime the subscribed coordinates change
func ItemFollow(o Moveable, coords *pair.RxInt32PairSubscriber, offsetX, offsetY int32) {
	for pos := range coords.C {
		o.Move(pos.L-offsetX, pos.R-offsetY)
	}
}

// CoordFollow updates the l and r pointers when the coordinates get updated
func CoordFollow(coords *pair.RxInt32PairSubscriber, l, r *int32) {
	for pos := range coords.C {
		*l = pos.L
		*r = pos.R
	}
}

// EnableDragging runs dragging subprocesses for the given object when the boolean subscriber is true
func EnableDragging(o drag.Draggable, mouse *gamekit.Mouse, rxdrag *drag.RxDraggable, toDrag *rx.BoolSubscriber) {
	go func() {

		var mouseLocation *pair.RxInt32PairSubscriber

		// wait for button click state changes
		for dragOn := range toDrag.C {

			if dragOn {

				// if clicked, subscribe on the mouse cursor position and update the object
				// when the mouse cursor moves

				mouseLocation = mouse.Position.Subscribe()
				ml := mouse.Position.Get()

				curX, curY := o.Position()

				// offset ensures the relative location of the item compared to the mouse stays the mouse
				offsetX := ml.L - curX
				offsetY := ml.R - curY

				// signal that dragging has started
				rxdrag.Set(o)

				// Make the item follow the mouse location
				go ItemFollow(o, mouseLocation, offsetX, offsetY)

			} else if mouseLocation != nil {

				// when the button is not clicked, close the mouseLocation subscription
				mouseLocation.Close()
				mouseLocation = nil

				// signal that dragging has stopped
				rxdrag.Set(nil)
			}
		}
	}()
}

// EnableDropping runs dropping subproceses for the given slot.
func EnableDropping(s *Slot, mouse *gamekit.Mouse, rxdrag *drag.RxDraggable, isDestination bool) {
	go func() {

		sub := rxdrag.Subscribe()

		var curDragging drag.Draggable
		var mouseLocation *pair.RxInt32PairSubscriber

		var dragX, dragY int32

		for o := range sub.C {

			if o != nil {
				// dragging has started, track the mouse for collision
				mouseLocation = mouse.Position.Subscribe()

				// and save the current dragging object
				curDragging = o

				// make the dragX and dragY follow the mouse
				go CoordFollow(mouseLocation, &dragX, &dragY)

				continue
			}

			// dragging has stopped...

			// stop tracking the mouse position
			mouseLocation.Close()
			mouseLocation = nil

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
	}()
}

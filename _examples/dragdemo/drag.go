package main

import (
	"golang.org/x/net/context"

	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/pair"
	"github.com/sheenobu/rxgen/rx"
)

// Draggable defines an object that can be dragged
type Draggable interface {
	Position() (int32, int32)
	Move(int32, int32)
}

// EnableDragging runs dragging subprocesses for the given object when the boolean subscriber is true
func EnableDragging(o Draggable, mouse *gamekit.Mouse, clickBool *rx.BoolSubscriber) {
	go func() {
		var mouseLocation *pair.RxInt32PairSubscriber

		ctx, cancel := context.WithCancel(context.Background())

		// wait for button click state changes
		for buttonClick := range clickBool.C {

			if buttonClick {

				// if clicked, subscribe on the mouse cursor position and update the object
				// when the mouse cursor moves

				ctx, cancel = context.WithCancel(context.Background())

				var offsetX int32
				var offsetY int32

				ml := mouse.Position.Get()

				curX, curY := o.Position()

				// offset ensures the item we are dragging stays in relative position of the mouse from where we picked it up
				offsetX = ml.L - curX
				offsetY = ml.R - curY

				mouseLocation = mouse.Position.Subscribe()

				// the process which runs the Dragging by listening on mouse X and Y and tying it to the objects X and Y
				go func(mouseLocation *pair.RxInt32PairSubscriber, offsetX, offsetY int32) {
					defer mouseLocation.Close()
					for {
						select {
						case mouseMove := <-mouseLocation.C:
							o.Move(mouseMove.L-offsetX, mouseMove.R-offsetY)
						case <-ctx.Done():
							return
						}
					}
				}(mouseLocation, offsetX, offsetY)

			} else if mouseLocation != nil {

				// when the button is not clicked, close the mouseLocation subscription

				cancel()
			}
		}
	}()
}

package main

import (
	"github.com/sheenobu/go-gamekit"
	"github.com/sheenobu/go-gamekit/pair"
	"github.com/sheenobu/rxgen/rx"
	"golang.org/x/net/context"
)

// Draggable defines an object that can be dragged
type Draggable interface {
	Position() (int32, int32)
	Move(int32, int32)
}

// EnableDragging runs dragging subprocesses for the given object when the boolean subscriber is true
func EnableDragging(ctx context.Context, o Draggable, mouse *gamekit.Mouse, clickBool *rx.BoolSubscriber) {
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
			case mouseMove := <-mouseLocation:
				o.Move(mouseMove.L-offsetX, mouseMove.R-offsetY)
			case buttonClick := <-clickBool.C:

				if buttonClick {

					// if clicked, subscribe on the mouse cursor position and update the object
					// when the mouse cursor moves
					ml := mouse.Position.Get()

					curX, curY := o.Position()

					// offset ensures the item we are dragging stays in relative position of the mouse from where we picked it up
					offsetX = ml.L - curX
					offsetY = ml.R - curY

					sub := mouse.Position.Subscribe()
					mouseLocation = sub.C
					mouseLocationClose = sub.Close
				} else if mouseLocation != nil {
					// when the button is not clicked, close the mouseLocation subscription
					mouseLocationClose()
					mouseLocation = nil
					mouseLocationClose = func() {}
				}
			}
		}
	}()
}

package loop

import "github.com/sheenobu/go-gamekit"

type simpleLoop struct {
	wm      *gamekit.WindowManager
	running bool
	render  func()
}

// Simple builds a simple loop, with no timing management
func Simple(wm *gamekit.WindowManager, fn func()) Loop {
	return &simpleLoop{
		wm:      wm,
		running: false,
		render:  fn,
	}
}

// run the simple loop
func (sl *simpleLoop) Run() {
	sl.running = true

	go func() {
		sub := sl.wm.WindowCount.Subscribe()
		defer sub.Close()

		for {
			select {
			case c := <-sub.C:
				if c == 0 {
					sl.running = false
					return
				}
			}
		}
	}()

	for sl.running {
		sl.wm.DispatchEvents()
		sl.render()

		sdl.Delay(2)

	}
}

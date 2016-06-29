package main

import (
	"fmt"

	"github.com/sheenobu/go-gamekit"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"golang.org/x/net/context"
)

type resolutionPicker struct {
	font        *ttf.Font
	bounds      sdl.Rect
	resolutions []sdl.DisplayMode

	arrowUpButton   *button
	arrowDownButton *button

	arrowUp   *sdl.Texture
	arrowDown *sdl.Texture

	scrollPosition int32
	scrollMax      int32

	offscreenHeight  int32
	offscreenTexture *sdl.Texture

	selected int32
}

func newResolutionPicker(bounds sdl.Rect, r *sdl.Renderer) *resolutionPicker {

	rs := &resolutionPicker{}
	rs.bounds = bounds

	// load font
	f, err := ttf.OpenFont("./data/font.ttf", 15)
	if err != nil {
		panic(err)
	}

	rs.font = f

	// load arrow up texture
	rs.arrowUp, err = img.LoadTexture(r, "./data/arrow_up.png")
	if err != nil {
		panic(err)
	}

	// load arrow down texture
	rs.arrowDown, err = img.LoadTexture(r, "./data/arrow_down.png")
	if err != nil {
		panic(err)
	}

	// enumerate the native machine modes

	displays, err := sdl.GetNumVideoDisplays()
	if err != nil {
		panic(err)
	}

	for i := 0; i < displays; i++ {
		modes, err := sdl.GetNumDisplayModes(i)
		if err != nil {
			panic(err)
		}

		for j := 0; j < modes; j++ {
			var dm sdl.DisplayMode
			if err := sdl.GetDisplayMode(i, j, &dm); err != nil {
				panic(err)
			}

			rs.resolutions = append(rs.resolutions, dm)
		}
	}

	// add some common really small modes now

	rs.resolutions = append(rs.resolutions, sdl.DisplayMode{
		Format:      0,
		W:           1024,
		H:           768,
		RefreshRate: 60,
	})
	rs.resolutions = append(rs.resolutions, sdl.DisplayMode{
		Format:      0,
		W:           800,
		H:           600,
		RefreshRate: 60,
	})
	rs.resolutions = append(rs.resolutions, sdl.DisplayMode{
		Format:      0,
		W:           640,
		H:           480,
		RefreshRate: 60,
	})
	rs.resolutions = append(rs.resolutions, sdl.DisplayMode{
		Format:      0,
		W:           360,
		H:           240,
		RefreshRate: 60,
	})

	actualHeight := int32(len(rs.resolutions) * 20)
	rs.offscreenHeight = actualHeight

	rs.offscreenTexture, err = r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET,
		int(bounds.W), int(rs.offscreenHeight))

	if err != nil {
		panic(err)
	}

	r.SetRenderTarget(rs.offscreenTexture)

	r.SetDrawColor(0, 0, 0, 255)
	r.FillRect(&sdl.Rect{
		X: 0, Y: 0, W: rs.bounds.W, H: actualHeight})

	for idx, dm := range rs.resolutions {

		sfc, err := rs.font.RenderUTF8_Solid(fmt.Sprintf("%d/%d (%d)", dm.W, dm.H, dm.RefreshRate), sdl.Color{R: 255, G: 255, B: 255})
		if err != nil {
			panic(err)
		}
		t, err := r.CreateTextureFromSurface(sfc)
		if err != nil {
			panic(err)
		}

		r.Copy(t, nil, &sdl.Rect{X: 5, Y: int32(5 + idx*20), W: sfc.W, H: sfc.H})

		sfc.Free()
	}
	r.Present()
	r.SetRenderTarget(nil)

	rs.scrollPosition = 0
	rs.scrollMax = int32(len(rs.resolutions) - int(rs.bounds.H/20))

	rs.selected = 3

	// create the scroll buttons
	rs.arrowUpButton = newButton(&sdl.Rect{X: 235 * 2, Y: 4 * 2, W: 9 * 2, H: 7 * 2}, rs.arrowUp)
	rs.arrowDownButton = newButton(&sdl.Rect{X: 235 * 2, Y: 38 * 2, W: 9 * 2, H: 7 * 2}, rs.arrowDown)

	return rs
}

func (rs *resolutionPicker) Run(ctx context.Context, m *gamekit.Mouse, res *launchResults) {

	// set the selected display
	res.ChosenResolution = rs.resolutions[rs.selected]

	// run the button processes
	go rs.arrowUpButton.Run(ctx, m)
	go rs.arrowDownButton.Run(ctx, m)

	arrowUpSub := rs.arrowUpButton.Clicked.Subscribe()
	defer arrowUpSub.Close()

	arrowDownSub := rs.arrowDownButton.Clicked.Subscribe()
	defer arrowDownSub.Close()

	for {
		select {
		case clicked := <-arrowUpSub.C:
			if clicked && rs.scrollPosition > 0 {
				rs.scrollPosition--
			}
		case clicked := <-arrowDownSub.C:
			if clicked && rs.scrollPosition < rs.scrollMax {
				rs.scrollPosition++
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rs *resolutionPicker) Render(r *sdl.Renderer) {
	r.Copy(rs.offscreenTexture, &sdl.Rect{X: 0, Y: rs.scrollPosition * 20, W: rs.bounds.W, H: rs.bounds.H}, &rs.bounds)

	rs.arrowUpButton.Render(r)
	rs.arrowDownButton.Render(r)

	if rs.selected >= rs.scrollPosition {
		r.SetDrawColor(255, 255, 255, 255)
		r.DrawRect(&sdl.Rect{X: rs.bounds.X, Y: rs.bounds.Y + 20*(rs.selected-rs.scrollPosition), W: rs.bounds.W, H: 20})
	}

	scrollPos := float64(rs.scrollPosition) / float64(rs.scrollMax)

	r.SetDrawColor(255, 255, 255, 255)
	r.FillRect(&sdl.Rect{X: 236 * 2, Y: 11*2 + int32(((2*27)-(9*2))*scrollPos), W: 7 * 2, H: 9 * 2})
}

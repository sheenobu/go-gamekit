package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type sheet struct {
	texture   *sdl.Texture
	positions []sdl.Rect
}

func newSheet(r *sdl.Renderer, file string) *sheet {

	t, err := img.LoadTexture(r, file)
	if err != nil {
		panic(err)
	}

	return &sheet{
		texture:   t,
		positions: make([]sdl.Rect, 0),
	}
}

func (s *sheet) Add(r *sdl.Rect) (idx int) {
	s.positions = append(s.positions, *r)
	return len(s.positions) - 1
}

func (s *sheet) Copy(r *sdl.Renderer, dest *sdl.Rect, idx int) error {
	dest.W = s.positions[idx].W * 2
	dest.H = s.positions[idx].H * 2
	return r.Copy(s.texture, &s.positions[idx], dest)
}

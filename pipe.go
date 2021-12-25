package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	mu       sync.RWMutex
	x        int32
	h        int32
	speed    int32
	inverted bool
	w        int32

	texture *sdl.Texture
}

func newPipe(r *sdl.Renderer) (*pipe, error) {
	texture, err := img.LoadTexture(r, "res/OldPipe.webp")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	return &pipe{
		x:        400,
		w:        52,
		h:        300,
		speed:    8,
		texture:  texture,
		inverted: true,
	}, nil
}

func (p *pipe) destroy() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.texture.Destroy()
}

func (p *pipe) paint(r *sdl.Renderer) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{W: p.w, H: p.h, X: p.x, Y: (600 - p.h)}
	if p.inverted {
		rect.Y = 0
	}

	if p.inverted {
		if err := r.CopyEx(p.texture, nil, rect, 0, nil, sdl.FLIP_VERTICAL); err != nil {
			return fmt.Errorf("could not draw pipe")
		}
	} else {
		if err := r.Copy(p.texture, nil, rect); err != nil {
			return fmt.Errorf("could not draw pipe")
		}
	}

	return nil
}

func (p *pipe) restart() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x = 400
}

func (p *pipe) update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= p.speed
}

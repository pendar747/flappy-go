package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	speed   int32

	pipes []*pipe
}

type pipe struct {
	mu       sync.RWMutex
	x        int32
	h        int32
	w        int32
	inverted bool
}

func newPipes(r *sdl.Renderer) (ps *pipes, err error) {
	texture, err := img.LoadTexture(r, "res/OldPipe.webp")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	ps = &pipes{
		texture: texture,
		speed:   8,
	}

	go func() {
		for {
			pipe := newPipe()
			ps.pipes = append(ps.pipes, pipe)
			time.Sleep(time.Second * 2)
		}
	}()

	return ps, nil
}

func newPipe() *pipe {
	return &pipe{
		x:        800,
		w:        52,
		h:        rand.Int31n(400),
		inverted: rand.Int31()%2 == 0,
	}
}

func (ps *pipes) destroy() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	return ps.texture.Destroy()
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.mu.RLock()

		rect := &sdl.Rect{W: p.w, H: p.h, X: p.x, Y: (600 - p.h)}
		if p.inverted {
			rect.Y = 0
		}

		if p.inverted {
			if err := r.CopyEx(ps.texture, nil, rect, 0, nil, sdl.FLIP_VERTICAL); err != nil {
				return fmt.Errorf("could not draw pipe")
			}
		} else {
			if err := r.Copy(ps.texture, nil, rect); err != nil {
				return fmt.Errorf("could not draw pipe")
			}
		}

		p.mu.RUnlock()
	}

	return nil
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, p := range ps.pipes {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.x -= ps.speed
	}
}

func (ps *pipes) touch(b *bird) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, p := range ps.pipes {
		b.touch(p)
	}
}

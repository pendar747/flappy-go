package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/jpg/Mountain1.jpg")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	birds := make([]*sdl.Texture, 4)
	for i := 0; i < 4; i++ {
		bird, err := img.LoadTexture(r, fmt.Sprintf("res/PNG/Frame-%d.png", i+1))
		if err != nil {
			return nil, fmt.Errorf("could not load image: %v", err)
		}
		birds[i] = bird
	}

	return &scene{bg: bg, birds: birds}, nil
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not draw background")
	}

	rect := &sdl.Rect{W: 50, H: 43, X: 10, Y: 300 - 43/2}

	i := s.time % len(s.birds)
	if err := r.Copy(s.birds[i], nil, rect); err != nil {
		return fmt.Errorf("could not draw bird")
	}

	r.Present()
	return nil
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(100 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- errors.Wrap(err, "could not paint scene")
				}
			}
		}
	}()

	return errc
}

func (s *scene) destroy() error {
	return s.bg.Destroy()
}

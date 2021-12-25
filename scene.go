package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg   *sdl.Texture
	bird *bird
	pipe *pipe
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/jpg/Mountain1.jpg")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}

	pipe, err := newPipe(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: bird, pipe: pipe}, nil
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not draw background")
	}

	err := s.bird.paint(r)
	if err != nil {
		return err
	}

	err = s.pipe.paint(r)
	if err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(100 * time.Millisecond)
		done := false
		for !done {
			select {
			case e := <-events:
				done = s.handleEvent(e)
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game over")
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- errors.Wrap(err, "could not paint scene")
				}
			}
		}
	}()

	return errc
}

func (s *scene) update() {
	s.bird.update()
	s.pipe.update()
	s.bird.touch(s.pipe)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipe.restart()
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
		return false
	default:
		log.Printf("unknown event %T", e)
	}
	return false
}

func (s *scene) destroy() error {
	err := s.bg.Destroy()
	if err != nil {
		return err
	}

	err = s.bird.destroy()
	if err != nil {
		return err
	}

	err = s.pipe.destroy()
	if err != nil {
		return err
	}

	return nil
}

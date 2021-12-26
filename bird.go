package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const gravity = 0.1
const jumpSpeed = 3

type bird struct {
	time     int
	textures []*sdl.Texture
	dead     bool
	mu       sync.RWMutex

	x, y, w, h int32

	speed float64
}

func newBird(r *sdl.Renderer) (*bird, error) {
	textures := make([]*sdl.Texture, 4)
	for i := 0; i < 4; i++ {
		texture, err := img.LoadTexture(r, fmt.Sprintf("res/PNG/Frame-%d.png", i+1))
		if err != nil {
			return nil, fmt.Errorf("could not load image: %v", err)
		}
		textures[i] = texture
	}

	return &bird{textures: textures, y: 300, speed: 1, x: 10, w: 50, h: 43}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	b.speed += gravity
	if b.y < 43 {
		b.dead = true
	}
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	rect := &sdl.Rect{W: b.w, H: b.h, X: b.x, Y: (600 - int32(b.y) - b.h)}

	i := (b.time / 10) % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not draw bird")
	}

	return nil
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = -1 * jumpSpeed
}

func (b *bird) destroy() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, texture := range b.textures {
		err := texture.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *bird) isDead() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.dead
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	p.mu.RLock()
	defer p.mu.RUnlock()

	heightConstraint := b.y+b.h < p.h
	if p.inverted {
		heightConstraint = (600 - b.y + b.h) < p.h
	}
	if b.x+b.w > p.x && heightConstraint && b.x+b.w < p.x+p.w {
		log.Printf("bird touched the pipe!!")
		b.dead = true
	}
}

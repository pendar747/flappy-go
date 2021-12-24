package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const gravity = 3
const jumpSpeed = 10

type bird struct {
	time     int
	textures []*sdl.Texture

	y, speed float64
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

	return &bird{textures: textures, y: 300, speed: 1}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	b.speed += gravity
	if b.y < 43 {
		b.speed = -10
		b.y = 43
	}

	rect := &sdl.Rect{W: 50, H: 43, X: 10, Y: (600 - int32(b.y))}

	i := b.time % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not draw bird")
	}

	return nil
}

func (b *bird) jump() {
	b.speed = -1 * jumpSpeed
}

func (b *bird) destroy() {
	for _, texture := range b.textures {
		texture.Destroy()
	}
}

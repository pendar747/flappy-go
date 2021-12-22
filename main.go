package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

func keepAlive(duration *time.Duration) {
	start := time.Now()
	running := true
	for running {
		if duration != nil && time.Since(start) > *duration {
			fmt.Printf("time passed: %s > %s\n", time.Since(start), *duration)
			running = false
			break
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	// w2, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
	// 	800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer window.Destroy()

	if err := drawTitle(renderer); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	d := 2 * time.Second
	keepAlive(&d)

	s, err := newScene(renderer)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	time.AfterFunc(5*time.Second, func() {
		cancel()
	})

	return <-s.run(ctx, renderer)
}

func drawTitle(renderer *sdl.Renderer) error {
	renderer.Clear()

	font, err := ttf.OpenFont("res/fonts/YujiHentaiganaAkari-Regular.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}

	s, err := font.RenderUTF8Solid("Flappy Gopher", sdl.Color{
		R: 73,
		G: 229,
		B: 156,
	})
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}

	texture, err := renderer.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	if err := renderer.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("could not render texture: %v", err)
	}

	renderer.Present()

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in run: %v", err)
		os.Exit(2)
	}
}

package main

import (
	"fmt"
	"log"
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

	if err := drawTitle(renderer, "Flappy Gopher"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	d := 2 * time.Second
	keepAlive(&d)

	s, err := newScene(renderer)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	events := make(chan sdl.Event)
	errc := s.run(events, renderer)

	_ = errc

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			events <- event
			log.Printf("received event is %T\n", event)
			// TODO: fix bug with handling error
			// if err := <-errc; err != nil {
			// log.Printf("error rendering scene: %s", err.Error())
			// running = false
			// return err
			// }
			switch event.(type) {
			case *sdl.QuitEvent:
				log.Print("Quit\n")
				running = false
			}
		}
	}

	return nil
}

func drawTitle(renderer *sdl.Renderer, text string) error {
	renderer.Clear()

	font, err := ttf.OpenFont("res/fonts/YujiHentaiganaAkari-Regular.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}

	s, err := font.RenderUTF8Solid(text, sdl.Color{
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

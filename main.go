package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600

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

	window, renderer, err := sdl.CreateWindowAndRenderer(WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
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

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()
	texture, textW, textH, err := getTextTexture(r, text)
	if err != nil {
		return err
	}
	textRect := &sdl.Rect{
		X: WINDOW_WIDTH/2 - textW/2,
		Y: WINDOW_HEIGHT/2 - textH/2,
		W: textW,
		H: textH,
	}
	err = r.Copy(texture, nil, textRect)
	r.Present()
	return err
}

func drawGameOver(r *sdl.Renderer) error {
	r.Clear()
	s, err := img.Load("res/PNG/Frame-3.png")
	if err != nil {
		return fmt.Errorf("could not load surface %v", err)
	}
	texture, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not load image: %v", err)
	}

	w := int32(200)
	h := w * (s.W / s.H)
	rect := &sdl.Rect{
		W: w,
		H: h,
		X: WINDOW_WIDTH/2 - w/2,
		Y: WINDOW_HEIGHT/3 - h/2,
	}
	err = r.CopyEx(texture, nil, rect, 0, nil, sdl.FLIP_VERTICAL)
	if err != nil {
		return fmt.Errorf("could not render game over image: %v", err)
	}

	texture, textW, textH, err := getTextTexture(r, "Game over")
	if err != nil {
		return err
	}
	textRect := &sdl.Rect{
		X: WINDOW_WIDTH/2 - textW/2,
		Y: rect.Y + rect.H + 10,
		W: textW,
		H: textH,
	}
	err = r.Copy(texture, nil, textRect)
	r.Present()
	return err
}

func getTextTexture(r *sdl.Renderer, text string) (texture *sdl.Texture, w, h int32, err error) {
	font, err := ttf.OpenFont("res/fonts/Pushster-Regular.ttf", 120)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("could not open font: %v", err)
	}

	s, err := font.RenderUTF8Blended(text, sdl.Color{
		R: 73,
		G: 229,
		B: 156,
	})
	if err != nil {
		return nil, 0, 0, fmt.Errorf("could not render title: %v", err)
	}

	texture, err = r.CreateTextureFromSurface(s)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("could not create texture: %v", err)
	}

	return texture, s.W, s.H, nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in run: %v", err)
		os.Exit(2)
	}
}

package main

import (
	"fmt"
	"log"
	"github.com/veandco/go-sdl2/sdl"
	"doom-clone/internal"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	fmt.Println("Game is starting...")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("Error initializing SDL: %s", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Doom-like Game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Error creating window: %s", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Error creating renderer: %s", err)
	}
	defer renderer.Destroy()

	// Start the game loop
	if err := internal.RunGame(renderer); err != nil {
		log.Fatalf("Error running the game: %s", err)
	}
}


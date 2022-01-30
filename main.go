package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetFullscreen(true)
	w, h := ebiten.ScreenSizeInFullscreen()
	toRun := &Game{
		sprites: newSprites(),
		width:   w / 2,
		height:  h / 2,
	}
	toRun.reset()
	if err := ebiten.RunGame(toRun); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

const (
	starChance = 0.0005
)

func newStars(width, height int) []iVec2 {
	var s []iVec2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Float64() < starChance {
				s = append(s, iVec2{x, y})
			}
		}
	}
	return s
}

func (g *Game) updateBackgroundStars() {
	for i, star := range g.bgStars {
		layer := 3 + (star.x % 3)
		if g.tick%layer == 0 {
			star.y++
			if star.y == g.height {
				star.y = 0
			}
			g.bgStars[i] = star
		}
	}
}

func (g *Game) drawBackgroundStars(screen *ebiten.Image) {
	for _, star := range g.bgStars {
		opts := &ebiten.DrawImageOptions{}
		g.drawSprite(screen, g.sprites.star, opts, fVec2{float64(star.x), float64(star.y)}, fVec2{1, 1})
	}
}

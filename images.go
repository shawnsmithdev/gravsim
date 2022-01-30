package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

type sprites struct {
	ship   *ebiten.Image
	planet *ebiten.Image
	hole   *ebiten.Image
	star   *ebiten.Image
}

const (
	shipSpriteSize   = 10
	planetSpriteSize = 128
	planetSpriteR2   = planetSpriteSize * planetSpriteSize / 4
)

// its a blue triangle
func newShipSprite() *ebiten.Image {
	pixels := make([]byte, 4*shipSpriteSize*shipSpriteSize)
	for y := 0; y < shipSpriteSize; y++ {
		ybase := y * shipSpriteSize * 4
		for x := 0; x < shipSpriteSize; x++ {
			xbase := ybase + (x * 4)
			if y > 2*x-shipSpriteSize && y > -2*x+shipSpriteSize {
				pixels[xbase] = 64
				pixels[xbase+1] = 64
				pixels[xbase+2] = 255
				pixels[xbase+3] = 255
			}
		}
	}
	result := ebiten.NewImage(shipSpriteSize, shipSpriteSize)
	result.ReplacePixels(pixels)
	return result
}

// its a green circle
func newPlanetSprite() *ebiten.Image {
	mid := planetSpriteSize / 2
	pixels := make([]byte, 4*planetSpriteSize*planetSpriteSize)
	for y := 0; y < planetSpriteSize; y++ {
		ybase := y * planetSpriteSize * 4
		for x := 0; x < planetSpriteSize; x++ {
			xbase := ybase + (x * 4)
			cx := x - mid
			cy := y - mid
			if cx*cx+cy*cy <= planetSpriteR2 {
				color := colornames.Darkgreen
				if cx*cx+cy*cy <= 16 {
					color = colornames.Darkolivegreen
				}
				pixels[xbase] = color.R
				pixels[xbase+1] = color.G
				pixels[xbase+2] = color.B
				pixels[xbase+3] = 255
			}
		}
	}
	result := ebiten.NewImage(planetSpriteSize, planetSpriteSize)
	result.ReplacePixels(pixels)
	return result
}

func newStarSprite() *ebiten.Image {
	result := ebiten.NewImage(1, 1)
	result.Set(0, 0, colornames.White)
	return result
}

func newBlackHoleSprite() *ebiten.Image {
	mid := planetSpriteSize / 2
	pixels := make([]byte, 4*planetSpriteSize*planetSpriteSize)
	for y := 0; y < planetSpriteSize; y++ {
		ybase := y * planetSpriteSize * 4
		for x := 0; x < planetSpriteSize; x++ {
			xbase := ybase + (x * 4)
			cx := x - mid
			cy := y - mid
			r2 := cx*cx + cy*cy
			if r2 <= planetSpriteR2 {
				z := math.Sqrt(float64(r2)) / (planetSpriteSize / 2)
				copy(pixels[xbase:xbase+4], []byte{
					32 - byte(32*z),
					0,
					32 - byte(32*z),
					255 - byte(255*z),
				})
			}
		}
	}
	result := ebiten.NewImage(planetSpriteSize, planetSpriteSize)
	result.ReplacePixels(pixels)
	return result
}

func newSprites() sprites {
	return sprites{
		ship:   newShipSprite(),
		planet: newPlanetSprite(),
		hole:   newBlackHoleSprite(),
		star:   newStarSprite(),
	}
}

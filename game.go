package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"math"
	"math/rand"
	"os"
)

const (
	defaultPlanetCount = 3
	maxPlanets         = 255
	levelClearBonus    = 250000
	scorePerExtraLife  = 1000000
)

type Game struct {
	width, height int
	planets       []*planet // planets with gravity
	player        *ship     // player controlled ship TODO: multiplayer
	sprites       sprites   // sprite cache
	bgStars       []iVec2
	hideStars     bool
	tick          int // used for animations
	lives         int // remaining lives
	score         uint64
	hiScore       uint64
}

func (g *Game) nearestVirtual(here, there fVec2) fVec2 {
	var dx, dy int
	var result fVec2
	if here.x < float64(g.width)/2 {
		dx = -1
	}
	if here.y < float64(g.height)/2 {
		dy = -1
	}
	smallest := math.Inf(1)
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			virtPos := there.add(fVec2{
				float64((x + dx) * g.width),
				float64((y + dy) * g.height),
			})
			mag := here.add(virtPos.negate()).magnitude()
			if mag < smallest {
				smallest = mag
				result = virtPos
			}
		}
	}
	return result
}

// finds if a point collides with a planet, returns planet index or -1 if no collision
func (g *Game) findCollision(pos fVec2) int {
	for i, here := range g.planets {
		diff := pos.add(g.nearestVirtual(pos, here.pos).negate())
		if diff.magnitude() < here.radius {
			return i
		}
	}
	return -1
}

// (bullet, planet)
func (g *Game) bulletCollides() (int, int) {
	for i, b := range g.player.bullets {
		found := g.findCollision(b.pos)
		if found >= 0 {
			return i, found
		}
	}
	return -1, -1
}

func (g *Game) movePlayer() {
	g.player.angPos = math.Pi
	g.player.accel.reset()
	g.player.vel.reset()
	g.player.bullets = nil
	var collision int
	for collision >= 0 {
		g.player.pos = g.randomPos()
		collision = g.findCollision(g.player.pos)
	}
}

func (g *Game) Update() error {
	g.tick++

	// key input
	g.updateKeyboard()

	// player ship angle
	g.player.updateAngle()

	// accel due to thruster
	g.player.accel.reset()
	g.player.applyThrust()

	// accel due to gravity
	for x, here := range g.planets {
		here.accel.reset()

		// ship gravity on planet
		virtualPos := g.nearestVirtual(here.pos, g.player.pos)
		here.gravityFor(&state{
			pos:  virtualPos,
			mass: g.player.mass,
		})

		// planet-planet gravity
		for y, there := range g.planets {
			if x != y {
				virtualPos = g.nearestVirtual(here.pos, there.pos)
				here.gravityFor(&state{
					pos:  virtualPos,
					mass: there.mass,
				})
			}
		}

		// planet gravity on ship
		virtualPos = g.nearestVirtual(g.player.pos, here.pos)
		g.player.gravityFor(&state{
			pos:  virtualPos,
			mass: here.mass,
		})
	}

	// velocity/position
	// planets
	for _, here := range g.planets {
		here.updateAngle()
		here.updateVelocityBound()
		here.updatePosition(g.width, g.height)
	}

	// ship
	g.player.updateVelocityBound()
	g.player.updatePosition(g.width, g.height)
	g.player.updateBullets(g.width, g.height)

	// stars
	g.updateBackgroundStars()

	// collisions
	if len(g.planets) > 0 {
		playerCollision := g.findCollision(g.player.pos)
		if playerCollision >= 0 {
			g.lives--
			if g.lives < 0 {
				g.reset()
				return nil
			}
			g.movePlayer()
		}

		bulletIdx, planetIdx := g.bulletCollides()
		if bulletIdx >= 0 {
			g.addScore(uint(g.planets[planetIdx].mass))
			g.destroyPlanet(planetIdx)
			if len(g.planets) == 0 {
				g.addScore(levelClearBonus)
				g.nextLevel()
			} else {
				g.player.removeBullet(bulletIdx)
			}
		}
	}

	// planet-planet
	if len(g.planets) > 1 {
		for x, here := range g.planets[:len(g.planets)-1] {
			for _, there := range g.planets[1+x:] {
				// continue if no collision
				diff := here.pos.add(g.nearestVirtual(here.pos, there.pos).negate())
				dist := diff.magnitude()
				nudge := here.radius + there.radius - dist
				if nudge < 0 {
					continue
				}

				// unit vector normal to collision
				normal := diff.scale(1.0 / dist)

				// move out of collision
				here.pos = here.pos.add(normal.scale(nudge))
				there.pos = there.pos.add(normal.scale(-nudge))

				// change normal part of velocity for elastic collision
				massDiff := here.mass - there.mass
				both := here.mass + there.mass
				hereNormal := here.vel.dotProduct(normal)
				thereNormal := there.vel.dotProduct(normal)
				// v1' = v1(m1-m2)/(m1+m2) + v2(2*m2)(m1+m2)
				hereReflect := normal.scale((hereNormal * massDiff / both) + (thereNormal * 2 * there.mass / both))
				// v2' = v1(2*m1)/(m1+m2) + v2(m2-m1)(m1+m2)
				thereReflect := normal.scale((hereNormal * 2 * here.mass / both) + (thereNormal * -massDiff / both))

				// a unit vector orthogonal to the collision, velocity in this direction stays the same
				orth := fVec2{-normal.y, normal.x}
				hereOrth := orth.scale(here.vel.dotProduct(orth))
				thereOrth := orth.scale(there.vel.dotProduct(orth))

				here.vel = hereReflect.add(hereOrth)
				there.vel = thereReflect.add(thereOrth)
			}
		}
	}
	return nil
}

func (g *Game) drawSprite(screen *ebiten.Image, sprite *ebiten.Image, opts *ebiten.DrawImageOptions, here, size fVec2) {
	// draw sprite at real location
	opts.GeoM.Translate(here.x, here.y)
	defer opts.GeoM.Translate(-here.x, -here.y)
	screen.DrawImage(sprite, opts)

	// draw around border wrap
	var dx, dy float64
	if here.x < size.x {
		dx = 1
	} else if here.x > float64(g.width)-size.x {
		dx = -1
	}
	if here.y < size.y {
		dy = 1
	} else if here.y > float64(g.height)-size.y {
		dy = -1
	}

	if dx == 0 && dy == 0 {
		return
	}

	if dx != 0 {
		opts.GeoM.Translate(dx*float64(g.width), 0)
		screen.DrawImage(sprite, opts)
		opts.GeoM.Translate(-dx*float64(g.width), 0)
	}

	if dy != 0 {
		opts.GeoM.Translate(0, dy*float64(g.height))
		screen.DrawImage(sprite, opts)
		if dx != 0 {
			opts.GeoM.Translate(dx*float64(g.width), 0)
			screen.DrawImage(sprite, opts)
			opts.GeoM.Translate(-dx*float64(g.width), 0)
		}
		opts.GeoM.Translate(0, -dy*float64(g.height))
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// stars
	if !g.hideStars {
		g.drawBackgroundStars(screen)
	}

	// planets
	for _, here := range g.planets {
		scale := here.mass / planetScaleToMass
		planetOpts := &ebiten.DrawImageOptions{}
		planetOpts.GeoM.Scale(scale, scale)
		planetOpts.GeoM.Translate(-scale*planetSpriteSize/2, -scale*planetSpriteSize/2)
		planetOpts.GeoM.Rotate(here.angPos)
		g.drawSprite(screen, g.sprites.planet, planetOpts, here.pos, fVec2{here.radius, here.radius})
	}

	// Text
	// Usage
	text.Draw(screen, "Quit:Q Reset:R, Thrust:UP, Turn:LEFT/RIGHT, Show/Hide Stars:S, Stop:F", basicfont.Face7x13,
		(g.width/2)-220, g.height-4,
		colornames.Limegreen)
	// Score
	text.Draw(screen, fmt.Sprintf("Top Score: %09d", g.hiScore), basicfont.Face7x13,
		4, 12, colornames.Limegreen)
	text.Draw(screen, fmt.Sprintf("Score: %09d", g.score), basicfont.Face7x13,
		(g.width/2)-60, 12, colornames.Limegreen)
	// Lives
	if g.lives > 0 {
		lifeOpts := &ebiten.DrawImageOptions{}
		lifeOpts.GeoM.Translate(float64(g.width)-42, 3)
		if g.lives > 3 {
			text.Draw(screen, fmt.Sprintf("x%02d", g.lives),
				basicfont.Face7x13,
				g.width-32, 12, colornames.Limegreen)
			screen.DrawImage(g.sprites.ship, lifeOpts)
		} else {
			for i := 0; i < g.lives; i++ {
				screen.DrawImage(g.sprites.ship, lifeOpts)
				lifeOpts.GeoM.Translate(shipSpriteSize+4, 0)
			}
		}
	}

	// ship
	shipOpts := &ebiten.DrawImageOptions{}
	shipOpts.GeoM.Translate(-shipSpriteSize/2, -shipSpriteSize/2)
	shipOpts.GeoM.Rotate(math.Pi - g.player.angPos)
	g.drawSprite(screen, g.sprites.ship, shipOpts, g.player.pos, fVec2{shipSpriteSize, shipSpriteSize})

	// ship thuster
	thrustColor := colornames.Cyan
	if g.player.thrust > 0 {
		thrustColor = colornames.Red
	}
	ebitenutil.DrawLine(screen,
		g.player.pos.x+(-shipSpriteSize*math.Sin(g.player.angPos)),
		g.player.pos.y+(-shipSpriteSize*math.Cos(g.player.angPos)),
		g.player.pos.x+(-shipSpriteSize*math.Sin(g.player.angPos)/2),
		g.player.pos.y+(-shipSpriteSize*math.Cos(g.player.angPos)/2),
		thrustColor)

	// bullets
	for _, b := range g.player.bullets {
		ebitenutil.DrawLine(screen,
			b.pos.x,
			b.pos.y,
			b.pos.x+(bulletLineLen*math.Sin(b.angPos)),
			b.pos.y+(bulletLineLen*math.Cos(b.angPos)),
			colornames.Red)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func (g *Game) updateKeyboard() {
	// quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
	// reset
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.reset()
	}
	// hide stars performance?
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.hideStars = !g.hideStars
	}
	// thrust
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.player.thrust = shipThrustAccel
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		g.player.thrust = 0
	}
	// rotate
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.player.startTurn(shipTurnAngVel)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.player.startTurn(-shipTurnAngVel)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) || inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		g.player.stopTurn()
	}
	// full stop
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.player.fullStop()
	}
	// shoot bullets
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.player.shoot()
	}
}

func (g *Game) randomPos() fVec2 {
	return fVec2{
		x: float64(rand.Intn(g.width)),
		y: float64(rand.Intn(g.height)),
	}
}

func (g *Game) addPlanet() {
	if len(g.planets) > maxPlanets {
		return
	}
	np := newPlanet()
	np.pos = g.randomPos()
	np.angVel = (rand.Float64() * 0.1) - 0.05
	g.planets = append(g.planets, np)
}

func (g *Game) destroyPlanet(planet int) {
	if planet < len(g.planets)-1 {
		copy(g.planets[planet:], g.planets[planet+1:])
	}
	g.planets = g.planets[:len(g.planets)-1]
}

func (g *Game) addScore(amount uint) {
	before := g.score / scorePerExtraLife
	now := (g.score + uint64(amount)) / scorePerExtraLife
	if now > before {
		g.lives += int(now - before)
	}
	g.score += uint64(amount)
	if g.score > g.hiScore {
		g.hiScore = g.score
	}
}

func (g *Game) reset() {
	if g.player == nil {
		g.player = newShip()
	}
	g.tick = 1
	g.lives = 2
	g.score = 0
	g.nextLevel()
}

func (g *Game) nextLevel() {
	g.planets = nil
	for i := 0; i < defaultPlanetCount; i++ {
		g.addPlanet()
	}
	g.movePlayer()
	g.bgStars = newStars(g.width, g.height)
}

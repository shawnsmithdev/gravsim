package main

import (
	"math"
)

const (
	shipMass        = 200
	shipTurnAngVel  = 0.08
	shipThrustAccel = 0.04
	maxShipVel      = 15

	maxBullets    = 3
	bulletVel     = 6
	bulletTicks   = 30
	bulletLineLen = 6
)

// bullet can destroy planets, eventually expires, and has no mass
type bullet struct {
	state
	ticksLeft int
}

// the ship is the player controlled object
type ship struct {
	state
	thrust  float64
	bullets []*bullet
}

func newShip() *ship {
	return &ship{
		state: state{
			angPos: math.Pi,
			mass:   shipMass,
		},
	}
}

func (s *ship) updateVelocityBound() {
	s.state.updateVelocity(maxShipVel)
}

func (s *ship) applyThrust() {
	if s.thrust != 0 {
		s.accel.x += s.thrust * math.Sin(s.angPos)
		s.accel.y += s.thrust * math.Cos(s.angPos)
	}
}

func (s *ship) shoot() {
	if len(s.bullets) < maxBullets {
		s.bullets = append(s.bullets, &bullet{
			state: state{
				angPos: s.angPos,
				pos:    s.pos,
				vel: fVec2{
					x: bulletVel * math.Sin(s.angPos),
					y: bulletVel * math.Cos(s.angPos),
				},
			},
			ticksLeft: bulletTicks,
		})
	}
}

// find and remove expired bullets, and update the rest
func (s *ship) updateBullets(boundx, boundy int) {
	var expired []int
	for i, b := range s.bullets {
		b.ticksLeft--
		if b.ticksLeft > 0 {
			b.updatePosition(boundx, boundy)
		} else {
			expired = append(expired, i)
		}
	}

	for i, x := range expired {
		s.removeBullet(i + x)
	}
}

func (s *ship) removeBullet(bullet int) {
	if bullet < len(s.bullets)-1 {
		copy(s.bullets[bullet:], s.bullets[bullet+1:])
	}
	s.bullets = s.bullets[:len(s.bullets)-1]
}

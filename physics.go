package main

import (
	"math"
)

const (
	bigG = 0.01
)

// handle torus behavior at edges
func boundPos(x, max float64) float64 {
	if x < 0 {
		return x + math.Ceil(-x/max)*max
	}
	if x > max {
		return x - math.Floor(x/max)*max
	}
	return x
}

// state is the basic physics state of all game objects
type state struct {
	pos    fVec2
	vel    fVec2
	accel  fVec2
	mass   float64
	angPos float64 // rads
	angVel float64 // rads per tick
}

func (s *state) startTurn(delta float64) {
	s.angVel += delta
}

func (s *state) stopTurn() {
	s.angVel = 0
}

func (s *state) fullStop() {
	s.stopTurn()
	s.vel = fVec2{0, 0}
}

func (s *state) updatePosition(boundx, boundy int) {
	s.pos.x = boundPos(s.pos.x+s.vel.x, float64(boundx))
	s.pos.y = boundPos(s.pos.y+s.vel.y, float64(boundy))
}

func (s *state) updateVelocity(bound float64) {
	s.vel = s.vel.add(s.accel)
	speed := s.vel.magnitude()
	if speed > bound {
		s.vel = s.vel.scale(bound / speed)
	}
}

func (s *state) updateAngle() {
	s.angPos += s.angVel
}

// gravity updates accel of t to reflect gravitational attraction of other
func (s *state) gravityFor(other *state) {
	if s.mass == 0 {
		return
	}
	r := fVec2{
		x: s.pos.x - other.pos.x,
		y: s.pos.y - other.pos.y,
	}
	rMag := r.magnitude()
	if rMag == 0 {
		return
	}
	rAccel := -bigG * other.mass / (rMag * rMag)
	s.accel.x += rAccel * (r.x) / rMag
	s.accel.y += rAccel * (r.y) / rMag
}

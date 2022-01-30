package main

import "math/rand"

const (
	maxPlanetVelocity = 5
	radiusPerMass     = 300

	minPlanetMass      = 3000.0
	maxExtraPlanetMass = 15000.0
	planetScaleToMass  = 20000.0
)

// a planet is just a circle for now
type planet struct {
	state
	radius float64
}

func (p *planet) updateVelocityBound() {
	p.state.updateVelocity(maxPlanetVelocity)
}

func newPlanet() *planet {
	mass := minPlanetMass + (maxExtraPlanetMass * rand.Float64())
	return &planet{
		state: state{
			vel:  fVec2{rand.Float64() - 0.5, rand.Float64() - 0.5},
			mass: mass,
		},
		radius: mass / radiusPerMass,
	}
}

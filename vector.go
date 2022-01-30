package main

import "math"

// 2d vectors
type fVec2 struct {
	x, y float64
}

func (v *fVec2) reset() {
	v.x, v.y = 0, 0
}

func (v fVec2) negate() fVec2 {
	return fVec2{-v.x, -v.y}
}

func (v fVec2) add(other fVec2) fVec2 {
	return fVec2{v.x + other.x, v.y + other.y}
}

func (v fVec2) scale(scaleBy float64) fVec2 {
	return fVec2{v.x * scaleBy, v.y * scaleBy}
}

func (v fVec2) dotProduct(other fVec2) float64 {
	return v.y*other.y + v.x*other.x
}

func (v fVec2) magnitude() float64 {
	return math.Sqrt(v.dotProduct(v))
}

type iVec2 struct {
	x, y int
}

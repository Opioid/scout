package math

import (
	"github.com/Opioid/math32"
)

type OptimizedRay struct {
	Ray
	Time float32	
	Depth uint32	
	ReciprocalDirection Vector3
	Sign [3]uint32
}

func MakeOptimizedRay(origin, direction Vector3, mint, maxt float32, time float32, depth uint32) OptimizedRay {
	r := OptimizedRay{}
	r.Origin = origin
	r.SetDirection(direction)
	r.MinT = mint
	r.MaxT = maxt
	r.Time = time
	r.Depth = depth
	return r
}

func (r *OptimizedRay) SetDirection(direction Vector3) {
	r.Direction = direction
	r.ReciprocalDirection = MakeVector3(1.0 / direction.X, 1.0 / direction.Y, 1.0 / direction.Z)

	r.Sign[0] = math32.Signbit(direction.X)
	r.Sign[1] = math32.Signbit(direction.Y)
	r.Sign[2] = math32.Signbit(direction.Z)
}
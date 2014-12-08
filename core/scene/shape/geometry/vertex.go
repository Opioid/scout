package geometry

import (
	"github.com/Opioid/scout/base/math"
)

type Vertex struct {
	P, N, T math.Vector3
	BitangentSign float32
	UV math.Vector2
}
package geometry

import (
	"github.com/Opioid/scout/base/math"
)

type Differential struct {
	P math.Vector3			// posisition in world space
	T, B, N math.Vector3	// tangent frame in world space
	UV math.Vector2			// texture coordinates
}
package geometry

import (
	"github.com/Opioid/scout/base/math"
)

type Differential struct {
	P math.Vector3			// posisition in world space
	T, B, N math.Vector3	// tangent frame in world space
	UV math.Vector2			// texture coordinates
}

func (dg *Differential) TangentToWorld(v math.Vector3) math.Vector3 {
	return math.MakeVector3(
		v.X * dg.T.X + v.Y * dg.B.X + v.Z * dg.N.X,
		v.X * dg.T.Y + v.Y * dg.B.Y + v.Z * dg.N.Y,
		v.X * dg.T.Z + v.Y * dg.B.Z + v.Z * dg.N.Z)
}
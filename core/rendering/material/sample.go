package material

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Sample interface {
	Evaluate(l math.Vector3) math.Vector3
	Values() *Values

	SampleEvaluate(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32)

	TangentToWorld(v math.Vector3) math.Vector3

	CoordinateSystem() (math.Vector3, math.Vector3, math.Vector3)
}

type SampleBase struct {
	T, B, N math.Vector3

	Wo math.Vector3
}

func (s *SampleBase) TangentToWorld(v math.Vector3) math.Vector3 {
	return math.MakeVector3(
		v.X * s.T.X + v.Y * s.B.X + v.Z * s.N.X,
		v.X * s.T.Y + v.Y * s.B.Y + v.Z * s.N.Y,
		v.X * s.T.Z + v.Y * s.B.Z + v.Z * s.N.Z)
}

func (s *SampleBase) CoordinateSystem() (math.Vector3, math.Vector3, math.Vector3) {
	return s.T, s.B, s.N
}

func (s *SampleBase) SameHemisphere(v math.Vector3) bool {
	if s.N.Dot(v) > 0.0 {
		return true
	}

	return false
}
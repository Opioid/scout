package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
)

type Material interface {
	Evaluate(dg *geometry.Differential, l, v math.Vector3) (math.Vector3, float32)
	EvaluateAmbient(dg *geometry.Differential) (math.Vector3, float32)

	EvaluateSpecular(dg *geometry.Differential, l, v math.Vector3) math.Vector3

	Roughness() float32

	IsMirror() bool
}
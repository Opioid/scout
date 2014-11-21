package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	renderermaterial "github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
)

type Material interface {
	Evaluate(dg *geometry.Differential, v math.Vector3) renderermaterial.SubstituteBrdf

	IsMirror() bool
}
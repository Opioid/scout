package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	renderingmaterial "github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
)

type Material interface {
	Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D) renderingmaterial.Values

	IsMirror() bool
}

type Sample interface {
	Evaluate(l math.Vector3) math.Vector3
}
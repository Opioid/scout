package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	renderingmaterial "github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
)

type Material interface {
	Energy() math.Vector3

	Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerId uint32) renderingmaterial.Sample

	Free(sample renderingmaterial.Sample, workerId uint32)

	IsMirror() bool
	IsLight() bool
}


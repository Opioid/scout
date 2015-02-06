package substitute

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type ColorMap_NormalMap struct {
	base
	metallic float32
	roughness float32
	colorMap, normalMap *texture.Texture2D
}

func NewColorMap_NormalMap(roughness, metallic float32, colorMap, normalMap *texture.Texture2D, stack *BinnedStack) *ColorMap_NormalMap {
	m := new(ColorMap_NormalMap)
	m.stack = stack
	m.metallic = metallic
	m.roughness = math32.Max(roughness, minRoughness)
	m.colorMap = colorMap
	m.normalMap = normalMap
	return m
}

func (m *ColorMap_NormalMap) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerID uint32) material.Sample {
	color := sampler.Sample(m.colorMap, dg.UV).Vector3()

	nm := sampler.Sample(m.normalMap, dg.UV).Vector3()

	n := dg.TangentToWorld(nm).Normalized()

	s := m.stack.Pop(workerID)

	s.N = n
	s.T, s.B = math.CoordinateSystem(n)

	s.values.Set(color, 1.0, m.roughness, m.metallic, n, v)
	return s		
}

func (m *ColorMap_NormalMap) IsMirror() bool {
	return m.roughness <= minRoughness
}
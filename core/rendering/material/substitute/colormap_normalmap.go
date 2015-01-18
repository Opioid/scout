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

func NewColorMap_NormalMap(roughness, metallic float32, colorMap, normalMap *texture.Texture2D, pool *Pool) *ColorMap_NormalMap {
	m := new(ColorMap_NormalMap)
	m.pool = pool
	m.metallic = metallic
	m.roughness = math32.Max(roughness, minRoughness)
	m.colorMap = colorMap
	m.normalMap = normalMap
	return m
}

func (m *ColorMap_NormalMap) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerId uint32) material.Sample {
	nm := sampler.Sample(m.normalMap, dg.UV).Vector3()

	tangentToWorldSpace := math.MakeMatrix3x3FromAxes(dg.T, dg.B, dg.N)

	n := tangentToWorldSpace.TransformVector3(nm).Normalized()

//	s := new(Sample)
	s := m.pool.Get(workerId)
	s.values.Set(math.MakeVector3(0.75, 0.75, 0.75), 1, m.roughness, m.metallic, n, v)
	return s		
}

func (m *ColorMap_NormalMap) IsMirror() bool {
	return m.roughness <= minRoughness
}
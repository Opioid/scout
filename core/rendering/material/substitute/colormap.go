package substitute

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type ColorMap struct {
	base
	metallic float32
	roughness float32
	colorMap *texture.Texture2D
}

func NewColorMap(roughness, metallic float32, colorMap *texture.Texture2D, stack *BinnedStack) *ColorMap {
	m := new(ColorMap)
	m.stack = stack
	m.metallic = metallic
	m.roughness = math32.Max(roughness, minRoughness)
	m.colorMap = colorMap
	return m
}

func (m *ColorMap) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerID uint32) material.Sample {
	cs := sampler.Sample(m.colorMap, dg.UV)
	s := m.stack.Pop(workerID)
	s.T = dg.T
	s.B = dg.B
	s.N = dg.N
	s.Wo = v	
	s.Set(cs.Vector3(), cs.W, m.roughness, m.metallic)
	return s	
}

func (m *ColorMap) IsMirror() bool {
	return m.roughness <= minRoughness
}
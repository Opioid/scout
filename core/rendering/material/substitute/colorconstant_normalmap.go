package substitute

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type ColorConstant_NormalMap struct {
	base
	color math.Vector3
	metallic float32
	roughness float32
	normalMap *texture.Texture2D
}

func NewColorConstant_NormalMap(color math.Vector3, roughness, metallic float32, normalMap *texture.Texture2D, stack *BinnedStack) *ColorConstant_NormalMap {
	m := new(ColorConstant_NormalMap)
	m.stack = stack
	m.color = color
	m.metallic = metallic
	m.roughness = math32.Max(roughness, minRoughness)
	m.normalMap = normalMap
	return m
}

func (m *ColorConstant_NormalMap) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerID uint32) material.Sample {
	nm := sampler.Sample3(m.normalMap, dg.UV)

	n := dg.TangentToWorld(nm).Normalized()	

	s := m.stack.Pop(workerID)

	s.N = n
	s.T, s.B = math.CoordinateSystem(n)

	s.Wo = v

	s.Set(m.color, 1.0, m.roughness, m.metallic)
	return s	
}

func (m *ColorConstant_NormalMap) IsMirror() bool {
	return m.roughness <= minRoughness
}
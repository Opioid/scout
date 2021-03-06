package substitute

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_ "fmt"
)

type ColorConstant struct {
	base
	color math.Vector3
	metallic float32
	roughness float32
}

func NewColorConstant(color math.Vector3, roughness, metallic float32, stack *BinnedStack) *ColorConstant {
	m := new(ColorConstant)
	m.stack = stack
	m.color = color
	m.metallic = metallic	
	m.roughness = math32.Max(roughness, minRoughness)
	return m
}

func (m *ColorConstant) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerID uint32) material.Sample {
	s := m.stack.Pop(workerID)
	s.T = dg.T
	s.B = dg.B
	s.N = dg.N
	s.Wo = v
	s.Set(m.color, 1.0, m.roughness, m.metallic)
	return s
}

func (m *ColorConstant) IsMirror() bool {
	return m.roughness <= minRoughness
}
package substitute

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type ColorConstant struct {
	base
	color math.Vector3
	metallic float32
	roughness float32

}

func NewColorConstant(color math.Vector3, roughness, metallic float32, pool *Pool) *ColorConstant {
	m := new(ColorConstant)
	m.pool = pool
	m.color = color
	m.metallic = metallic	
	m.roughness = math32.Max(roughness, minRoughness)
	return m
}

func (m *ColorConstant) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerId uint32) material.Sample {
//	s := new(Sample)
	s := m.pool.Get(workerId)
	s.values.Set(m.color, 1, m.roughness, m.metallic, dg.N, v)
	return s
}

func (m *ColorConstant) IsMirror() bool {
	return m.roughness <= minRoughness
}
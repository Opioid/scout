package glass

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/math32"
)

type ColorConstant_NormalMap struct {
	base
	color math.Vector3
	normalMap *texture.Texture2D
}

func NewColorConstant_NormalMap(color math.Vector3, normalMap *texture.Texture2D, pool *Pool) *ColorConstant_NormalMap {
	m := new(ColorConstant_NormalMap)
	m.pool = pool
	m.color = color
	m.normalMap = normalMap
	return m
}

func (m *ColorConstant_NormalMap) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerId uint32) material.Sample {
	nm := sampler.Sample(m.normalMap, dg.UV).Vector3()

	tangentToWorldSpace := math.MakeMatrix3x3FromAxes(dg.T, dg.B, dg.N)

	n := tangentToWorldSpace.TransformVector3(nm).Normalized()

	s := m.pool.Get(workerId)
	s.values.Set(m.color, 1.0, 0.0, 0.0, n, v)
	return s
}

func (m *ColorConstant_NormalMap) IsMirror() bool {
	return true
}
package light

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
)

type ColorConstant struct {
	base
}

func NewColorConstant(l light.Light) *ColorConstant {
	m := new(ColorConstant)
	m.light = l
	return m
}

func (m *ColorConstant) Energy() math.Vector3 {
	return m.light.Color().Scale(m.light.Lumen())
}

func (m *ColorConstant) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerId uint32) material.Sample {
	return nil
}
package glass

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/math32"
)

type ColorConstant struct {
	base
	color math.Vector3
}

func NewColorConstant(color math.Vector3, stack *BinnedStack) *ColorConstant {
	m := new(ColorConstant)
	m.stack = stack
	m.color = color
	return m
}

func (m *ColorConstant) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D, workerID uint32) material.Sample {
	s := m.stack.Pop(workerID)
	s.values.Set(m.color, 1.0, 0.0, 0.0, dg.N, v)
	return s
}

func (m *ColorConstant) IsMirror() bool {
	return true
}
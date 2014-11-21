package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
)

type Substitute_ColorConstant struct {
	color math.Vector3
	metallic float32
	roughness, a2 float32

}

func NewSubstitute_ColorConstant(color math.Vector3, roughness, metallic float32) *Substitute_ColorConstant {
	m := new(Substitute_ColorConstant)
	m.color = color
	m.metallic = metallic	
	m.roughness = roughness
	a := roughness * roughness
	m.a2 = a * a
	return m
}

func (m *Substitute_ColorConstant) Evaluate(dg *geometry.Differential, v math.Vector3) SubstituteBrdf {
	return MakeSubstituteBrdf(m.color, 1, m.roughness, m.metallic, dg.N, v)
}

func (m *Substitute_ColorConstant) IsMirror() bool {
	return m.a2 == 0
}
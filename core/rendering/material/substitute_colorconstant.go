package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Substitute_ColorConstant struct {
	color math.Vector3
	metallic float32
	roughness float32

}

func NewSubstitute_ColorConstant(color math.Vector3, roughness, metallic float32) *Substitute_ColorConstant {
	m := new(Substitute_ColorConstant)
	m.color = color
	m.metallic = metallic	
	m.roughness = math32.Max(roughness, minRoughness)
	return m
}

func (m *Substitute_ColorConstant) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D) SubstituteBrdf {
	return MakeSubstituteBrdf(m.color, 1, m.roughness, m.metallic, dg.N, v)
}

func (m *Substitute_ColorConstant) IsMirror() bool {
	return m.roughness <= minRoughness
}
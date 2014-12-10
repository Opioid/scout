package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
)

type Substitute_ColorConstant_NormalMap struct {
	color math.Vector3
	metallic float32
	roughness, a2 float32
	normalMap *texture.Texture2D
}

func NewSubstitute_ColorConstant_NormalMap(color math.Vector3, roughness, metallic float32, normalMap *texture.Texture2D) *Substitute_ColorConstant_NormalMap {
	m := new(Substitute_ColorConstant_NormalMap)
	m.color = color
	m.metallic = metallic
	m.roughness = roughness
	a := roughness * roughness
	m.a2 = a * a
	m.normalMap = normalMap
	return m
}

func (m *Substitute_ColorConstant_NormalMap) Evaluate(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D) SubstituteBrdf {
	nm := sampler.Sample(m.normalMap, dg.UV).Vector3()

	tangentToWorldSpace := math.MakeMatrix3x3FromAxes(dg.T, dg.B, dg.N)

	n := tangentToWorldSpace.TransformVector3(nm).Normalized()

	return MakeSubstituteBrdf(m.color, 1, m.roughness, m.metallic, n, v)
}

func (m *Substitute_ColorConstant_NormalMap) IsMirror() bool {
	return m.a2 == 0
}
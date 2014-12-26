package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Substitute_ColorMap_NormalMap struct {
	metallic float32
	roughness float32
	colorMap, normalMap *texture.Texture2D
}

func NewSubstitute_ColorMap_NormalMap(roughness, metallic float32, colorMap, normalMap *texture.Texture2D) *Substitute_ColorMap_NormalMap {
	m := new(Substitute_ColorMap_NormalMap)
	m.metallic = metallic
	m.roughness = math32.Max(roughness, minRoughness)
	m.colorMap = colorMap
	m.normalMap = normalMap
	return m
}

func (m *Substitute_ColorMap_NormalMap) Sample(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D) SubstituteBrdf {
//	cs  := sampler.Sample(m.colorMap, dg.UV)

	nm := sampler.Sample(m.normalMap, dg.UV).Vector3()

	tangentToWorldSpace := math.MakeMatrix3x3FromAxes(dg.T, dg.B, dg.N)

	n := tangentToWorldSpace.TransformVector3(nm).Normalized()

	return MakeSubstituteBrdf(math.MakeVector3(0.75, 0.75, 0.75), 1, m.roughness, m.metallic, n, v)
//	return MakeSubstituteBrdf(cs.Vector3(), cs.W, m.roughness, m.metallic, n, v)
}

func (m *Substitute_ColorMap_NormalMap) IsMirror() bool {
	return m.roughness <= minRoughness
}
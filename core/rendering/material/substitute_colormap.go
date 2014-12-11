package material

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Substitute_ColorMap struct {
	metallic float32
	roughness float32
	colorMap *texture.Texture2D
}

func NewSubstitute_ColorMap(roughness, metallic float32, colorMap *texture.Texture2D) *Substitute_ColorMap {
	m := new(Substitute_ColorMap)
	m.metallic = metallic
	m.roughness = math32.Max(roughness, minRoughness)
	m.colorMap = colorMap
	return m
}

func (m *Substitute_ColorMap) Evaluate(dg *geometry.Differential, v math.Vector3, sampler texture.Sampler2D) SubstituteBrdf {
	cs  := sampler.Sample(m.colorMap, dg.UV)
	return MakeSubstituteBrdf(cs.Vector3(), cs.W, m.roughness, m.metallic, dg.N, v)
}

func (m *Substitute_ColorMap) IsMirror() bool {
	return m.roughness <= minRoughness
}
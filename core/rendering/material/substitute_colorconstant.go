package material

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
)

type Substitute_ColorConstant struct {
	color math.Vector3
	// roughness float32
	a2 float32
}

func NewSubstitute_ColorConstant(color math.Vector3, roughness float32) *Substitute_ColorConstant {
	m := new(Substitute_ColorConstant)
	m.color = color
	a := roughness * roughness
	m.a2 = a * a
	return m
}

func (m *Substitute_ColorConstant) Evaluate(dg *shape.DifferentialGeometry, l, v math.Vector3) math.Vector3 {
	n_dot_l := math.Max(dg.Nn.Dot(l), 0.00001)
	n_dot_v := math.Max(dg.Nn.Dot(v), 0.0)

	h := v.Add(l).Normalized()

	n_dot_h := dg.Nn.Dot(h)
	v_dot_h := v.Dot(h)

	f0 := math.Vector3{0.03, 0.03, 0.03}

	specular := specular_f(v_dot_h, f0).Scale(specular_d(n_dot_h, m.a2)).Scale(specular_g(n_dot_l, n_dot_v, m.a2))

	return m.color.Add(specular).Scale(n_dot_l)
}

func (m *Substitute_ColorConstant) IsMirror() bool {
	return m.a2 == 0.0
}
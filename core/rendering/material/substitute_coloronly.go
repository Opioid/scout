package material

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type Substitute_ColorOnly struct {
	color math.Vector3
	// roughness float32
	a2 float32
}

func NewSubstitute_ColorOnly(color math.Vector3, roughness float32) *Substitute_ColorOnly {
	m := new(Substitute_ColorOnly)
	m.color = color
	a := roughness * roughness
	m.a2 = a * a
	return m
}

func (m *Substitute_ColorOnly) Evaluate(n, l, v math.Vector3) math.Vector3 {
	n_dot_l := math.Max(n.Dot(l), 0.00001)
	n_dot_v := math.Max(n.Dot(v), 0.0)

	h := v.Add(l).Normalized()

	n_dot_h := n.Dot(h)
	v_dot_h := v.Dot(h)

	f0 := math.Vector3{0.03, 0.03, 0.03}

//	a := m.Roughness * m.Roughness
//	a2 := a * a

	specular := specular_f(v_dot_h, f0).Scale(specular_d(n_dot_h, m.a2)).Scale(specular_g(n_dot_l, n_dot_v, m.a2))

	return m.color.Add(specular).Scale(n_dot_l)

	
/*
	n_dot_l := math.Max(n.Dot(l), 0.0)

	return m.Color.Scale(n_dot_l)
	*/
}

func (m *Substitute_ColorOnly) IsMirror() bool {
	return m.a2 == 0.0
}

func specular_f(v_dot_h float32, f0 math.Vector3) math.Vector3 {
	return f0.Add(math.Vector3{1.0 - f0.X, 1.0 - f0.Y, 1.0 - f0.Z}.Scale(math.Exp2((-5.55473 * v_dot_h - 6.98316) * v_dot_h)))
}

func specular_d(n_dot_h, a2 float32) float32 {
	d := n_dot_h * n_dot_h * (a2 - 1.0) + 1.0
	return a2 / (gomath.Pi * d * d)
}

func specular_g(n_dot_l, n_dot_v, a2 float32) float32 {
	g_v := n_dot_v + math.Sqrt((n_dot_v - n_dot_v * a2) * n_dot_v + a2)
	g_l := n_dot_l + math.Sqrt((n_dot_l - n_dot_l * a2) * n_dot_l + a2)
	return math.InverseSqrt(g_v * g_l)
}

/*
// GGX/Trowbridge-Reitz
float specular_d(float n_dot_h, float a2)
{
	float d = n_dot_h * n_dot_h * (a2 - 1.f) + 1.f;
	return a2 / (pi * d * d);
}

vec3 specular_f(float v_dot_h, vec3 f0)
{
	return f0 + (1.f - f0) * exp2((-5.55473 * v_dot_h - 6.98316) * v_dot_h);
}

float specular_g(float n_dot_l, float n_dot_v, float a2)
{
	float G_V = n_dot_v + sqrt((n_dot_v - n_dot_v * a2) * n_dot_v + a2);
	float G_L = n_dot_l + sqrt((n_dot_l - n_dot_l * a2) * n_dot_l + a2);
	return inversesqrt(G_V * G_L);
}

*/
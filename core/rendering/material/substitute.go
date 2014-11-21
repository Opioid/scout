package material

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

func specular_f(v_dot_h float32, f0 math.Vector3) math.Vector3 {
	return f0.Add(math.MakeVector3(1.0 - f0.X, 1.0 - f0.Y, 1.0 - f0.Z).Scale(math.Exp2((-5.55473 * v_dot_h - 6.98316) * v_dot_h)))
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

type SubstituteBrdf struct {
	Color   math.Vector3
	DiffuseColor math.Vector3
 	Opacity float32

 	n, v math.Vector3
 	N_dot_v float32

 	F0 math.Vector3
 	Roughness float32
 	a2 float32
}

func MakeSubstituteBrdf(color math.Vector3, opacity, roughness, metallic float32, n, v math.Vector3) SubstituteBrdf {
	brdf := SubstituteBrdf{}

	brdf.Color = color
	brdf.DiffuseColor = color.Scale(1 - metallic).Scale(opacity)
	brdf.Opacity = opacity
	brdf.n = n
	brdf.v = v
	brdf.N_dot_v = math.Maxf(n.Dot(v), 0)

	brdf.F0 = math.MakeVector3(0.03, 0.03, 0.03).Lerp(color, metallic).Scale(opacity)
	
	brdf.Roughness = roughness
	a := roughness * roughness
	brdf.a2 = a * a

	return brdf
}

func (brdf *SubstituteBrdf) Evaluate(l math.Vector3) math.Vector3 {
//	return math.MakeVector3(1, 1, 1)

/*
	n_dot_l := math.Maxf(dg.N.Dot(l), 0.00001)
	n_dot_v := math.Maxf(dg.N.Dot(v), 0.0)

	h := v.Add(l).Normalized()

	n_dot_h := dg.N.Dot(h)
	v_dot_h := v.Dot(h)

	f0 := math.MakeVector3(0.03, 0.03, 0.03).Lerp(m.color, m.metallic)

	specular := specular_f(v_dot_h, f0).Scale(specular_d(n_dot_h, m.a2)).Scale(specular_g(n_dot_l, n_dot_v, m.a2))

	diffuse := m.color.Scale(1.0 - m.metallic)

	return diffuse.Add(specular).Scale(n_dot_l), 1.0
	*/

	n_dot_l := math.Maxf(brdf.n.Dot(l), 0.00001)

	h := brdf.v.Add(l).Normalized()

	n_dot_h := brdf.n.Dot(h)
	v_dot_h := brdf.v.Dot(h)

	specular := specular_f(v_dot_h, brdf.F0).Scale(specular_d(n_dot_h, brdf.a2)).Scale(specular_g(n_dot_l, brdf.N_dot_v, brdf.a2))

	return brdf.DiffuseColor.Add(specular).Scale(n_dot_l)
}
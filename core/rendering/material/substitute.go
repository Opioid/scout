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
package ggx

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
)


func SpecularF(WoDotH float32, f0 math.Vector3) math.Vector3 {
	return f0.Add(math.MakeVector3(1.0 - f0.X, 1.0 - f0.Y, 1.0 - f0.Z).Scale(math.Exp2((-5.55473 * WoDotH - 6.98316) * WoDotH)))
}

func SpecularD(NdotH, a2 float32) float32 {
	d := NdotH * NdotH * (a2 - 1.0) + 1.0
//	return a2 / math32.Max((gomath.Pi * d * d), gomath.SmallestNonzeroFloat32)
	return a2 / (math32.Pi * d * d)
}

func SpecularG(NdotWi, NdotWo, a2 float32) float32 {
	g_v := NdotWo + math32.Sqrt((NdotWo - NdotWo * a2) * NdotWo + a2)
	g_l := NdotWi + math32.Sqrt((NdotWi - NdotWi * a2) * NdotWi + a2)
	return math32.Rsqrt(g_v * g_l)
}


/*
vec3 importance_sample_GGX(vec2 xi, float roughness, vec3 n)
{
	float a = roughness * roughness;

	float phi = 2.f * pi * xi.x;

	float cos_theta = sqrt((1.f - xi.y) / (1.f + (a * a - 1.f) * xi.y ));

	float sin_theta = sqrt(1.f - cos_theta * cos_theta );

	vec3 h = vec3(sin_theta * cos(phi), sin_theta * sin(phi), cos_theta);

	vec3 up = abs(n.z) < 0.999f ? vec3(0.f, 0, 1.f) : vec3(1.f, 0.f, 0.f);

	vec3 tangent_x = normalize(cross(up, n));

	vec3 tangent_y = cross(n, tangent_x);

	// Tangent to world space
	return h.x * tangent_x + h.y * tangent_y + h.z * n;
}
*/

func ImportanceSample(xi math.Vector2, roughness float32, n math.Vector3) math.Vector3 {
	a := roughness * roughness

	phi := 2.0 * gomath.Pi * xi.X

	cos_theta := math32.Sqrt((1.0 - xi.Y) / (1.0 + (a * a - 1.0) * xi.Y))

	sin_theta := math32.Sqrt(1.0 - cos_theta * cos_theta)

	h := math.MakeVector3(sin_theta * math.Cos(phi), sin_theta * math.Sin(phi), cos_theta)

	var up math.Vector3

	if math32.Abs(n.Z) < 0.999 {
		up = math.MakeVector3(0.0, 0.0, 1.0)
	} else {
		up = math.MakeVector3(1.0, 0.0, 0.0)
	}

	tangent_x := up.Cross(n).Normalized()

	tangent_y := n.Cross(tangent_x)

	// Tangent to world space
	return tangent_x.Scale(h.X).Add(tangent_y.Scale(h.Y)).Add(n.Scale(h.Z))
}

/*
float g1(float n_dot_v, float k)
{
	return n_dot_v / (n_dot_v * (1.f - k) + k);
}

float g_smith(float roughness, float n_dot_l, float n_dot_v)
{
	float r1 = roughness + 1.f;
	float k = (r1 * r1) / 8.f;

	return g1(n_dot_l, k) * g1(n_dot_v, k);
}
*/

func g1(n_dot_v, k float32) float32 {
	return n_dot_v / (n_dot_v * (1.0 - k) + k)
}

func G_smith(roughness, n_dot_l, n_dot_v float32) float32 {
	r1 := roughness + 1.0
	k  := (r1 + r1) / 8.0

	return g1(n_dot_l, k) * g1(n_dot_v, k)
}
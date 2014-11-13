package ggx

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

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

	cos_theta := math.Sqrt((1.0 - xi.Y) / (1.0 + (a * a - 1.0) * xi.Y))

	sin_theta := math.Sqrt(1.0 - cos_theta * cos_theta)

	h := math.MakeVector3(sin_theta * math.Cos(phi), sin_theta * math.Sin(phi), cos_theta)

	var up math.Vector3

	if math.Absf(n.Z) < 0.999 {
		up = math.MakeVector3(0.0, 0.0, 1.0)
	} else {
		up = math.MakeVector3(1.0, 0.0, 0.0)
	}

	tangent_x := up.Cross(n).Normalized()

	tangent_y := n.Cross(tangent_x)

	// Tangent to world space
	return tangent_x.Scale(h.X).Add(tangent_y.Scale(h.Y)).Add(n.Scale(h.Z))
}
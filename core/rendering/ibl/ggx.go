package ibl

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material/ggx"
	_ "github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/scout/base/math/random"
	_ "math"
	"os"
	"image/png"
	_ "runtime"
	_ "sync"	
	_ "strconv"
	_ "fmt"
)

func IntegrateGgxBrdf(numSamples uint32, buffer *texture.Buffer) {
	dimensions := buffer.Dimensions()

	dx := float32(dimensions.X - 1)
	dy := float32(dimensions.Y - 1)

	for y := 0; y < dimensions.Y; y++ {
		v := float32(y) / dy

		for x := 0; x < dimensions.X; x++ {
			roughness := float32(x) / dx

			brdf := intefrateBrdf(roughness, v, numSamples)

			buffer.Set(x, y, math.MakeVector4(brdf.X, brdf.Y, 0, 1))
		}
	}

	image := buffer.RGBA()

	fo, err := os.Create("ggx_brdf.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}

/*
float2 integrate_brdf(float roughness, float n_dot_v, uint32_t num_samples)
{
	const float3 n(0.f, 0.f, 1.f);

	n_dot_v = std::max(n_dot_v, 0.00001f);

	float3 v;
	v.x = sqrt(1.f - n_dot_v * n_dot_v); // sin
	v.y = 0.f;
	v.z = n_dot_v; // cos

	float a = 0.f;
	float b = 0.f;

	for (uint32_t i = 0; i < num_samples; ++i)
	{
		float2 xi = math::hammersley( i, num_samples);
		float3 h = importance_sample_GGX(xi, roughness, n);
		float3 l = 2.f * dot(v, h) * h - v;
		float n_dot_l = math::saturate(l.z);
		float n_dot_h = math::saturate(h.z);
		float v_dot_h = math::saturate(dot(v, h));

		if (n_dot_l > 0.f)
		{
			float g = g_smith(roughness, n_dot_v, n_dot_l);
			float g_vis = g * v_dot_h / (n_dot_h * n_dot_v);
			float fc = pow(1.f - v_dot_h, 5.f);
			a += (1.f - fc) * g_vis;
			b += fc * g_vis;
		}
	}

	return float2(a, b) / float(num_samples);
}*/

func intefrateBrdf(roughness, n_dot_v float32, numSamples uint32) math.Vector2 {
	n := math.MakeVector3(0, 0, 1)

	n_dot_v = math.Maxf(n_dot_v, 0.00001)

	var v math.Vector3
	v.X = math.Sqrt(1 - n_dot_v * n_dot_v) // sin
	v.Y = 0
	v.Z = n_dot_v // cos

	a := float32(0)
	b := float32(0)

	for i := uint32(0); i < numSamples; i++ {
		xi := math.Hammersley(i, numSamples)
		h  := ggx.ImportanceSample(xi, roughness, n)
		l  := h.Scale(2.0 * v.Dot(h)).Sub(v)

		n_dot_l := math.Saturate(l.Z)
		n_dot_h := math.Saturate(h.Z)
		v_dot_h := math.Saturate(v.Dot(h))	

		if n_dot_l > 0 {
			g := ggx.G_smith(roughness, n_dot_v, n_dot_l)
			g_vis := g * v_dot_h / (n_dot_h * n_dot_v)
			fc := math.Pow(1 - v_dot_h, 5)

			a += (1 - fc) * g_vis
			b += fc * g_vis
		}	

	}

	return math.MakeVector2(a, b).Div(float32(numSamples))
}
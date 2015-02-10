package substitute

import (
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_ "fmt"
)

const (
	// magic roughness constant that doesn't cause INF in specular_d
	// instead there is a max() now
	// 0.01313900625 
//	minRoughness = 0.0// 0.01313900625
	// Ran into another issue in specular_g, which doesn't (??) require such a high minRoughness. 
	// Mabye keeping it above 0 makes sense anyway. Don't know about the specular_d issue now.
	minRoughness = 1.0 / 255.0
)

func specular_f(v_dot_h float32, f0 math.Vector3) math.Vector3 {
	return f0.Add(math.MakeVector3(1.0 - f0.X, 1.0 - f0.Y, 1.0 - f0.Z).Scale(math.Exp2((-5.55473 * v_dot_h - 6.98316) * v_dot_h)))
}

func specular_d(n_dot_h, a2 float32) float32 {
	d := n_dot_h * n_dot_h * (a2 - 1.0) + 1.0
//	return a2 / math32.Max((gomath.Pi * d * d), gomath.SmallestNonzeroFloat32)
	return a2 / (math32.Pi * d * d)
}

func specular_g(n_dot_l, n_dot_v, a2 float32) float32 {
	g_v := n_dot_v + math32.Sqrt((n_dot_v - n_dot_v * a2) * n_dot_v + a2)
	g_l := n_dot_l + math32.Sqrt((n_dot_l - n_dot_l * a2) * n_dot_l + a2)
	return math32.Rsqrt(g_v * g_l)
}

type base struct {
	stack *BinnedStack
}

func (b *base) Free(sample material.Sample, workerID uint32) {
	b.stack.Push(workerID)
}

func (b *base) Energy() math.Vector3 {
	return math.MakeVector3(0.0, 0.0, 0.0)
}

func (b *base) IsLight() bool {
	return false
}

type Sample struct {
	material.SampleBase
	values material.Values

	metallic float32

	lambert LambertBxdf
	ggx     GgxBxdf
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	nDotL := math32.Max(s.values.N.Dot(l), 0.00001)
/*
	h := s.values.V.Add(l).Normalized()

	nDotH := s.values.N.Dot(h)
	vDotH := s.values.V.Dot(h)

	specular := specular_f(vDotH, s.values.F0).Scale(specular_d(nDotH, s.values.A2)).Scale(specular_g(nDotL, s.values.N_dot_v, s.values.A2))

	return s.values.DiffuseColor.Scale(math32.InvPi).Add(specular).Scale(nDotL)
*/
	return s.values.DiffuseColor.Scale(math32.InvPi).Scale(nDotL)
}

func (s *Sample) Values() *material.Values {
	return &s.values
}

func (s *Sample) MonteCarloBxdf(subsample uint32, sampler sampler.Sampler) (material.Bxdf, float32) {
/*	if s.metallic == 1.0 {
		s.ggx.set(s.values.V, s.values.N, s.values.F0, s.values.A2)
		return &s.ggx, 1.0
	} else {

		p := sampler.GenerateSample1D(0, 0)

		if p < 0.5 {
			s.lambert.set(s.values.DiffuseColor)
			return &s.lambert, 0.5
		} else {
			s.ggx.set(s.values.V, s.values.N, s.values.F0, s.values.A2)
			return &s.ggx, 0.5
		}
	}
*/

	s.lambert.set(s.values.DiffuseColor)
	return &s.lambert, 1.0
}

type LambertBxdf struct {
	color math.Vector3
}

func (b *LambertBxdf) set(color math.Vector3) {
	b.color = color
} 

func (b *LambertBxdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, float32) {
	sample := sampler.GenerateSample2D(0, subsample) 

	hs := math.SampleHemisphereCosine(sample.X, sample.Y)
	pdf := float32(1.0)

	return hs, pdf
}

func (b *LambertBxdf) Evaluate(l math.Vector3) math.Vector3 {
	// Div by Pi is not neccessary because it is implicitly handled by the cosine distributed importance sample!
//	return b.color.Div(gomath.Pi)

	return b.color
}

type GgxBxdf struct {
	v, n math.Vector3
	f0 math.Vector3
	a2 float32
}

func (b *GgxBxdf) set(v, n math.Vector3, f0 math.Vector3, a2 float32) {
	b.v = v
	b.n = n
	b.f0 = f0
	b.a2 = a2
} 

func (b *GgxBxdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, float32) {
	xi := sampler.GenerateSample2D(0, subsample) 

	phi := 2.0 * math32.Pi * xi.X

	cos_theta := math32.Sqrt((1.0 - xi.Y) / (1.0 + (b.a2 - 1.0) * xi.Y))

	sin_theta := math32.Sqrt(1.0 - cos_theta * cos_theta)

	sin_phi, cos_phi := math.Sincos(phi)

	h := math.MakeVector3(sin_theta * cos_phi, sin_theta * sin_phi, cos_theta)	

	return h, 1.0
}

func (b *GgxBxdf) Evaluate(l math.Vector3) math.Vector3 {
	n_dot_l := math32.Max(b.n.Dot(l), 0.00001)
	n_dot_v := math32.Max(b.n.Dot(b.v), 0.0)

	h := b.v.Add(l).Normalized()

	v_dot_h := b.v.Dot(h)

	specular := specular_f(v_dot_h, b.f0).Scale(specular_g(n_dot_l, n_dot_v, b.a2))

	return specular.Scale(n_dot_l)
}
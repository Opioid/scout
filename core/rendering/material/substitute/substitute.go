package substitute

import (
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/material/ggx"
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

	lambert LambertBrdf
	ggx     GgxBrdf
}

func NewSample() *Sample {
	s := &Sample{}
	s.lambert.sample = s
	s.ggx.sample = s
	return s
}

func (s *Sample) Set(color math.Vector3, opacity, roughness, metallic float32) {
	s.metallic = metallic
	s.values.Set(color, opacity, roughness, metallic)
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	NdotWi := math32.Max(s.N.Dot(l), 0.00001)
	NdotWo := math32.Max(s.N.Dot(s.Wo), 0.00001)

	h := s.Wo.Add(l).Normalized()

	NdotH  := s.N.Dot(h)
	WoDotH := s.Wo.Dot(h)

	specular := ggx.F(WoDotH, s.values.F0).Scale(ggx.D(NdotH, s.values.A2) * ggx.G(NdotWi, NdotWo, s.values.A2))

	return s.values.DiffuseColor.Scale(math32.InvPi).Add(specular).Scale(NdotWi)

//	return s.values.DiffuseColor.Scale(math32.InvPi).Scale(NdotWi)

//	return specular.Scale(NdotWi)

//	f := ggx.F(WoDotH, s.values.F0)
//	d := ggx.D(NdotH, s.values.A2)
//	g := ggx.G_smith(s.values.Roughness, NdotWi, s.values.NdotWo)

//	specular := f.Scale(d * g).Div(4.0 * NdotWi * s.values.NdotWo)
//	specular := math.MakeVector3(d, d, d)

//	return specular.Scale(NdotWi)
}

func (s *Sample) Values() *material.Values {
	return &s.values
}

func (s *Sample) SampleEvaluate(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32) {
/*	if s.metallic == 1.0 {
		r, wi, _, pdf := s.ggx.ImportanceSample(subsample, sampler)
		return r, wi, pdf
	}

	p := sampler.GenerateSample1D(0, 0)

	if p < 0.5 {
		r0, wi, NdotWi, pdf0 := s.lambert.ImportanceSampleTest(subsample, sampler)

		r1, pdf1 := s.ggx.Evaluate(wi, NdotWi)

		return r0.Add(r1), wi, (pdf0 + pdf1) * 0.5
	} else {
		r0, wi, NdotWi, pdf0 := s.ggx.ImportanceSampleTest(subsample, sampler)

		r1, pdf1 := s.lambert.Evaluate(wi, NdotWi)

		return r0.Add(r1), wi, (pdf0 + pdf1) * 0.5
	}
*/

	if !s.SameHemisphere(s.Wo) {
		return math.MakeVector3(0.0, 0.0, 0.0), math.MakeVector3(0.0, 0.0, 0.0), 0.0
	}

	if s.metallic == 1.0 {
		r, wi, _, pdf := s.ggx.ImportanceSample(subsample, sampler)
		return r, wi, pdf
	} else {
		p := sampler.GenerateSample1D(0, 0)

		if p < 0.5 {
			r, wi, _, pdf := s.lambert.ImportanceSample(subsample, sampler)
			return r, wi, pdf * 0.5
		} else {
			r, wi, _, pdf := s.ggx.ImportanceSample(subsample, sampler)
			return r, wi, pdf * 0.5
		}
	}

//	r0, wi, _, pdf := s.lambert.ImportanceSampleTest(subsample, sampler)
//	return r0, wi, pdf

//	r0, wi, _, pdf := s.ggx.ImportanceSampleTest(subsample, sampler)
//	return r0, wi, pdf	
}

type LambertBrdf struct {
	sample *Sample
}

func (b *LambertBrdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32) {
	sample := sampler.GenerateSample2D(0, subsample) 

	is := math.SampleHemisphereCosine(sample.X, sample.Y)
	wi := b.sample.TangentToWorld(is).Normalized()

	NdotWi := math32.Max(b.sample.N.Dot(wi), 0.00001)

//	return b.sample.values.DiffuseColor.Scale(math32.InvPi * NdotWi), wi, NdotWi, math32.InvPi * NdotWi

	return b.sample.values.DiffuseColor, wi, NdotWi, 1.0
}

func (b *LambertBrdf) ImportanceSampleTest(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32) {
	sample := sampler.GenerateSample2D(0, subsample) 

	is := math.SampleHemisphereCosine(sample.X, sample.Y)
	wi := b.sample.TangentToWorld(is).Normalized()

	NdotWi := math32.Max(b.sample.N.Dot(wi), 0.00001)

	return b.sample.values.DiffuseColor.Scale(math32.InvPi * NdotWi), wi, NdotWi, math32.InvPi * NdotWi
}

func (b *LambertBrdf) Evaluate(wi math.Vector3, NdotWi float32) (math.Vector3, float32) {
	return b.sample.values.DiffuseColor.Scale(math32.InvPi * NdotWi), math32.InvPi * NdotWi
}

type GgxBrdf struct {
	sample *Sample
}

func (b *GgxBrdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32) {
	xi := sampler.GenerateSample2D(0, subsample) 

	NdotH := math32.Sqrt((1.0 - xi.Y) / ((b.sample.values.A2 - 1.0) * xi.Y + 1.0))
	sintheta := math32.Sqrt(1.0 - NdotH * NdotH)
	phi := 2.0 * math32.Pi * xi.X
	sinphi, cosphi := math.Sincos(phi)

	is := math.MakeVector3(sintheta * cosphi, sintheta * sinphi, NdotH)	
	h  := b.sample.TangentToWorld(is)

	WoDotH := b.sample.Wo.Dot(h)

	wi := h.Scale(2.0 * WoDotH).Sub(b.sample.Wo).Normalized()

	NdotWi := math32.Max(b.sample.N.Dot(wi), 0.00001)
	NdotWo := math32.Max(b.sample.N.Dot(b.sample.Wo), 0.0)

	f := ggx.F(WoDotH, b.sample.values.F0)
	g := ggx.G(NdotWi, NdotWo, b.sample.values.A2)

	specular := f.Scale(g)
	r := specular.Scale(NdotWi)
	return r, wi, NdotWi, NdotH / (4.0 * WoDotH)
}

func (b *GgxBrdf) ImportanceSampleTest(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32) {
	xi := sampler.GenerateSample2D(0, subsample) 

	NdotH := math32.Sqrt((1.0 - xi.Y) / ((b.sample.values.A2 - 1.0) * xi.Y + 1.0))
	sintheta := math32.Sqrt(1.0 - NdotH * NdotH)
	phi := 2.0 * math32.Pi * xi.X
	sinphi, cosphi := math.Sincos(phi)

	is := math.MakeVector3(sintheta * cosphi, sintheta * sinphi, NdotH)	
	h  := b.sample.TangentToWorld(is)

	WoDotH := b.sample.Wo.Dot(h)

	wi := h.Scale(2.0 * WoDotH).Sub(b.sample.Wo).Normalized()

	NdotWi := math32.Max(b.sample.N.Dot(wi), 0.00001)
	NdotWo := math32.Max(b.sample.N.Dot(b.sample.Wo), 0.0)

	f := ggx.F(WoDotH, b.sample.values.F0)
	d := ggx.D(NdotH, b.sample.values.A2)
	g := ggx.G(NdotWi, NdotWo, b.sample.values.A2)

	specular := f.Scale(d * g)
	r := specular.Scale(NdotWi)
	return r, wi, NdotWi, d * NdotH / (4.0 * WoDotH)
}

func (b *GgxBrdf) Evaluate(l math.Vector3, NdotWi float32) (math.Vector3, float32) {
//	NdotWi := math32.Max(b.sample.values.N.Dot(l), 0.00001)
	NdotWo := math32.Max(b.sample.N.Dot(b.sample.Wo), 0.00001)

	h := b.sample.Wo.Add(l).Normalized()

	NdotH := math32.Max(b.sample.N.Dot(h), 0.0)

	WoDotH := b.sample.Wo.Dot(h)

	f := ggx.F(WoDotH, b.sample.values.F0)
	d := ggx.D(NdotH, b.sample.values.A2)
	g := ggx.G(NdotWi, NdotWo, b.sample.values.A2)

	specular := f.Scale(d * g)

	return specular.Scale(NdotWi), d * NdotH / (4.0 * WoDotH)
}
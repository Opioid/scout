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

func (s *Sample) Set(color math.Vector3, opacity, roughness, metallic float32, n, wo math.Vector3) {
	s.metallic = metallic
	s.values.Set(color, opacity, roughness, metallic, n, wo)
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	NdotL := math32.Max(s.values.N.Dot(l), 0.00001)

	h := s.values.Wo.Add(l).Normalized()

	NdotH := s.values.N.Dot(h)
	WoDotH := s.values.Wo.Dot(h)

	specular := ggx.SpecularF(WoDotH, s.values.F0).Scale(ggx.SpecularD(NdotH, s.values.A2)).Scale(ggx.SpecularG(NdotL, s.values.NdotWo, s.values.A2))

	return s.values.DiffuseColor.Scale(math32.InvPi).Add(specular).Scale(NdotL)

//	return s.values.DiffuseColor.Scale(math32.InvPi).Scale(NdotL)

//	return specular.Scale(NdotL)
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
		r0, wi, NdotWi, pdf0 := s.lambert.ImportanceSample(subsample, sampler)

		r1, pdf1 := s.ggx.Evaluate(wi, NdotWi)

		return r0.Add(r1), wi, (pdf0 + pdf1)
	} else {
		r0, wi, NdotWi, pdf0 := s.ggx.ImportanceSample(subsample, sampler)

		r1, pdf1 := s.lambert.Evaluate(wi, NdotWi)

		return r0.Add(r1), wi, (pdf0 + pdf1)
	}
*/

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

//	r0, wi, _, pdf := s.lambert.ImportanceSample(subsample, sampler)
//	return r0, wi, pdf

//	r0, wi, _, pdf := s.ggx.ImportanceSample(subsample, sampler)
//	return r0, wi, pdf	
}

type LambertBrdf struct {
	sample *Sample
}

func (b *LambertBrdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32) {
	sample := sampler.GenerateSample2D(0, subsample) 

	is := math.SampleHemisphereCosine(sample.X, sample.Y)
	wi := b.sample.TangentToWorld(is).Normalized()

	NdotWi := math32.Max(b.sample.values.N.Dot(wi), 0.00001)

//	return b.sample.values.DiffuseColor.Scale(math32.InvPi * NdotWi), wi, NdotWi, math32.InvPi * NdotWi

	return b.sample.values.DiffuseColor, wi, NdotWi, 1.0
}


func (b *LambertBrdf) Evaluate(wi math.Vector3, NdotWi float32) (math.Vector3, float32) {
	return b.sample.values.DiffuseColor.Scale(math32.InvPi * NdotWi), math32.InvPi * NdotWi
}

type GgxBrdf struct {
	sample *Sample
}

func (b *GgxBrdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32) {
	xi := sampler.GenerateSample2D(0, subsample) 

	phi := 2.0 * math32.Pi * xi.X

	costheta := math32.Sqrt((1.0 - xi.Y) / ((b.sample.values.A2 - 1.0) * xi.Y + 1.0))
	sintheta := math32.Sqrt(1.0 - costheta * costheta)
	sinphi, cosphi := math.Sincos(phi)

	is := math.MakeVector3(sintheta * cosphi, sintheta * sinphi, costheta)	
	is  = b.sample.TangentToWorld(is)

//	is := math.SampleHemisphereUniform(xi.X, xi.Y)
//	is  = b.sample.TangentToWorld(is).Normalized()	

	// trying to avoid division by zero here, doesn't seem to fix the firefly problem though
	WoDotIs := math32.Max(b.sample.values.Wo.Dot(is), 0.00001)

	wi := is.Scale(2.0 * WoDotIs).Sub(b.sample.values.Wo).Normalized()

	NdotWi := math32.Max(b.sample.values.N.Dot(wi), 0.00001)
	NdotWo := math32.Max(b.sample.values.N.Dot(b.sample.values.Wo), 0.0)

	h := b.sample.values.Wo.Add(wi).Normalized()

	WoDotH := b.sample.values.Wo.Dot(h)

	specular := ggx.SpecularF(WoDotH, b.sample.values.F0).Scale(ggx.SpecularG(NdotWi, NdotWo, b.sample.values.A2))
	r := specular.Scale(NdotWi)
	return r, wi, NdotWi, costheta / (4.0 * WoDotH)


/*	costheta := b.sample.values.N.Dot(h)
	specular := ggx.SpecularF(WoDotH, b.sample.values.F0).Scale(ggx.SpecularD(costheta, b.sample.values.A2)).Scale(ggx.SpecularG(NdotWi, NdotWo, b.sample.values.A2))

	return specular.Scale(WoDotH / (costheta * NdotWo)), wi, NdotWi, 1.0 / (2.0 * math32.Pi)
	*/
}

func (b *GgxBrdf) Evaluate(l math.Vector3, NdotWi float32) (math.Vector3, float32) {
	NdotL := math32.Max(b.sample.values.N.Dot(l), 0.00001)
	NdotWo := math32.Max(b.sample.values.N.Dot(b.sample.values.Wo), 0.0)

	h := b.sample.values.Wo.Add(l).Normalized()

	WoDotH := b.sample.values.Wo.Dot(h)

	specular := ggx.SpecularF(WoDotH, b.sample.values.F0).Scale(ggx.SpecularG(NdotL, NdotWo,  b.sample.values.A2))

	costheta := math32.Abs(b.sample.values.N.Dot(h))

//	return specular.Scale(NdotL), costheta / (4.0 * WoDotH)
	
	return specular.Scale(NdotL), ggx.SpecularD(costheta, b.sample.values.A2) * costheta / (4.0 * WoDotH)
}
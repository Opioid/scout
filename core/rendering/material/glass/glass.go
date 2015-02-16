package glass

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/math32"
	_ "math"
	_ "fmt"
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

	btdf Btdf
}

func NewSample() *Sample {
	s := &Sample{}
	s.btdf.sample = s
//	s.ggx.sample = s
	return s
}

func (s *Sample) Set(color math.Vector3, opacity, roughness, metallic float32, n, wo math.Vector3) {
	s.values.Set(color, opacity, roughness, metallic, n, wo)
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	return math.MakeVector3(0.3, 0.3, 0.3)
}

func (s *Sample) Values() *material.Values {
	return &s.values
}

func (s *Sample) MonteCarloBxdf(subsample uint32, sampler sampler.Sampler) (material.Bxdf, float32) {
	return &s.btdf, 1.0
}

type Btdf struct {
	sample *Sample
}

func (b *Btdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, float32) {

	eta := float32(1.0)

	cosi := b.sample.values.Wo.Scale(-1.0).Dot(b.sample.values.N)
//	cost2 := 1.0 - eta * eta * (1.0 - cosi * cosi)
	t := b.sample.values.Wo.Scale(eta).Add(b.sample.values.N.Scale(eta * cosi /*- math32.Sqrt(math32.Abs(cost2))*/))

	return t.Normalized(), 1.0

}

func (b *Btdf) Evaluate(l math.Vector3) math.Vector3 {
	return b.sample.values.DiffuseColor
}

/*
float3 refract( float3 i, float3 n, float eta )
{
  float cosi = dot(-i, n);
  float cost2 = 1.0f - eta * eta * (1.0f - cosi*cosi);
  float3 t = eta*i + ((eta*cosi - sqrt(abs(cost2))) * n);
  return t * (float3)(cost2 > 0);
}*/
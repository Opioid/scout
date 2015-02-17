package glass

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/material"
	_ "github.com/Opioid/scout/core/rendering/material/ggx"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
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

func (b *Btdf) ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32) {

	etat := float32(1.3)

	eta := float32(1.0 / etat)

	n := b.sample.values.N.Scale(1.0)


	incident := b.sample.values.Wo.Scale(-1.0)

	cosi := -incident.Dot(n)

	if cosi < 0.0 {
		// hit from the inside
		cosi = -cosi
		n.ScaleAssign(-1.0)
		eta = float32(etat / 1.0)
		
	//	fmt.Println("From inside")
	} 

	cost2 := 1.0 - eta * eta * (1.0 - cosi * cosi)

	if cost2 < 0.0 {
		// total inner reflection
		return math.MakeVector3(0.0, 0.0, 0.0), math.MakeVector3(0.0, 0.0, 0.0), 0.0
	}


	t := incident.Scale(eta).Add(n.Scale(eta * cosi - math32.Sqrt(cost2)))

//	t := b.sample.values.Wo.Scale(eta).Add(n.Scale(eta * cosi - math32.Sqrt(math32.Abs(cost2))))

	wi := t.Normalized()

	return b.sample.values.DiffuseColor, wi, 1.0

/*
	f0 := math.MakeVector3(0.3, 0.3, 0.3)

	h := b.sample.values.Wo.Add(wi).Normalized()

	WoDotH := b.sample.values.Wo.Dot(h)

	fresnel := ggx.SpecularF(WoDotH, f0)

	return b.sample.values.DiffuseColor.Scale(1.0 - fresnel.X).Div(wi.Dot(n)), wi, 1.0
*/
	
	
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

/*
inline void Refract(
  VEC3 &out, const VEC3 &incidentVec, const VEC3 &normal, float eta)
{
  float N_dot_I = Dot(normal, incidentVec);
  float k = 1.f - eta * eta * (1.f - N_dot_I * N_dot_I);
  if (k < 0.f)
    out = VEC3(0.f, 0.f, 0.f);
  else
    out = eta * incidentVec - (eta * N_dot_I + sqrtf(k)) * N;
}*/
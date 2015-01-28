package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type normalSettings struct {
	linearSampler_repeat texture.Sampler2D
}

type normal struct {
	integrator
	normalSettings
}

func (ni *normal) StartNewPixel(numSamples uint32) {
}

func (ni *normal) Li(worker *rendering.Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	material := intersection.Material()

	v := ray.Direction.Scale(-1.0)
	brdf := material.Sample(&intersection.Geo.Differential, v, ni.linearSampler_repeat, ni.id)
	values := brdf.Values()

	result := values.N

	material.Free(brdf, ni.id)

	return result
}

func (ni *normal) MaxBounces() uint32 {
	return 0
}

func (ni *normal) PrimaryVisibility() uint8 {
	return prop.Primary
}

func (ni *normal) SecondaryVisibility() uint8 {
	return prop.Secondary
}

type normalFactory struct {
	normalSettings
}

func NewNormalFactory() *normalFactory {
	f := new(normalFactory)
	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat)) 
	return f
}

func (f *normalFactory) New(id uint32, rng *random.Generator) rendering.Integrator {
	ni := new(normal)

	ni.id = id
	ni.rng = rng
	ni.linearSampler_repeat = f.linearSampler_repeat


	return ni
}
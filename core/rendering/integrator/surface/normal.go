package surface

import (
	"github.com/Opioid/scout/core/rendering/integrator"
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
	integrator.Integrator
	normalSettings
}

func (ni *normal) StartNewPixel(numSamples uint32) {
}

func (ni *normal) Li(worker *rendering.Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	material := intersection.Material()

	v := ray.Direction.Scale(-1.0)
	brdf := material.Sample(&intersection.Geo.Differential, v, ni.linearSampler_repeat, ni.ID)

	_, _, n := brdf.CoordinateSystem()

	result := n.AddS(1.0).Scale(0.5)

	material.Free(brdf, ni.ID)

	return result
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

func (f *normalFactory) New(id uint32, rng *random.Generator) rendering.SurfaceIntegrator {
	ni := new(normal)

	ni.ID = id
	ni.Rng = rng
	ni.linearSampler_repeat = f.linearSampler_repeat

	return ni
}
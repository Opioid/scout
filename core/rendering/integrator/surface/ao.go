package surface

import (
	"github.com/Opioid/scout/core/rendering/integrator"
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type aoSettings struct {
	numSamples uint32
	numSamplesReciprocal float32
	samples []math.Vector2
	radius float32
}

type ao struct {
	integrator.Integrator
	sampler *pkgsampler.Random
	aoSettings
}

func (a *ao) StartNewPixel(numSamples uint32) {
	a.sampler.Restart(numSamples)
}

func (a *ao) Li(worker *rendering.Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	occlusionRay := math.OptimizedRay{}
	occlusionRay.Origin = intersection.Geo.P
	occlusionRay.MinT = intersection.Geo.Epsilon
	occlusionRay.MaxT = a.radius

//	basis := math.Matrix3x3{}
//	basis.SetBasis(intersection.Geo.N)

	result := float32(0.0)

	samples := a.sampler.GenerateSamples(subsample, a.samples) 

	for _, sample := range samples {
		s := math.SampleHemisphereCosine(sample.X, sample.Y)

	//	v := basis.TransformVector3(s)
		v := intersection.Geo.TangentToWorld(s)

		occlusionRay.SetDirection(v)

		if worker.Visibility(&occlusionRay) {
			result += a.numSamplesReciprocal
		}
	}

	return math.MakeVector3(result, result, result)
}

func (a *ao) PrimaryVisibility() uint8 {
	return prop.Primary
}

func (a *ao) SecondaryVisibility() uint8 {
	return prop.Secondary
}

type aoFactory struct {
	aoSettings
}

func NewAoFactory(numSamples uint32, radius float32) *aoFactory {
	f := new(aoFactory)

	f.numSamples = numSamples
	f.numSamplesReciprocal = 1.0 / float32(numSamples)
	f.radius = radius

	return f
}

func (f *aoFactory) New(id uint32, rng *random.Generator) rendering.SurfaceIntegrator {
	a := new(ao)

	a.ID = id
	a.Rng = rng
	a.sampler = pkgsampler.NewRandom(f.numSamples, rng)
	a.numSamples = f.numSamples	
	a.numSamplesReciprocal = f.numSamplesReciprocal
	a.samples = make([]math.Vector2, f.numSamples)
	a.radius = f.radius

	return a
}
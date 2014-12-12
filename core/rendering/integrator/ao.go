package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type aoSettings struct {
	numSamples uint32
	numSamplesReciprocal float32
	radius float32
}

type ao struct {
	integrator
	sampler pkgsampler.ScrambledHammersley
	aoSettings
}

func (a *ao) FirstSample(numSamples uint32) {
	a.sampler.Restart(numSamples)
}

func (a *ao) Li(scene *pkgscene.Scene, task *rendering.RenderTask, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	occlusionRay := math.OptimizedRay{}
	occlusionRay.Origin = intersection.Dg.P
	occlusionRay.MinT = intersection.Epsilon
	occlusionRay.MaxT = a.radius

	basis := math.Matrix3x3{}
	basis.SetBasis(intersection.Dg.N)

	result := float32(0.0)

	samples := a.sampler.GenerateSamples(subsample) 

	for _, sample := range samples {
		s := math.HemisphereSample_cos(sample.X, sample.Y)

		v := basis.TransformVector3(s)
		occlusionRay.SetDirection(v)

		if !scene.IntersectP(&occlusionRay) {
			result += a.numSamplesReciprocal
		}
	}

	return math.MakeVector3(result, result, result)
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

func (f *aoFactory) New(rng *random.Generator) rendering.Integrator {
	a := new(ao)

	a.rng = rng
	a.sampler = pkgsampler.MakeScrambledHammersley(rng)
	a.sampler.Resize(f.numSamples)
	a.numSamples = f.numSamples	
	a.numSamplesReciprocal = 1.0 / float32(a.numSamples)
	a.radius = f.radius

	return a
}
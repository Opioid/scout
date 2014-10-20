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
	sampler pkgsampler.Stratified
	aoSettings
}

func (a *ao) Li(scene *pkgscene.Scene, task *rendering.RenderTask, subsample, numSamples uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	occlusionRay := math.OptimizedRay{}
	occlusionRay.Origin = intersection.Dg.P
	occlusionRay.MinT = intersection.Epsilon
	occlusionRay.MaxT = a.radius

	basis := math.Matrix3x3{}
	basis.SetBasis(intersection.Dg.N)

	result := float32(0.0)

//	offset          := a.numSamples * subsample
//	numTotalSamples := a.numSamples * numSamples

	a.sampler.Restart()

	sample := pkgsampler.Sample{}

//	for i := uint32(0); i < a.numSamples; i++ {
	for a.sampler.GenerateNewSample(&sample) {
	//	s := math.HemisphereSample_cos(a.rng.RandomFloat32(), a.rng.RandomFloat32())
		s := math.HemisphereSample_cos(sample.Coordinates.X, sample.Coordinates.Y)

	//	h := math.Hammersley(i + offset, numTotalSamples)
	//	s := hemisphereSample_cos(h.X, h.Y)

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
	f := aoFactory{}

	f.numSamples = numSamples
	f.numSamplesReciprocal = 1.0 / float32(numSamples)
	f.radius = radius

	return &f
}

func (f *aoFactory) New(rng *random.Generator) rendering.Integrator {
	a := ao{}

	a.rng = rng
	a.sampler = pkgsampler.MakeStratified(rng)
	a.sampler.Resize(math.MakeVector2i(4, 4))
	a.numSamples = 16//f.numSamples
	a.numSamplesReciprocal = 1.0 / float32(a.numSamples)
	a.radius = f.radius

	return &a
}
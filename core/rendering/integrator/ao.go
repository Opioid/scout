package integrator

import (
	"github.com/Opioid/scout/core/rendering"
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
	aoSettings
}

func (a *ao) Li(scene *pkgscene.Scene, task *rendering.RenderTask, sample, numSamples uint32, ray *math.OptimizedRay, intersection *prop.Intersection, rng *random.Generator) math.Vector3 {
	occlusionRay := math.OptimizedRay{}
	occlusionRay.Origin = intersection.Dg.P
	occlusionRay.MinT = intersection.Epsilon
	occlusionRay.MaxT = a.radius

	basis := math.Matrix3x3{}
	basis.SetBasis(intersection.Dg.N)

	result := float32(0.0)

//	offset          := a.numSamples * sample
//	numTotalSamples := a.numSamples * numSamples

	for i := uint32(0); i < a.numSamples; i++ {
		s := math.HemisphereSample_cos(rng.RandomFloat32(), rng.RandomFloat32())

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

func (f *aoFactory) New() rendering.Integrator {
	a := ao{}

	a.numSamples = f.numSamples
	a.numSamplesReciprocal = 1.0 / float32(f.numSamples)
	a.radius = f.radius

	return &a
}
package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	gomath "math"
	_ "fmt"
)

type ao struct {
	numSamples uint32
	numSamplesReciprocal float32
	radius float32
}

func NewAo(numSamples uint32, radius float32) *ao {
	a := ao{}

	a.numSamples = numSamples
	a.numSamplesReciprocal = 1.0 / float32(numSamples)
	a.radius = radius

	return &a
}

func (a *ao) Li(scene *pkgscene.Scene, renderer *rendering.Renderer, sample, numSamples uint32, ray *math.OptimizedRay, intersection *prop.Intersection, rng *random.Generator) math.Vector3 {
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
		s := hemisphereSample_cos(rng.RandomFloat32(), rng.RandomFloat32())

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

func hemisphereSample_cos(r, s float32) math.Vector3 {
	r1 := s * 2.0 * gomath.Pi
	r2 := math.Sqrt(1.0 - r)
	sp := math.Sqrt(1.0 - r2 * r2)

	return math.MakeVector3(math.Cos(r1) * sp, math.Sin(r1) * sp, r2)
}
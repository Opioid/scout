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
	numSamples int
	radius float32
}

func NewAo(numSamples int, radius float32) *ao {
	return &ao{numSamples, radius}
}

func (a *ao) Li(scene *pkgscene.Scene, renderer *rendering.Renderer, ray *math.OptimizedRay, intersection *prop.Intersection, rng *random.Generator) math.Vector3 {
	occlusionRay := math.OptimizedRay{}
	occlusionRay.Origin = intersection.Dg.P
	occlusionRay.MinT = intersection.Epsilon
	occlusionRay.MaxT = a.radius

	result := float32(0.0)

	for i := 0; i < a.numSamples; i++ {
		v := hemisphereSample(intersection.Dg.N, rng.RandomFloat32(), rng.RandomFloat32())
		occlusionRay.SetDirection(v)

		if !scene.IntersectP(&occlusionRay) {
			result += 1.0
		}
	}

	result /= float32(a.numSamples)

	return math.Vector3{result, result, result}
}

func hemisphereSample(v math.Vector3, r, s float32) math.Vector3 {
	basis := math.Matrix3x3{}
	basis.SetBasis(v)

	r1 := s * 2.0 * gomath.Pi
	r2 := math.Sqrt(1.0 - r)
	sp := math.Sqrt(1.0 - r2 * r2)

	dir := math.Vector3{math.Cos(r1) * sp, math.Sin(r1) * sp, r2}

	return basis.TransformVector3(dir)
}
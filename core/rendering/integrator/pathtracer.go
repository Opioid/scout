package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type pathtracerSettings struct {
	numSamples uint32
	maxBounces uint32

	numSamplesReciprocal float32

	samples []math.Vector2

	secondaryRay math.OptimizedRay

	linearSampler_repeat texture.Sampler2D

	shadowRay math.OptimizedRay
}

type pathtracer struct {
	integrator
	sampler *pkgsampler.Random
	pathtracerSettings

	lightSamples []light.Sample
}

func (pt *pathtracer) StartNewPixel(numSamples uint32) {
	pt.sampler.Restart(numSamples * pt.maxBounces)
}

func (pt *pathtracer) Li(worker *rendering.Worker, subsample uint32, scene *pkgscene.Scene, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	material := intersection.Material()
/*
	if material.IsLight() {
		return material.Energy()
	}
*/
	result := math.MakeVector3(0.0, 0.0, 0.0)

	nextDepth := ray.Depth + 1

	if nextDepth > pt.maxBounces {
		return result
	}

	v := ray.Direction.Scale(-1.0)
	brdf := material.Sample(&intersection.Geo.Differential, v, pt.linearSampler_repeat, pt.id)

	l := scene.RandomLight(pt.rng.RandomFloat32())

	if l != nil {
		pt.shadowRay.Origin = intersection.Geo.P
		pt.shadowRay.MinT = intersection.Geo.Epsilon
		pt.shadowRay.MaxT = 1000.0
		pt.shadowRay.Time = ray.Time

		pt.lightSamples = l.Samples(intersection.Geo.P, ray.Time, subsample, 1, pt.sampler, pt.lightSamples)

		numSamplesReciprocal := 1.0 / float32(len(pt.lightSamples))

		for _, s := range pt.lightSamples {
			pt.shadowRay.SetDirection(s.L)

			if !scene.IntersectP(&pt.shadowRay) {
				r := brdf.Evaluate(s.L)

				result.AddAssign(s.Energy.Mul(r).Scale(numSamplesReciprocal))
			}
		}
	} else {
		
		samples := pt.sampler.GenerateSamples(ray.Depth + subsample * pt.maxBounces, pt.samples) 

		values := brdf.Values()

		basis := math.Matrix3x3{}
		basis.SetBasis(values.N)

		for _, sample := range samples {
			s := math.HemisphereSample_cos(sample.X, sample.Y)

			v := basis.TransformVector3(s)
		//	v := intersection.Geo.TangentToWorld(s)

			pt.secondaryRay.Set(intersection.Geo.P, v, intersection.Geo.Epsilon, 1000.0, ray.Time, nextDepth)

			environment := worker.Li(subsample, scene, &pt.secondaryRay)

			result.AddAssign((environment).Scale(pt.numSamplesReciprocal))
		}

		result.MulAssign(values.DiffuseColor)
	
	}

	material.Free(brdf, pt.id)

	return result

}

func (pt *pathtracer) MaxBounces() uint32 {
	return pt.maxBounces
}

type pathtracerFactory struct {
	pathtracerSettings
}

func NewPathtracerFactory(numSamples, maxBounces uint32) *pathtracerFactory {
	f := new(pathtracerFactory)

	f.numSamples = numSamples
	f.maxBounces = maxBounces

	f.numSamplesReciprocal = 1.0 / float32(numSamples)

	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat))

	return f
}

func (f *pathtracerFactory) New(id uint32, rng *random.Generator) rendering.Integrator {
	pt := new(pathtracer)

	pt.id = id
	pt.rng = rng
	pt.numSamples = f.numSamples
	pt.numSamplesReciprocal = f.numSamplesReciprocal
	pt.maxBounces = f.maxBounces
	pt.sampler = pkgsampler.NewRandom(f.numSamples, rng)
	pt.samples = make([]math.Vector2, f.numSamples)

	pt.lightSamples = make([]light.Sample, 0, 1)

	pt.linearSampler_repeat = f.linearSampler_repeat

	return pt
}
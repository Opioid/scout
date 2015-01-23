package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/prop"
	_ "github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type pathtracerSettings struct {
	maxBounces uint32

	secondaryRay math.OptimizedRay

	linearSampler_repeat texture.Sampler2D
}

type pathtracer struct {
	integrator
	sampler *pkgsampler.Random
	pathtracerSettings
}

func (pt *pathtracer) StartNewPixel(numSamples uint32) {
	pt.sampler.Restart(numSamples * pt.maxBounces)
}

func (pt *pathtracer) Li(worker *rendering.Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
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

	l, probability := worker.Scene.RandomLight(pt.rng.RandomFloat32())

	pt.secondaryRay.Origin = intersection.Geo.P
	pt.secondaryRay.MinT = intersection.Geo.Epsilon
	pt.secondaryRay.MaxT = 1000.0
	pt.secondaryRay.Time = ray.Time
	pt.secondaryRay.Depth = nextDepth

	if l != nil {
		ls := l.Sample(&worker.Transformation, intersection.Geo.P, ray.Time, subsample, pt.sampler)

		pt.secondaryRay.SetDirection(ls.L)

	//	if !scene.IntersectP(&pt.secondaryRay) {
		if !worker.Shadow(&pt.secondaryRay) {	
			r := brdf.Evaluate(ls.L)

			result.AddAssign(ls.Energy.Mul(r))
		}

	} else {
		values := brdf.Values()

		basis := math.Matrix3x3{}
		basis.SetBasis(values.N)

		sample := pt.sampler.GenerateSample(0, ray.Depth + subsample * pt.maxBounces) 
		hs := math.HemisphereSample_cos(sample.X, sample.Y)

		v := basis.TransformVector3(hs)
	//	v := intersection.Geo.TangentToWorld(s)

		pt.secondaryRay.SetDirection(v)

		environment := worker.Li(subsample, &pt.secondaryRay)

		result.AddAssign(values.DiffuseColor.Mul(environment))
	}

	result.DivAssign(probability)

	material.Free(brdf, pt.id)

	return result
}

func (pt *pathtracer) MaxBounces() uint32 {
	return pt.maxBounces
}

type pathtracerFactory struct {
	pathtracerSettings
}

func NewPathtracerFactory(maxBounces uint32) *pathtracerFactory {
	f := new(pathtracerFactory)

	f.maxBounces = maxBounces

	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat))

	return f
}

func (f *pathtracerFactory) New(id uint32, rng *random.Generator) rendering.Integrator {
	pt := new(pathtracer)

	pt.id = id
	pt.rng = rng
	pt.maxBounces = f.maxBounces
	pt.sampler = pkgsampler.NewRandom(1024, rng)

	pt.linearSampler_repeat = f.linearSampler_repeat

	return pt
}
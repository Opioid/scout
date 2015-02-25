package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/prop"
	_ "github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "github.com/Opioid/math32"
	_ "fmt"
)

type pathtracerSettings struct {
	minBounces, maxBounces uint32

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

	throughput := math.MakeVector3(1.0, 1.0, 1.0)
	result := math.MakeVector3(0.0, 0.0, 0.0)
	hit := true

	for i := uint32(0); i <= pt.maxBounces; i++ {
		material := intersection.Material()

		if material.IsLight() {
			l := ray.Direction.Scale(-1.0)
			nDotL := intersection.Geo.N.Dot(l)

			if nDotL > 0.0 {
				result.AddAssign(throughput.Mul(material.Energy()))
			} 

			break
		}

		nextDepth := ray.Depth + 1

		if nextDepth > pt.maxBounces {
			break
		}

		eye := ray.Direction.Scale(-1.0)
		materialSample := material.Sample(&intersection.Geo.Differential, eye, pt.linearSampler_repeat, pt.id)

		r, wi, pdf := materialSample.SampleEvaluate(ray.Depth + subsample * pt.maxBounces, pt.sampler)

		material.Free(materialSample, pt.id)

		if pdf == 0.0 {
			break
		}

		throughput.MulAssign(r.Div(pdf))

		ray.Origin = intersection.Geo.P
		ray.SetDirection(wi)
		ray.MinT = intersection.Geo.Epsilon
		ray.MaxT = 1000.0
		ray.Depth = nextDepth

	//	oldp := ray.Origin

		if hit, intersection = worker.Intersect(ray); !hit {
			r := worker.Scene.Surrounding.Sample(ray)
			result.AddAssign(throughput.Mul(r))			
			break
		} else {
		//	fmt.Printf("(%v, %v) -> %v\n", oldp, wi, intersection.Geo.P)
		}

	}

	return result
}

func (pt *pathtracer) PrimaryVisibility() uint8 {
	return prop.Primary
}

func (pt *pathtracer) SecondaryVisibility() uint8 {
	return prop.Secondary | prop.IsLight
}

type pathtracerFactory struct {
	pathtracerSettings
}

func NewPathtracerFactory(minBounces, maxBounces uint32) *pathtracerFactory {
	f := new(pathtracerFactory)

	f.minBounces = minBounces
	f.maxBounces = maxBounces

	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat))

	return f
}

func (f *pathtracerFactory) New(id uint32, rng *random.Generator) rendering.Integrator {
	pt := new(pathtracer)

	pt.id = id
	pt.rng = rng
	pt.minBounces = f.minBounces
	pt.maxBounces = f.maxBounces
	pt.sampler = pkgsampler.NewRandom(1024, rng)

	pt.linearSampler_repeat = f.linearSampler_repeat

	return pt
}
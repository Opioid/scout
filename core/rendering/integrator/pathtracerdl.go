package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/prop"
	_ "github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type pathtracerDlSettings struct {
	maxBounces uint32

	secondaryRay math.OptimizedRay

	linearSampler_repeat texture.Sampler2D
}

type pathtracerDl struct {
	integrator
	sampler *pkgsampler.Random
	pathtracerDlSettings
}

func (pt *pathtracerDl) StartNewPixel(numSamples uint32) {
	pt.sampler.Restart(numSamples * pt.maxBounces)
}

func (pt *pathtracerDl) Li(worker *rendering.Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	throughput := math.MakeVector3(1.0, 1.0, 1.0)
	result := math.MakeVector3(0.0, 0.0, 0.0)
	hit := true

	for i := uint32(0); i < pt.maxBounces; i++ {
		nextDepth := ray.Depth + 1
	/*	if nextDepth > pt.maxBounces {
			break
		}
*/
		// No handling of geometry from the "inside" for now
		if ray.Direction.Dot(intersection.Geo.N) > 0.0 {
			break
		}

		pt.secondaryRay.Origin = intersection.Geo.P
		pt.secondaryRay.MinT = intersection.Geo.Epsilon
		pt.secondaryRay.Time = ray.Time
		pt.secondaryRay.Depth = nextDepth

		eye := ray.Direction.Scale(-1.0)
		material := intersection.Material()
		materialSample := material.Sample(&intersection.Geo.Differential, eye, pt.linearSampler_repeat, pt.id)

		l, lightPdf := worker.Scene.MonteCarloLight(pt.rng.RandomFloat32())

		if l != nil {
			ls := l.Sample(&worker.ScratchBuffer.Transformation, intersection.Geo.P, ray.Time, subsample, pt.sampler)

			if ls.Pdf > 0.0 {
				pt.secondaryRay.SetDirection(ls.L)
				pt.secondaryRay.MaxT = ls.T

				if worker.Visibility(&pt.secondaryRay) {
					r := materialSample.Evaluate(ls.L)
					result.AddAssign(throughput.Mul(ls.Energy.Mul(r).Div(lightPdf * ls.Pdf)))
				}
			}
		} 

		bxdf, samplePdf := materialSample.MonteCarloBxdf(ray.Depth + subsample * pt.maxBounces, pt.sampler)

		wi, bxdfPdf := bxdf.ImportanceSample(ray.Depth + subsample * pt.maxBounces, pt.sampler)
		r := bxdf.Evaluate(wi)

		material.Free(materialSample, pt.id)

		throughput.MulAssign(r.Div(samplePdf * bxdfPdf))

		ray.Origin = intersection.Geo.P
		ray.SetDirection(wi)
		ray.MinT = intersection.Geo.Epsilon
		ray.MaxT = 1000.0
		ray.Depth = nextDepth

		if hit, intersection = worker.Intersect(ray); !hit {
			r := worker.Scene.Surrounding.Sample(ray)
			result.AddAssign(throughput.Mul(r))
			break
		} 
	}

	return result
}

func (pt *pathtracerDl) MaxBounces() uint32 {
	return pt.maxBounces
}

func (pt *pathtracerDl) PrimaryVisibility() uint8 {
	return prop.Primary
}

func (pt *pathtracerDl) SecondaryVisibility() uint8 {
	return prop.Secondary
}

type pathtracerDlFactory struct {
	pathtracerDlSettings
}

func NewPathtracerDlFactory(maxBounces uint32) *pathtracerDlFactory {
	f := new(pathtracerDlFactory)

	f.maxBounces = maxBounces

	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat))

	return f
}

func (f *pathtracerDlFactory) New(id uint32, rng *random.Generator) rendering.Integrator {
	pt := new(pathtracerDl)

	pt.id = id
	pt.rng = rng
	pt.maxBounces = f.maxBounces
	pt.sampler = pkgsampler.NewRandom(1024, rng)

	pt.linearSampler_repeat = f.linearSampler_repeat

	return pt
}
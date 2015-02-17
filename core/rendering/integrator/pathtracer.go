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

		// No handling of geometry from the "inside" for now
	/*	if eye.Dot(intersection.Geo.N) < 0.0 {
			fmt.Println("dobdob")
			break
		}*/
	
		materialSample := material.Sample(&intersection.Geo.Differential, eye, pt.linearSampler_repeat, pt.id)

		bxdf, samplePdf := materialSample.MonteCarloBxdf(ray.Depth + subsample * pt.maxBounces, pt.sampler)

		r, wi, bxdfPdf := bxdf.ImportanceSample(ray.Depth + subsample * pt.maxBounces, pt.sampler)

		material.Free(materialSample, pt.id)

		combinedPdf := samplePdf * bxdfPdf
		
		if combinedPdf == 0.0 {
			break
		}

		throughput.MulAssign(r.Div(combinedPdf))

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

/*
	material := intersection.Material()

	if material.IsLight() {
		l := ray.Direction.Scale(-1.0)
		nDotL := intersection.Geo.N.Dot(l)

		if nDotL > 0.0 {
			return material.Energy()
		} else {
			return math.MakeVector3(0.0, 0.0, 0.0)
		}
	}

	nextDepth := ray.Depth + 1

	if nextDepth > pt.maxBounces {
		return math.MakeVector3(0.0, 0.0, 0.0)
	}

	eye := ray.Direction.Scale(-1.0)

	// No handling of geometry from the "inside" for now
	if eye.Dot(intersection.Geo.N) < 0.0 {
		return math.MakeVector3(0.0, 0.0, 0.0)
	}

	materialSample := material.Sample(&intersection.Geo.Differential, eye, pt.linearSampler_repeat, pt.id)

	bxdf, samplePdf := materialSample.MonteCarloBxdf(ray.Depth + subsample * pt.maxBounces, pt.sampler)

	hs, bxdfPdf := bxdf.ImportanceSample(ray.Depth + subsample * pt.maxBounces, pt.sampler)
	v := materialSample.TangentToWorld(hs)

	r := bxdf.Evaluate(v)

	material.Free(materialSample, pt.id)

	pt.secondaryRay.Origin = intersection.Geo.P

	pt.secondaryRay.SetDirection(v)
	pt.secondaryRay.MinT = intersection.Geo.Epsilon
	pt.secondaryRay.MaxT = 1000.0
	pt.secondaryRay.Time = ray.Time
	pt.secondaryRay.Depth = nextDepth

	environment := worker.Li(subsample, &pt.secondaryRay)

	return r.Mul(environment).Div(samplePdf * bxdfPdf)
	*/
}

func (pt *pathtracer) MaxBounces() uint32 {
	return pt.maxBounces
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
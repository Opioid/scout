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

	if material.IsLight() {
		return material.Energy()
	}

	nextDepth := ray.Depth + 1

	if nextDepth > pt.maxBounces {
		return math.MakeVector3(0.0, 0.0, 0.0)
	}

	// No handling of geometry from the "inside" for now
	if ray.Direction.Dot(intersection.Geo.N) > 0.0 {
		return math.MakeVector3(0.0, 0.0, 0.0)
	}

	eye := ray.Direction.Scale(-1.0)
	materialSample := material.Sample(&intersection.Geo.Differential, eye, pt.linearSampler_repeat, pt.id)

	bxdf, bp := materialSample.MonteCarloBxdf(ray.Depth + subsample * pt.maxBounces, pt.sampler)

	hs := bxdf.ImportanceSample(ray.Depth + subsample * pt.maxBounces, pt.sampler)
	v := materialSample.TangentToWorld(hs)

	r := bxdf.Evaluate(v)

	material.Free(materialSample, pt.id)

	pt.secondaryRay.Origin = intersection.Geo.P

//	fmt.Println(intersection.Geo.P)
//	panic("ja")

	pt.secondaryRay.SetDirection(v)
	pt.secondaryRay.MinT = intersection.Geo.Epsilon
	pt.secondaryRay.MaxT = 1000.0
	pt.secondaryRay.Time = ray.Time
	pt.secondaryRay.Depth = nextDepth

	environment := worker.Li(subsample, &pt.secondaryRay)

	return r.Mul(environment).Div(bp)
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
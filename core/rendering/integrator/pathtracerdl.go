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
	material := intersection.Material()
/*
	if material.IsLight() {
		return material.Energy()
	}
*/
	
	nextDepth := ray.Depth + 1

	if nextDepth > pt.maxBounces {
		return math.MakeVector3(0.0, 0.0, 0.0)
	}

	// No handling of geometry from the "inside" for now
	if ray.Direction.Dot(intersection.Geo.N) > 0.0 {
		return math.MakeVector3(0.0, 0.0, 0.0)
	}

	result := math.MakeVector3(0.0, 0.0, 0.0)

	pt.secondaryRay.Origin = intersection.Geo.P
	pt.secondaryRay.MinT = intersection.Geo.Epsilon
	pt.secondaryRay.Time = ray.Time
	pt.secondaryRay.Depth = nextDepth

	v := ray.Direction.Scale(-1.0)
	materialSample := material.Sample(&intersection.Geo.Differential, v, pt.linearSampler_repeat, pt.id)

	l, lp := worker.Scene.MonteCarloLight(pt.rng.RandomFloat32())

	if l != nil {
		ls := l.Sample(&worker.ScratchBuffer.Transformation, intersection.Geo.P, ray.Time, subsample, pt.sampler)

		pt.secondaryRay.SetDirection(ls.L)
		pt.secondaryRay.MaxT = ls.T

		if !worker.Shadow(&pt.secondaryRay) {
			r := materialSample.Evaluate(ls.L)

			result.AddAssign(ls.Energy.Mul(r).Div(lp))
		}

	} 

	/*else*/ {
		pt.secondaryRay.MaxT = 1000.0

	//	values := materialSample.Values()

	//	basis := math.Matrix3x3{}
	//	basis.SetBasis(values.N)

	//	sample := pt.sampler.GenerateSample(0, ray.Depth + subsample * pt.maxBounces) 
	//	hs := math.HemisphereSample_cos(sample.X, sample.Y)

		bxdf, bp := materialSample.MonteCarloBxdf(ray.Depth + subsample * pt.maxBounces, pt.sampler)

		hs := bxdf.ImportanceSample(ray.Depth + subsample * pt.maxBounces, pt.sampler)

	//	v := basis.TransformVector3(hs)
		v := materialSample.TangentToWorld(hs)

		r := bxdf.Evaluate(v)

		material.Free(materialSample, pt.id)

		pt.secondaryRay.SetDirection(v)

		environment := worker.Li(subsample, &pt.secondaryRay)

		result.AddAssign(r.Mul(environment).Div(bp))
	}

//	result.DivAssign(lp)



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
package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/texture/buffer"
	"github.com/Opioid/scout/core/rendering/ibl"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type whittedSettings struct {
	maxBounces uint32

	maxLightSamples uint32

	shadowRay math.OptimizedRay
	secondaryRay math.OptimizedRay

	linearSampler_repeat texture.Sampler2D
	linearSampler_clamp texture.Sampler2D

	brdf *texture.Texture2D	
}

type whitted struct {
	integrator
	whittedSettings

	sampler *pkgsampler.ScrambledHammersley

	lightSamples []light.Sample
}

func (w *whitted) StartNewPixel(numSamples uint32) {
	w.sampler.Restart(numSamples)
}

func (w *whitted) Li(worker *rendering.Worker, subsample uint32, scene *pkgscene.Scene, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	result := math.MakeVector3(0.0, 0.0, 0.0)

	w.shadowRay.Origin = intersection.Geo.P
	w.shadowRay.MinT = intersection.Geo.Epsilon
	w.shadowRay.MaxT = 1000.0
	w.shadowRay.Time = ray.Time

	v := ray.Direction.Scale(-1.0)

	material := intersection.Material()
	brdf := material.Sample(&intersection.Geo.Differential, v, w.linearSampler_repeat, w.id)
	values := brdf.Values()

	for _, l := range scene.Lights {
		w.lightSamples = l.Samples(intersection.Geo.P, ray.Time, subsample, w.maxLightSamples, w.sampler, w.lightSamples)

		numSamplesReciprocal := 1.0 / float32(len(w.lightSamples))

		for _, s := range w.lightSamples {
			w.shadowRay.SetDirection(s.L)

			if !scene.IntersectP(&w.shadowRay) {
				r := brdf.Evaluate(s.L)

				result.AddAssign(s.Energy.Mul(r).Scale(numSamplesReciprocal))
			}
		}
	}

	ambientColor := scene.Surrounding.SampleDiffuse(values.N)
	result.AddAssign(ambientColor.Mul(values.DiffuseColor))

	reflection := values.N.Reflect(ray.Direction).Normalized()

	var environment math.Vector3

	if material.IsMirror() && ray.Depth < w.maxBounces {
	//	secondaryRay := math.MakeOptimizedRay(intersection.Dg.P, reflection, intersection.Epsilon, 1000.0, ray.Time, ray.Depth + 1)

		// If this is the second (or more) bounce, this will overwrite ray, because they are actually refering to the same memory!
		w.secondaryRay.Set(intersection.Geo.P, reflection, intersection.Geo.Epsilon, 1000.0, ray.Time, ray.Depth + 1)

		environment = worker.Li(subsample, scene, &w.secondaryRay)
	} else {
		environment = scene.Surrounding.SampleSpecular(reflection, values.Roughness)
	}

	pi_brdf := w.linearSampler_clamp.Sample(w.brdf, math.MakeVector2(values.Roughness, values.N_dot_v))

	result.AddAssign(environment.Mul(values.F0.Scale(pi_brdf.X).AddS(pi_brdf.Y)))

	material.Free(brdf, w.id)

	return result
}

func (w *whitted) MaxBounces() uint32 {
	return w.maxBounces
}

type whittedFactory struct {
	whittedSettings
}

func NewWhittedFactory(maxBounces, maxLightSamples uint32) *whittedFactory {
	f := new(whittedFactory)

	f.maxBounces = maxBounces
	f.maxLightSamples = maxLightSamples

	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat)) 
	f.linearSampler_clamp = texture.NewSampler2D_linear(new(texture.AddressMode_clamp)) 

	f.brdf = texture.NewTexture2D(buffer.Float4, math.MakeVector2i(128, 128), 1)
	ibl.IntegrateGgxBrdf(1024, f.brdf.Image.Buffers[0])

	return f
}

func (f *whittedFactory) New(id uint32, rng *random.Generator) rendering.Integrator {
	w := new(whitted)

	w.id = id
	w.rng = rng
	w.maxBounces = f.maxBounces
	w.maxLightSamples = f.maxLightSamples
//	w.sampler = pkgsampler.MakeStratified(rng)
//	w.sampler.Resize(math.MakeVector2i(4, 4))
	w.sampler = pkgsampler.NewScrambledHammersley(w.maxLightSamples, rng)
	w.lightSamples = make([]light.Sample, 0, w.maxLightSamples)

	w.linearSampler_repeat = f.linearSampler_repeat
	w.linearSampler_clamp = f.linearSampler_clamp
	w.brdf = f.brdf

	return w
}
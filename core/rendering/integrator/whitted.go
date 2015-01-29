package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/texture/buffer"
	"github.com/Opioid/scout/core/rendering/ibl"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type whittedSettings struct {
	maxBounces uint32

	shadowRay math.OptimizedRay
	secondaryRay math.OptimizedRay

	linearSampler_repeat texture.Sampler2D
	linearSampler_clamp texture.Sampler2D

	brdf *texture.Texture2D	
}

type whitted struct {
	integrator
	whittedSettings

	sampler pkgsampler.Sampler
}

func (w *whitted) StartNewPixel(numSamples uint32) {
	w.sampler.Restart(numSamples)
}

func (w *whitted) Li(worker *rendering.Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	result := math.MakeVector3(0.0, 0.0, 0.0)

	w.shadowRay.Origin = intersection.Geo.P
	w.shadowRay.MinT = intersection.Geo.Epsilon
	w.shadowRay.Time = ray.Time

	material := intersection.Material()

	if material.IsLight() {
		return material.Energy()
	}

	v := ray.Direction.Scale(-1.0)
	brdf := material.Sample(&intersection.Geo.Differential, v, w.linearSampler_repeat, w.id)
	values := brdf.Values()

	for _, l := range worker.Scene.Lights {
		ls := l.Sample(&worker.ScratchBuffer.Transformation, intersection.Geo.P, ray.Time, subsample, w.sampler)

		w.shadowRay.SetDirection(ls.L)
		w.shadowRay.MaxT = ls.T

		if !worker.Shadow(&w.shadowRay) {
			r := brdf.Evaluate(ls.L)

			result.AddAssign(ls.Energy.Mul(r))
		}
	}

	ambientColor := worker.Scene.Surrounding.SampleDiffuse(values.N)
	result.AddAssign(ambientColor.Mul(values.DiffuseColor))

	reflection := values.N.Reflect(ray.Direction).Normalized()

	var environment math.Vector3

	if material.IsMirror() && ray.Depth < w.maxBounces {
	//	secondaryRay := math.MakeOptimizedRay(intersection.Dg.P, reflection, intersection.Epsilon, 1000.0, ray.Time, ray.Depth + 1)

		// If this is the second (or more) bounce, this will overwrite ray, because they are actually refering to the same memory!
		w.secondaryRay.Set(intersection.Geo.P, reflection, intersection.Geo.Epsilon, 1000.0, ray.Time, ray.Depth + 1)

		environment = worker.Li(subsample, &w.secondaryRay)
	} else {
		environment = worker.Scene.Surrounding.SampleSpecular(reflection, values.Roughness)
	}

	pi_brdf := w.linearSampler_clamp.Sample(w.brdf, math.MakeVector2(values.Roughness, values.N_dot_v))

	result.AddAssign(environment.Mul(values.F0.Scale(pi_brdf.X).AddS(pi_brdf.Y)))

	material.Free(brdf, w.id)

	return result
}

func (w *whitted) MaxBounces() uint32 {
	return w.maxBounces
}

func (w *whitted) PrimaryVisibility() uint8 {
	return prop.Primary
}

func (w *whitted) SecondaryVisibility() uint8 {
	return prop.Secondary
}

type whittedFactory struct {
	whittedSettings
}

func NewWhittedFactory(maxBounces uint32) *whittedFactory {
	f := new(whittedFactory)

	f.maxBounces = maxBounces

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
	w.sampler = pkgsampler.NewRandom(1024, rng)
	w.linearSampler_repeat = f.linearSampler_repeat
	w.linearSampler_clamp = f.linearSampler_clamp
	w.brdf = f.brdf

	return w
}
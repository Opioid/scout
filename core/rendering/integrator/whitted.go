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
	bounceDepth uint32

	maxLightSamples uint32

	linearSampler_repeat texture.Sampler2D
	linearSampler_clamp  texture.Sampler2D

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

func (w *whitted) Li(scene *pkgscene.Scene, task *rendering.Task, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	result := math.MakeVector3(0.0, 0.0, 0.0)

	material := intersection.Prop.Material

	shadowRay := math.OptimizedRay{}
	shadowRay.Origin = intersection.Dg.P
	shadowRay.MinT = intersection.Epsilon
	shadowRay.MaxT = 1000.0
	shadowRay.Time = ray.Time

	v := ray.Direction.Scale(-1.0)

	brdf := material.Sample(&intersection.Dg, v, w.linearSampler_repeat)

	for _, l := range scene.Lights {
		w.lightSamples = w.lightSamples[:0]

		l.Samples(intersection.Dg.P, subsample, 0.0, w.sampler, &w.lightSamples)

		numSamplesReciprocal := 1.0 / float32(len(w.lightSamples))

		for _, s := range w.lightSamples {
			shadowRay.SetDirection(s.L)

			if !scene.IntersectP(&shadowRay) {
				r := brdf.Evaluate(s.L)

				result.AddAssign(s.Energy.Mul(r).Scale(numSamplesReciprocal))		
			}
		}
	}

	ambientColor := scene.Surrounding.SampleDiffuse(brdf.N)
	result.AddAssign(ambientColor.Mul(brdf.DiffuseColor))

	reflection := brdf.N.Reflect(ray.Direction).Normalized()

	var environment math.Vector3

	if material.IsMirror() && ray.Depth < w.bounceDepth {
		secondaryRay := math.MakeOptimizedRay(intersection.Dg.P, reflection, intersection.Epsilon, 1000.0, ray.Time, ray.Depth + 1)

		environment = task.Li(scene, subsample, &secondaryRay)
	} else {
		environment = scene.Surrounding.SampleSpecular(reflection, brdf.Roughness)
	}

	pi_brdf := w.linearSampler_clamp.Sample(w.brdf, math.MakeVector2(brdf.Roughness, brdf.N_dot_v))

	result.AddAssign(environment.Scale(pi_brdf.X + pi_brdf.Y).Mul(brdf.F0))

	return result
}

type whittedFactory struct {
	whittedSettings
}

func NewWhittedFactory(bounceDepth, maxLightSamples uint32) *whittedFactory {
	f := new(whittedFactory)

	f.bounceDepth = bounceDepth
	f.maxLightSamples = maxLightSamples

	f.linearSampler_repeat = texture.NewSampler2D_linear(new(texture.AddressMode_repeat)) 
	f.linearSampler_clamp = texture.NewSampler2D_linear(new(texture.AddressMode_clamp)) 

	f.brdf = texture.NewTexture2D(buffer.Float4, math.MakeVector2i(32, 32), 1)
	ibl.IntegrateGgxBrdf(1024, f.brdf.Image.Buffers[0])

	return f
}

func (f *whittedFactory) New(rng *random.Generator) rendering.Integrator {
	w := new(whitted)

	w.rng = rng
	w.bounceDepth = f.bounceDepth
	w.maxLightSamples = f.maxLightSamples
//	w.sampler = pkgsampler.MakeStratified(rng)
//	w.sampler.Resize(math.MakeVector2i(4, 4))
	w.sampler = pkgsampler.NewScrambledHammersley(w.maxLightSamples, rng)
	w.lightSamples = make([]light.Sample, 0, w.maxLightSamples)

	w.linearSampler_repeat = f.linearSampler_repeat
	w.linearSampler_clamp  = f.linearSampler_clamp
	w.brdf = f.brdf

	return w
}
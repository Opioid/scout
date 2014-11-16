package integrator

import (
	"github.com/Opioid/scout/core/rendering"
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/texture"
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
}

type whitted struct {
	integrator
	whittedSettings

	sampler pkgsampler.ScrambledHammersley

	lightSamples []light.Sample

	brdf texture.Sampler2D
}

func (w *whitted) FirstSample(numSamples uint32) {
	w.sampler.Restart(numSamples)
}

func (w *whitted) Li(scene *pkgscene.Scene, task *rendering.RenderTask, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3 {
	result := math.MakeVector3(0.0, 0.0, 0.0)

	material := intersection.Prop.Material

	shadowRay := math.OptimizedRay{}
	shadowRay.Origin = intersection.Dg.P
	shadowRay.MinT = intersection.Epsilon
	shadowRay.MaxT = 1000.0

	v := ray.Direction.Scale(-1.0)

/*
	for _, l := range scene.Lights {
		lightVector := l.Vector(intersection.Dg.P)

		shadowRay.SetDirection(lightVector)

		if !scene.IntersectP(&shadowRay) {
			color, opacity := material.Evaluate(&intersection.Dg, lightVector, v)

			result.AddAssign(l.Light(intersection.Dg.P, color.Scale(opacity)))

		//	if opacity < 1.0 {
		//		secondaryRay := *ray
		//		secondaryRay.MinT = ray.MaxT + intersection.Epsilon
		//		secondaryRay.MaxT = 1000.0
				
		//		secondaryColor := r.li(scene, &secondaryRay, depth)
	
		//		result = result.Add(secondaryColor.Scale(1.0 - opacity))
		//	}
		}
	}
	*/

	for _, l := range scene.Lights {
		w.lightSamples = w.lightSamples[:0]

		l.Samples(intersection.Dg.P, subsample, &w.sampler, &w.lightSamples)

		numSamplesReciprocal := 1.0 / float32(len(w.lightSamples))

		for _, s := range w.lightSamples {
			shadowRay.SetDirection(s.L)

			if !scene.IntersectP(&shadowRay) {
				color, opacity := material.Evaluate(&intersection.Dg, s.L, v)
				result.AddAssign(s.Energy.Mul(color.Scale(opacity)).Scale(numSamplesReciprocal))
			}

		}
	}

	ambientColor := scene.Surrounding.SampleDiffuse(intersection.Dg.N)
	color, opacity := material.EvaluateAmbient(&intersection.Dg)
	result.AddAssign(ambientColor.Mul(color.Scale(opacity)))


/*
	numSamplesReciprocal := 1.0 / float32(w.maxLightSamples)

	basis := math.Matrix3x3{}

	basis.SetBasis(intersection.Dg.N)

	rn := w.rng.RandomUint32()

	sresult := math.MakeVector3(0.0, 0.0, 0.0)

	for i := uint32(0); i < w.maxLightSamples; i++ {
		sample := math.ScrambledHammersley(i, w.maxLightSamples, rn)

		s := math.HemisphereSample_cos(sample.X, sample.Y)

		l := basis.TransformVector3(s).Normalized()

		ambientColor = scene.Surrounding.SampleSpecular(l, 0)
			
		color = material.EvaluateSpecular(&intersection.Dg, l, v)

		sresult.AddAssign(ambientColor.Mul(color.Scale(opacity)).Scale(numSamplesReciprocal))
	}


	result.AddAssign(sresult)
*/

	reflection := intersection.Dg.N.Reflect(ray.Direction)

	if material.IsMirror() && ray.Depth < w.bounceDepth {
		secondaryRay := math.MakeOptimizedRay(intersection.Dg.P, reflection, intersection.Epsilon, 1000.0, ray.Depth + 1)

		result.AddAssign(task.Li(scene, subsample, &secondaryRay))
	} else {
		specularLight := scene.Surrounding.SampleSpecular(reflection, material.Roughness())
	//	color, opacity := material.EvaluateAmbient(&intersection.Dg)
		result.AddAssign(specularLight)
	}

	return result
}

type whittedFactory struct {
	whittedSettings

}

func NewWhittedFactory(bounceDepth, maxLightSamples uint32) *whittedFactory {
	f := whittedFactory{}

	f.bounceDepth = bounceDepth
	f.maxLightSamples = maxLightSamples

	return &f
}

func (f *whittedFactory) New(rng *random.Generator) rendering.Integrator {
	w := whitted{}

	w.rng = rng
	w.bounceDepth = f.bounceDepth
	w.maxLightSamples = f.maxLightSamples
//	w.sampler = pkgsampler.MakeStratified(rng)
//	w.sampler.Resize(math.MakeVector2i(4, 4))
	w.sampler = pkgsampler.MakeScrambledHammersley(rng)
	w.sampler.Resize(w.maxLightSamples)
	w.lightSamples = make([]light.Sample, 0, w.maxLightSamples)

	brdf := texture.NewTexture2D(math.MakeVector2i(32, 32), 1)

	ibl.IntegrateGgxBrdf(16, &brdf.Image.Buffers[0])

	w.brdf = texture.NewSampler2D_linear(brdf, new(texture.AddressMode_clamp)) 

	return &w
}
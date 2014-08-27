package rendering

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	"sync"
	_ "fmt"
)

type Renderer struct {
	Integrator Integrator

	samplerDimensions math.Vector2i
	currentPixel math.Vector2i
}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context) {
	dimensions := context.Camera.Film().Dimensions()

	r.currentPixel = math.Vector2i{0, 0}

	r.samplerDimensions = math.Vector2i{32, 32}

	wg := sync.WaitGroup{}

	for {
		sampler := r.newSubSampler(context.Sampler, dimensions)

		if sampler == nil {
			break
		}

		wg.Add(1)

		go func () {
			r.render(scene, context.Camera, sampler)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (r *Renderer) render(scene *pkgscene.Scene, camera camera.Camera, sampler pkgsampler.Sampler) {
	rng := random.Generator{}

	start := sampler.Start()
	rng.Seed(uint32(start.X) + 0, uint32(start.Y) + 1, uint32(start.X) + 2, uint32(start.Y) + 3)

	film := camera.Film()

	var ray math.OptimizedRay
	var sample pkgsampler.Sample

	for sampler.GenerateNewSample(&sample) {
		camera.GenerateRay(&sample, &ray)

		color := r.Li(scene, &ray, &rng) 

		film.AddSample(&sample, color)
	}
}

func (r *Renderer) Li(scene *pkgscene.Scene, ray *math.OptimizedRay, rng *random.Generator) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		return r.Integrator.Li(scene, r, ray, &intersection, rng) 
	} else {
		return scene.Surrounding.Sample(ray)
	}
}

func (r *Renderer) newSubSampler(s pkgsampler.Sampler, dimensions math.Vector2i) pkgsampler.Sampler {
	if r.currentPixel.X >= dimensions.X {
		r.currentPixel.X = 0
		r.currentPixel.Y += r.samplerDimensions.Y
	}

	if r.currentPixel.Y >= dimensions.Y {
		return nil
	}

	end := r.currentPixel.Add(r.samplerDimensions).Min(dimensions)

	sampler := s.SubSampler(r.currentPixel, end)

	r.currentPixel.X += r.samplerDimensions.X

	return sampler
}
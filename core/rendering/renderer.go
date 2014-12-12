package rendering

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/core/progress"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	"sync"
	_ "fmt"
)

type Renderer struct {
	IntegratorFactory IntegratorFactory

	samplerDimensions math.Vector2i
	currentPixel math.Vector2i
}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context, progressor progress.Sink) {
	dimensions := context.Camera.Film().Dimensions()

	r.currentPixel = math.MakeVector2i(0, 0)

	r.samplerDimensions = math.MakeVector2i(32, 32)

	numSamplers := int(float32(dimensions.X) / float32(r.samplerDimensions.X) + 0.5) * int(float32(dimensions.Y) / float32(r.samplerDimensions.Y) + 0.5)
	progressor.Start(numSamplers)

	wg := sync.WaitGroup{}

	for {

		if !r.advanceCurrentPixel(dimensions) {
			break
		}

		end := r.currentPixel.Add(r.samplerDimensions).Min(dimensions)

		sampler := context.Sampler.SubSampler(r.currentPixel, end)

		wg.Add(1)

		go func () {
			r.render(scene, context.Camera, sampler)
			progressor.Tick()
			wg.Done()
		}()
	}

	wg.Wait()
	progressor.End()
}

func (r *Renderer) render(scene *pkgscene.Scene, camera camera.Camera, sampler pkgsampler.Sampler) {
	task := RenderTask{}
	task.renderer = r

	start := sampler.Start()
	rng := random.Generator{}
	rng.Seed(uint32(start.X) + 0, uint32(start.Y) + 1, uint32(start.X) + 2, uint32(start.Y) + 3)	

	task.integrator = r.IntegratorFactory.New(&rng)

	film := camera.Film()

	var ray math.OptimizedRay
	var sample pkgsampler.Sample

	numSamples := sampler.NumSamplesPerPixel()

	for sampler.GenerateNewSample(&sample) {
		camera.GenerateRay(&sample, &ray)

		if 0 == sample.Id {
			task.FirstSample(numSamples)
		}

		color := task.Li(scene, sample.Id, &ray) 

		film.AddSample(&sample, color)
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

func (r *Renderer) advanceCurrentPixel(dimensions math.Vector2i) bool {
	if r.currentPixel.X >= dimensions.X {
		r.currentPixel.X = 0
		r.currentPixel.Y += r.samplerDimensions.Y
	}

	if r.currentPixel.Y >= dimensions.Y {
		return false
	}

	r.currentPixel.X += r.samplerDimensions.X

	return true
}
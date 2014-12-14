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
		rng := random.Generator{}
		rng.Seed(uint32(r.currentPixel.X) + 0, uint32(r.currentPixel.Y) + 1, uint32(r.currentPixel.X) + 2, uint32(r.currentPixel.Y) + 3)	

		end := r.currentPixel.Add(r.samplerDimensions).Min(dimensions)
		sampler := context.Sampler.Clone(&rng)

		wg.Add(1)

		go func (tileStart, tileEnd math.Vector2i) {
			r.render(scene, context.Camera, tileStart, tileEnd, sampler, &rng)
			progressor.Tick()
			wg.Done()
		}(r.currentPixel, end)

		if !r.advanceCurrentPixel(dimensions) {
			break
		}
	}

	wg.Wait()
	progressor.End()
}

func (r *Renderer) render(scene *pkgscene.Scene, camera camera.Camera, start, end math.Vector2i,
						  sampler pkgsampler.Sampler, rng *random.Generator) {
	task := makeTask(r, r.IntegratorFactory.New(rng))

	task.render(scene, camera, start, end, sampler)
/*
	film := camera.Film()

	var sample pkgsampler.Sample
	var ray math.OptimizedRay

	numSamples := sampler.NumSamplesPerPixel()

	for sampler.GenerateNewSample(&sample) {
		camera.GenerateRay(&sample, &ray)

		if 0 == sample.Id {
			task.FirstSample(numSamples)
		}

		color := task.Li(scene, sample.Id, &ray) 

		film.AddSample(&sample, color)
	}*/
}

func (r *Renderer) advanceCurrentPixel(dimensions math.Vector2i) bool {
	r.currentPixel.X += r.samplerDimensions.X

	if r.currentPixel.X >= dimensions.X {
		r.currentPixel.X = 0
		r.currentPixel.Y += r.samplerDimensions.Y
	}

	if r.currentPixel.Y >= dimensions.Y {
		return false
	}

	return true
}
package rendering

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Task struct {
	renderer *Renderer
	integrator Integrator
} 

func makeTask(renderer *Renderer, integrator Integrator) Task {
	t := Task{}
	t.renderer = renderer
	t.integrator = integrator
	return t
}

func (t *Task) render(scene *pkgscene.Scene, camera camera.Camera, start, end math.Vector2i, sampler pkgsampler.Sampler) {
	film := camera.Film()

	numSamples := sampler.NumSamplesPerPixel()

	var offset math.Vector2
	var sample pkgsampler.Sample
	var ray math.OptimizedRay

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			sampler.Restart()
			t.integrator.FirstSample(numSamples)
			sampleId := uint32(0)

			for sampler.GenerateNewSample(&offset) {
				sample.RelativeOffset = offset.Scale(0.5)
				sample.Coordinates = math.MakeVector2(float32(x) + 0.5 + sample.RelativeOffset.X, 
													  float32(y) + 0.5 + sample.RelativeOffset.Y)
				
				camera.GenerateRay(sample.Coordinates, &ray)

				color := t.Li(scene, sampleId, &ray) 

				film.AddSample(&sample, color)

				sampleId++
			}
		}
	}
}

func (t *Task) Li(scene *pkgscene.Scene, subsample uint32, ray *math.OptimizedRay) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		c := t.integrator.Li(scene, t, subsample, ray, &intersection) 
		return c
	} else {
		c := scene.Surrounding.Sample(ray)
		return c
	}
}
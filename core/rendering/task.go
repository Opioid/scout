package rendering

import (
	_ "github.com/Opioid/scout/core/rendering/film"
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

func (t *Task) render(scene *pkgscene.Scene, camera camera.Camera, shutterOpen, shutterClose float32, start, end math.Vector2i, sampler pkgsampler.Sampler) {
	f := camera.Film()

	numSamples := sampler.NumSamplesPerIteration()

//	var offset math.Vector2
	var sample pkgsampler.CameraSample
	var ray math.OptimizedRay

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			sampler.Restart(1)
			t.integrator.StartNewPixel(numSamples)
			sampleId := uint32(0)

			for sampler.GenerateNewSample(math.MakeVector2(float32(x), float32(y)), &sample) {
			//	sample.Coordinates = math.MakeVector2(float32(x) + offset.X, float32(y) + offset.Y)
			//	sample.LensUv = math.MakeVector2(sample.Coordinates.X / dx, sample.Coordinates.Y / dy)
			//	sample.RelativeOffset = offset.SubS(0.5)

				camera.GenerateRay(&sample, shutterOpen, shutterClose, &ray)

				color := t.Li(scene, sampleId, &ray) 

				f.AddSample(&sample, color, start, end)

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
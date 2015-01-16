package rendering

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Worker struct {
	integrator Integrator
} 

func makeWorker(integrator Integrator) Worker {
	w := Worker{integrator}
	return w
}

func (w *Worker) render(scene *pkgscene.Scene, camera camera.Camera, shutterOpen, shutterClose float32, start, end math.Vector2i, sampler pkgsampler.Sampler) {
	f := camera.Film()

	numSamples := sampler.NumSamplesPerIteration()

	var sample pkgsampler.CameraSample
	var ray math.OptimizedRay

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			sampler.Restart(1)
			w.integrator.StartNewPixel(numSamples)
			sampleId := uint32(0)

			for sampler.GenerateNewSample(math.MakeVector2(float32(x), float32(y)), &sample) {
				camera.GenerateRay(&sample, shutterOpen, shutterClose, &ray)

				color := w.Li(sampleId, scene, &ray) 

				f.AddSample(&sample, color, start, end)

				sampleId++
			}
		}
	}
}

func (w *Worker) Li(subsample uint32, scene *pkgscene.Scene, ray *math.OptimizedRay) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		c := w.integrator.Li(w, subsample, scene, ray, &intersection) 
		return c
	} else {
		c := scene.Surrounding.Sample(ray)
		return c
	}
}
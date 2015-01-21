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
	sample pkgsampler.CameraSample
	ray math.OptimizedRay
	intersections []prop.Intersection
} 

func makeWorker(integrator Integrator) Worker {
	w := Worker{integrator: integrator}

	// To reduce strain on the GC the max amount of intersection is alloceted once only.
	// I believe it would be possible to design the integrators so that only one intersection is ever needed regardless of depth.
	w.intersections = make([]prop.Intersection, integrator.MaxBounces() + 1)

	return w
}

func (w *Worker) render(scene *pkgscene.Scene, camera camera.Camera, shutterOpen, shutterClose float32, start, end math.Vector2i, sampler pkgsampler.Sampler) {
	f := camera.Film()

	numSamples := sampler.NumSamplesPerIteration()

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			sampler.Restart(1)
			w.integrator.StartNewPixel(numSamples)
			sampleId := uint32(0)
			offset := math.MakeVector2(float32(x), float32(y))

			for sampler.GenerateCameraSample(offset, &w.sample) {
				camera.GenerateRay(&w.sample, shutterOpen, shutterClose, &w.ray)

				color := w.Li(sampleId, scene, &w.ray) 

				f.AddSample(&w.sample, color, start, end)

				sampleId++
			}
		}
	}
}

func (w *Worker) Li(subsample uint32, scene *pkgscene.Scene, ray *math.OptimizedRay) math.Vector3 {
	intersection := &w.intersections[ray.Depth]

	if scene.Intersect(ray, intersection) {
		c := w.integrator.Li(w, subsample, scene, ray, intersection) 
		return c
	} else {
		c := scene.Surrounding.Sample(ray)
		return c
	}
}
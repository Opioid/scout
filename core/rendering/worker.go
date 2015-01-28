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

	Scene *pkgscene.Scene

	ScratchBuffer prop.ScratchBuffer
} 

func makeWorker(integrator Integrator) Worker {
	w := Worker{integrator: integrator}

	// To reduce strain on the GC the max amount of intersection is allocated once only.
	// I believe it would be possible to design the integrators so that only one intersection is ever needed regardless of depth.
	w.intersections = make([]prop.Intersection, integrator.MaxBounces() + 1)

	return w
}

func (w *Worker) render(scene *pkgscene.Scene, camera camera.Camera, shutterOpen, shutterClose float32, start, end math.Vector2i, sampler pkgsampler.Sampler) {
	w.Scene = scene

	f := camera.Film()

	numSamples := sampler.NumSamplesPerIteration()

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			sampler.Restart(1)
			w.integrator.StartNewPixel(numSamples)
			sampleId := uint32(0)
			offset := math.MakeVector2(float32(x), float32(y))

			for sampler.GenerateCameraSample(offset, &w.sample) {
				camera.GenerateRay(&w.sample, shutterOpen, shutterClose, &w.ScratchBuffer.Transformation, &w.ray)

				color := w.Li(sampleId, &w.ray) 

				f.AddSample(&w.sample, color, start, end)

				sampleId++
			}
		}
	}
}

func (w *Worker) Li(subsample uint32, ray *math.OptimizedRay) math.Vector3 {
	intersection := &w.intersections[ray.Depth]

	var visibility uint8

	if ray.Depth == 0 {
		visibility = w.integrator.PrimaryVisibility()
	} else {
		visibility = w.integrator.SecondaryVisibility()
	}

	if w.Scene.Intersect(ray, visibility, &w.ScratchBuffer, intersection) {
		c := w.integrator.Li(w, subsample, ray, intersection) 
		return c
	} else {
		c := w.Scene.Surrounding.Sample(ray)
		return c
	}
}

func (w *Worker) Shadow(ray *math.OptimizedRay) bool {
	return w.Scene.IntersectP(ray, &w.ScratchBuffer)
}
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

type Tile struct {
	renderer *Renderer
	integrator Integrator
} 

func makeTile(renderer *Renderer, integrator Integrator) Tile {
	t := Tile{}
	t.renderer = renderer
	t.integrator = integrator
	return t
}

func (t *Tile) render(scene *pkgscene.Scene, camera camera.Camera, shutterOpen, shutterClose float32, start, end math.Vector2i, sampler pkgsampler.Sampler) {
	f := camera.Film()

	numSamples := sampler.NumSamplesPerIteration()

	var sample pkgsampler.CameraSample
	var ray math.OptimizedRay

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			sampler.Restart(1)
			t.integrator.StartNewPixel(numSamples)
			sampleId := uint32(0)

			for sampler.GenerateNewSample(math.MakeVector2(float32(x), float32(y)), &sample) {
					camera.GenerateRay(&sample, shutterOpen, shutterClose, &ray)

				color := t.Li(scene, sampleId, &ray) 

				f.AddSample(&sample, color, start, end)

				sampleId++
			}
		}
	}
}

func (t *Tile) Li(scene *pkgscene.Scene, subsample uint32, ray *math.OptimizedRay) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		c := t.integrator.Li(scene, t, subsample, ray, &intersection) 
		return c
	} else {
		c := scene.Surrounding.Sample(ray)
		return c
	}
}
package rendering

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type RenderTask struct {
	renderer *Renderer
	integrator Integrator
} 

func (r *RenderTask) FirstSample(numSamples uint32) {
	r.integrator.FirstSample(numSamples)
}

func (r *RenderTask) Li(scene *pkgscene.Scene, subsample uint32, ray *math.OptimizedRay) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		c := r.integrator.Li(scene, r, subsample, ray, &intersection) 
		return c
	} else {
		c := scene.Surrounding.Sample(ray)
		return c
	}
}
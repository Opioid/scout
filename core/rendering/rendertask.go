package rendering

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
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
		return r.integrator.Li(scene, r, subsample, ray, &intersection) 
	} else {
		c, _ := scene.Surrounding.Sample(ray)
		return c
	}
}
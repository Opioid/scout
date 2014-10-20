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

func (r *RenderTask) Li(scene *pkgscene.Scene, sample, numSamples uint32, ray *math.OptimizedRay) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		return r.integrator.Li(scene, r, sample, numSamples, ray, &intersection) 
	} else {
		return scene.Surrounding.Sample(ray)
	}
}
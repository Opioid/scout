package rendering

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type RenderTask struct {
	renderer *Renderer
	integrator Integrator
} 

func (r *RenderTask) Li(scene *pkgscene.Scene, sample, numSamples uint32, ray *math.OptimizedRay, rng *random.Generator) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		return r.integrator.Li(scene, r, sample, numSamples, ray, &intersection, rng) 
	} else {
		return scene.Surrounding.Sample(ray)
	}
}
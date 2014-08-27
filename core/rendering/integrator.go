package rendering

import (
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Integrator interface {
	Li(scene *scene.Scene, renderer *Renderer, ray *math.OptimizedRay, intersection *prop.Intersection, rng *random.Generator) math.Vector3
}
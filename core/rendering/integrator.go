package rendering

import (
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Integrator interface {
	StartNewPixel(numSamples uint32)

	Li(w *Worker, subsample uint32, scene *scene.Scene, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3

	MaxBounces() uint32
}

type IntegratorFactory interface {
	New(rng *random.Generator) Integrator
}
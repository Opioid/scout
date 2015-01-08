package rendering

import (
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Integrator interface {
	StartNewPixel(numSamples uint32)

	Li(scene *scene.Scene, tile *Tile, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3
}

type IntegratorFactory interface {
	New(rng *random.Generator) Integrator
}
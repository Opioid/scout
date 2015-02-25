package rendering

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Integrator interface {
	StartNewPixel(numSamples uint32)

	Li(w *Worker, subsample uint32, ray *math.OptimizedRay, intersection *prop.Intersection) math.Vector3

	PrimaryVisibility() uint8
	SecondaryVisibility() uint8
}

type IntegratorFactory interface {
	New(id uint32, rng *random.Generator) Integrator
}
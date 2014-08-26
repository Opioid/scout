package integrator

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Integrator interface {
	Li(scene *pkgscene.Scene, ray *math.OptimizedRay, rng *random.Generator) math.Vector3
}
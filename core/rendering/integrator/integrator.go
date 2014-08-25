package integrator

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/base/math"
)

type Integrator interface {
	Li(scene *pkgscene.Scene, ray *math.OptimizedRay) math.Vector3
}
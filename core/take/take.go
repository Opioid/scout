package take

import (
	"github.com/Opioid/scout/core/rendering"
)

type Take struct {
	Scene string
	Context rendering.Context
	SurfaceIntegratorFactory rendering.SurfaceIntegratorFactory
}
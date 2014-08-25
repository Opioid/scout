package take

import (
	"github.com/Opioid/scout/core/rendering"
	"github.com/Opioid/scout/core/rendering/integrator"
)

type Take struct {
	Scene  string
	Context rendering.Context
	Integrator integrator.Integrator
}
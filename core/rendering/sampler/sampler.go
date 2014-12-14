package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Sampler interface {
	Clone(rng *random.Generator) Sampler 
	NumSamplesPerPixel() uint32

	Restart()

	GenerateNewSample(sample *math.Vector2) bool
}
package integrator

import (
	"github.com/Opioid/scout/base/math/random"
)

type Integrator struct {
	ID uint32
	Rng *random.Generator
}
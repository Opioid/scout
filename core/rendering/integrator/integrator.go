package integrator

import (
	"github.com/Opioid/scout/base/math/random"
)

type integrator struct {
	id uint32
	rng *random.Generator
}
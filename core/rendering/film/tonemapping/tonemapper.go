package tonemapping

import (
	"github.com/Opioid/scout/base/math"
)

type Tonemapper interface {
	Tonemap(color math.Vector3) math.Vector3
}
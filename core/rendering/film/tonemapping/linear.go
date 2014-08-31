package tonemapping 

import (
	"github.com/Opioid/scout/base/math"
)

type linear struct {

}

func NewLinear() *linear {
	return new(linear)
}

func (l *linear) Tonemap(color math.Vector3) math.Vector3 {
	return color
}
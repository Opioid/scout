package tonemapping 

import (
	"github.com/Opioid/scout/base/math"
)

type identity struct {

}

func NewIdentity() *identity {
	return new(identity)
}

func (i *identity) Tonemap(color math.Vector3) math.Vector3 {
	return color
}
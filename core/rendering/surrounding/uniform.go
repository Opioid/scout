package surrounding

import (
	"github.com/Opioid/scout/base/math"
)

type uniform struct {
	color math.Vector3
}

func NewUniform(color math.Vector3) *uniform {
	return &uniform{color}
}

func (u *uniform) Sample(ray *math.OptimizedRay) math.Vector3 {
	return u.color
}

func (u *uniform) SampleSecondary(ray *math.OptimizedRay) (math.Vector3, float32) {
	return u.color, 1
}

func (u *uniform) SampleDiffuse(v math.Vector3) math.Vector3 {
	return u.color
}

func (u *uniform) SampleSpecular(v math.Vector3, roughness float32) math.Vector3 {
	return u.color
}
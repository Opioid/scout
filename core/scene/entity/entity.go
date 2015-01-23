package entity

import (
	"github.com/Opioid/scout/base/math"
)

type Entity struct {
	Transformation math.ComposedTransformation
	Animation Animation
}

func NewEntity() *Entity {
	return &Entity{}
}

func (e *Entity) TransformationAt(time float32, transformation *math.ComposedTransformation)  {
	if !e.Animation.empty() {
		e.Animation.at1(time, transformation)
	} else {
		*transformation = e.Transformation
	}
}

func (e *Entity) SetTransformation(position, scale math.Vector3, rotation math.Quaternion) {
	e.Transformation.Set(position, scale, rotation)
}
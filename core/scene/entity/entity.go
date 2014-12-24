package entity

import (
	"github.com/Opioid/scout/base/math"
)

type Entity struct {
	transformation ComposedTransformation
	Animation Animation
}

func (e *Entity) TransformationAt(time float32) ComposedTransformation {
	if !e.Animation.empty() {
		return e.Animation.at(time)
	} 

	return e.transformation
}

func (e *Entity) SetTransformation(position, scale math.Vector3, rotation math.Quaternion) {
	e.transformation.Set(position, scale, rotation)
}
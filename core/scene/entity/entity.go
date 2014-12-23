package entity

import (
	
)

type Entity struct {
	Transformation ComposedTransformation
	animation animation
}

func (e *Entity) TransformationAt(time float32) *ComposedTransformation {
	if !e.animation.empty() {
		
	} 

	return &e.Transformation
}
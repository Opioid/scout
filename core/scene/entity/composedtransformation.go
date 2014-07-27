package entity

import (
	"github.com/Opioid/scout/base/math"
)

type ComposedTransformation struct {
	math.Transformation
	Matrix math.Matrix4x4
}

/*
func (t *ComposedTransformation) Update(transformation *math.Transformation) {
	t.Transformation = *transformation
}
*/

func (t *ComposedTransformation) Update() {
	rotation := math.NewMatrix3x3FromQuaternion(t.Rotation)
	t.Matrix.SetBasis(rotation)
	t.Matrix.SetOrigin(t.Position)
}
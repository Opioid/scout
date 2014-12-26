package entity

import (
	"github.com/Opioid/scout/base/math"
)

type ComposedTransformation struct {
	Position math.Vector3
	Scale math.Vector3
	Rotation math.Matrix3x3
	ObjectToWorld math.Matrix4x4
	WorldToObject math.Matrix4x4
}

func (t *ComposedTransformation) Set(position, scale math.Vector3, rotation math.Quaternion) {
	t.Position = position
	t.Scale = scale
	t.Rotation.SetFromQuaternion(rotation)

	t.ObjectToWorld.SetBasis(&t.Rotation)
	t.ObjectToWorld.Scale(scale)
	t.ObjectToWorld.SetOrigin(position)

	t.WorldToObject = t.ObjectToWorld.Inverted()
}

func (t *ComposedTransformation) SetFromTransformation(transformation *math.Transformation) {
	t.Set(transformation.Position, transformation.Scale, transformation.Rotation)
}
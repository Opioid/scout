package math

import (
	_ "fmt"
)

type ComposedTransformation struct {
	Position Vector3
	Scale Vector3
	Rotation Matrix3x3
	ObjectToWorld Matrix4x4
	WorldToObject Matrix4x4
}

func MakeComposedTransformation(transformation *Transformation) ComposedTransformation {
	rotation := MakeMatrix3x3FromQuaternion(transformation.Rotation)
	objectToWorld := MakeMatrix4x4FromBasisScaleOrigin(&rotation, transformation.Scale, transformation.Position)

	return ComposedTransformation{
		transformation.Position,
		transformation.Scale,
		rotation,
		objectToWorld,
		objectToWorld.Inverted(),
	}
}

func (t *ComposedTransformation) Set(position, scale Vector3, rotation Quaternion) {
	t.Position = position
	t.Scale = scale
	t.Rotation.SetFromQuaternion(rotation)
	t.ObjectToWorld.SetFromBasisScaleOrigin(&t.Rotation, scale, position)
	t.WorldToObject = t.ObjectToWorld.Inverted()
}

func (t *ComposedTransformation) SetFromTransformation(transformation *Transformation) {
	t.Set(transformation.Position, transformation.Scale, transformation.Rotation)
}
package entity

import (
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type ComposedTransformation struct {
	Position math.Vector3
	Scale math.Vector3
	Rotation math.Matrix3x3
	ObjectToWorld math.Matrix4x4
	WorldToObject math.Matrix4x4
}

func MakeComposedTransformation(transformation math.Transformation) ComposedTransformation {
	rotation := math.MakeMatrix3x3FromQuaternion(transformation.Rotation)
	objectToWorld := math.MakeMatrix4x4FromBasisScaleOrigin(rotation, transformation.Scale, transformation.Position)

	return ComposedTransformation{
		transformation.Position,
		transformation.Scale,
		rotation,
		objectToWorld,
		objectToWorld.Inverted(),
	}
}

func (t *ComposedTransformation) Set(position, scale math.Vector3, rotation math.Quaternion) {
	t.Position = position
	t.Scale = scale
	t.Rotation.SetFromQuaternion(rotation)

	t.ObjectToWorld.SetIdentity()
	t.ObjectToWorld.SetBasis(&t.Rotation)
	t.ObjectToWorld.Scale(scale)
	t.ObjectToWorld.SetOrigin(position)

	t.WorldToObject = t.ObjectToWorld.Inverted()
}

func (t *ComposedTransformation) SetFromTransformation(transformation math.Transformation) {
	t.Set(transformation.Position, transformation.Scale, transformation.Rotation)
}
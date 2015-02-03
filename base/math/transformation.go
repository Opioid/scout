package math

type Transformation struct {
	Position Vector3
	Scale    Vector3
	Rotation Quaternion
}

func (trans *Transformation) Lerp(other *Transformation, t float32) Transformation {
	r := Transformation{}
	r.Position = trans.Position.Lerp(other.Position, t)
	r.Scale = trans.Scale.Lerp(other.Scale, t)

//	r.Position = Vector3Lerp(trans.Position, other.Position, t)
//	r.Scale = Vector3Lerp(trans.Scale, other.Scale, t)

	r.Rotation = trans.Rotation.Slerp(other.Rotation, t)
	return r
}
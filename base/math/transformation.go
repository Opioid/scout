package math

type Transformation struct {
	Position Vector3
	Scale    Vector3
	Rotation Quaternion
}

func (this *Transformation) Lerp(other *Transformation, t float32) Transformation {
	r := Transformation{}
	r.Position = this.Position.Lerp(other.Position, t)
	r.Scale = this.Scale.Lerp(other.Scale, t)
	r.Rotation = this.Rotation.Slerp(other.Rotation, t)
	return r
}
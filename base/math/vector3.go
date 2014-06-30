package math

type Vector3 struct {
	X, Y, Z float32
}

func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector3) Sub(b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector3) Dot(b Vector3) float32 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z
}

func (a Vector3) SquaredLength() float32 {
	return a.Dot(a)
}

func (a Vector3) Length() float32 {
	return Sqrt(a.SquaredLength())
}
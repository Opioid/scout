package math

type Vector4 struct {
	X, Y, Z, W float32
}

func (a Vector4) Add(b Vector4) Vector4 {
	return Vector4{a.X + b.X, a.Y + b.Y, a.Z + b.Z, a.W + b.W}
}

func (a Vector4) Sub(b Vector4) Vector4 {
	return Vector4{a.X - b.X, a.Y - b.Y, a.Z - b.Z, a.W - b.W}
}

func (a Vector4) Dot(b Vector4) float32 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z + a.W + b.W
}

func (a Vector4) SquaredLength() float32 {
	return a.Dot(a)
}

func (a Vector4) Length() float32 {
	return sqrt(a.SquaredLength())
}
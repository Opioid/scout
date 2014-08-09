package math

type Vector3 struct {
	X, Y, Z float32
}

func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a *Vector3) AddAssign(b Vector3) Vector3 {
	a.X += b.X
	a.Y += b.Y
	a.Z += b.Z
	return *a
}

func (a Vector3) Sub(b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector3) Scale(s float32) Vector3 {
	return Vector3{a.X * s, a.Y * s, a.Z * s}
}

func (a Vector3) Mul(b Vector3) Vector3 {
	return Vector3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

func (a Vector3) Div(s float32) Vector3 {
	return a.Scale(1.0 / s)
}

func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{
		a.Y * b.Z - a.Z * b.Y,
		a.Z * b.X - a.X * b.Z,
		a.X * b.Y - a.Y * b.X,
	}
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

func (a Vector3) Normalized() Vector3 {
	return a.Div(a.Length())
}

func (a Vector3) Saturated() Vector3 {
	return Vector3{Clamp(a.X, 0.0, 1.0), Clamp(a.Y, 0.0, 1.0), Clamp(a.Z, 0.0, 1.0)}
}

func (a Vector3) Reflect(b Vector3) Vector3 {
	return b.Sub(a.Scale(2.0 * b.Dot(a)))
}

func (a Vector3) Min(b Vector3) Vector3 {
	return Vector3{Min(a.X, b.X), Min(a.Y, b.Y), Min(a.Z, b.Z)}
}

func (a Vector3) Max(b Vector3) Vector3 {
	return Vector3{Max(a.X, b.X), Max(a.Y, b.Y), Max(a.Z, b.Z)}
}
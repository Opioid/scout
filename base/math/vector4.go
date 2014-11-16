package math

type Vector4 struct {
	X, Y, Z, W float32
}

func MakeVector4(x, y, z, w float32) Vector4 {
	return Vector4{x, y, z, w}
}

func (v Vector4) Vector3() Vector3 {
	return Vector3{v.X, v.Y, v.Z}
}

func (a Vector4) At(i int) float32 {
	switch i {
	case 0:
		return a.X
	case 1:
		return a.Y
	case 2:
		return a.Z
	default:
		return a.W
	}
}

func (a *Vector4) Set(i int, value float32) {
	switch i {
	case 0:
		a.X = value
	case 1:
		a.Y = value
	case 2:
		a.Z = value
	default:
		a.W = value
	}
}

func (a Vector4) Add(b Vector4) Vector4 {
	return Vector4{a.X + b.X, a.Y + b.Y, a.Z + b.Z, a.W + b.W}
}

func (a *Vector4) AddAssign(b Vector4) Vector4 {
	a.X += b.X
	a.Y += b.Y
	a.Z += b.Z
	a.W += b.W
	return *a
}

func (a Vector4) Sub(b Vector4) Vector4 {
	return Vector4{a.X - b.X, a.Y - b.Y, a.Z - b.Z, a.W - b.W}
}

func (a Vector4) Scale(s float32) Vector4 {
	return Vector4{a.X * s, a.Y * s, a.Z * s, a.W * s}
}

func (a Vector4) Dot(b Vector4) float32 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z + a.W + b.W
}

func (a Vector4) SquaredLength() float32 {
	return a.Dot(a)
}

func (a Vector4) Length() float32 {
	return Sqrt(a.SquaredLength())
}

func (a Vector4) Lerp(b Vector4, t float32) Vector4 {
	_t := 1.0 - t
	return Vector4{_t * a.X + t * b.X, _t * a.Y + t * b.Y, _t * a.Z + t * b.Z, _t * a.W + t * b.W}
}
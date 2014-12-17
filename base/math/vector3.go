package math

import (
	"github.com/Opioid/math32"
)

type Vector3 struct {
	X, Y, Z float32
}

func MakeVector3(x, y, z float32) Vector3 {
	return Vector3{x, y, z}
}

func MakeIdentityVector3() Vector3 {
	return Vector3{1, 1, 1}
}

func MakeVector3All(s float32) Vector3 {
	return Vector3{s, s, s}
}

func (a Vector3) At(i int32) float32 {
	switch i {
	case 0:
		return a.X
	case 1:
		return a.Y
	default:
		return a.Z
	}
}

func (a *Vector3) Set(i int32, value float32) {
	switch i {
	case 0:
		a.X = value
	case 1:
		a.Y = value
	default:
		a.Z = value
	}
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

func (a Vector3) AddS(s float32) Vector3 {
	return Vector3{a.X + s, a.Y + s, a.Z + s}
}

func (a Vector3) Sub(b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector3) SubS(s float32) Vector3 {
	return Vector3{a.X - s, a.Y - s, a.Z - s}
}

func (a Vector3) Scale(s float32) Vector3 {
	return Vector3{a.X * s, a.Y * s, a.Z * s}
}

func (a *Vector3) ScaleAssign(s float32) Vector3 {
	a.X *= s
	a.Y *= s
	a.Z *= s
	return *a
}

func (a Vector3) Mul(b Vector3) Vector3 {
	return Vector3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

func (a Vector3) Div(s float32) Vector3 {
	return a.Scale(1 / s)
}

func (a Vector3) DivV(b Vector3) Vector3 {
	return Vector3{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
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
	return math32.Sqrt(a.SquaredLength())
}

func (a Vector3) SquaredDistance(b Vector3) float32 {
	return b.Sub(a).SquaredLength()
}

func (a Vector3) Normalized() Vector3 {
//	return a.Div(a.Length())
	rl := math32.Rsqrt(a.X * a.X + a.Y * a.Y + a.Z * a.Z)
	return Vector3{a.X * rl, a.Y * rl, a.Z * rl}
}

func (a Vector3) Saturated() Vector3 {
	return Vector3{math32.Clamp(a.X, 0, 1), math32.Clamp(a.Y, 0, 1), math32.Clamp(a.Z, 0, 1)}
}

func (a Vector3) Reflect(b Vector3) Vector3 {
	return b.Sub(a.Scale(2 * b.Dot(a)))
}

func (a Vector3) Min(b Vector3) Vector3 {
	return Vector3{math32.Min(a.X, b.X), math32.Min(a.Y, b.Y), math32.Min(a.Z, b.Z)}
}

func (a Vector3) Max(b Vector3) Vector3 {
	return Vector3{math32.Max(a.X, b.X), math32.Max(a.Y, b.Y), math32.Max(a.Z, b.Z)}
}

func (a Vector3) Lerp(b Vector3, t float32) Vector3 {
	_t := 1 - t
	return Vector3{_t * a.X + t * b.X, _t * a.Y + t * b.Y, _t * a.Z + t * b.Z}
}

func (a Vector3) ContainsNaN() bool {
	return math32.IsNaN(a.X) || math32.IsNaN(a.Y) || math32.IsNaN(a.Z)
}

func (a Vector3) ContainsInf() bool {
	return IsInf(a.X) || IsInf(a.Y) || IsInf(a.Z)
}

type Vector3i struct {
	X, Y, Z int32
}

func MakeVector3i(x, y, z int32) Vector3i {
	return Vector3i{x, y, z}
}

func (a Vector3i) Vector2i() Vector2i {
	return Vector2i{a.X, a.Y}
}

func (a Vector3i) Add(b Vector3i) Vector3i {
	return Vector3i{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector3i) Min(b Vector3i) Vector3i {
	return Vector3i{Mini(a.X, b.X), Mini(a.Y, b.Y), Mini(a.Z, b.Z)}
}

func (a Vector3i) Max(b Vector3i) Vector3i {
	return Vector3i{Maxi(a.X, b.X), Maxi(a.Y, b.Y), Maxi(a.Z, b.Z)}
}
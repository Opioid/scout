package math

type Vector2 struct {
	X, Y float32
}

func MakeVector2(x, y float32) Vector2 {
	return Vector2{x, y}
}

func (a Vector2) Add(b Vector2) Vector2 {
	return Vector2{a.X + b.X, a.Y + b.Y}
}

func (a Vector2) SubS(s float32) Vector2 {
	return Vector2{a.X - s, a.Y - s}
}

func (a Vector2) Scale(s float32) Vector2 {
	return Vector2{a.X * s, a.Y * s}
}

func (a Vector2) Div(s float32) Vector2 {
	return a.Scale(1 / s)
}

type Vector2i struct {
	X, Y int32
}

func MakeVector2i(x, y int32) Vector2i {
	return Vector2i{x, y}
}

func (a Vector2i) Add(b Vector2i) Vector2i {
	return Vector2i{a.X + b.X, a.Y + b.Y}
}

func (a Vector2i) Min(b Vector2i) Vector2i {
	return Vector2i{Mini(a.X, b.X), Mini(a.Y, b.Y)}
}

func (a Vector2i) Max(b Vector2i) Vector2i {
	return Vector2i{Maxi(a.X, b.X), Maxi(a.Y, b.Y)}
}
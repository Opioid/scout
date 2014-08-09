package math

type Vector2 struct {
	X, Y float32
}

type Vector2i struct {
	X, Y int
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
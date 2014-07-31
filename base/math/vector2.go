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
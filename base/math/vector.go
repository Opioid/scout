package math

type Vector3 struct {
	X, Y, Z float32
}

func V3Add(a, b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func V3Sub(a, b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func V3Dot(a, b Vector3) float32 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z
}
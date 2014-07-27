package math

type Ray struct {
	Origin, Direction Vector3
	MinT, MaxT float32
}

func (r *Ray) Point(t float32) Vector3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
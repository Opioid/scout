package math

type Plane struct {
	A, B, C, D float32
}

func MakePlane(normal, point Vector3) Plane {
	return Plane{A: normal.X, B: normal.Y, C: normal.Z, D: -normal.Dot(point)}
}

func (p *Plane) Behind(v Vector3) bool {
	d := p.A * v.X + p.B * v.Y + p.C * v.Z + p.D;
	return d < 0.0
}
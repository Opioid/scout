package shape

import (
	"github.com/Opioid/scout/base/math"
)

type Sphere struct {
	Position math.Vector3
	Radius float32
}

// Won't work from the inside!
func (s *Sphere) Intersect(ray *math.Ray, thit *float32) bool {
	v := ray.Origin.Sub(s.Position)
	b := -v.Dot(ray.Direction)
	det := (b * b) - v.Dot(v) + (s.Radius * s.Radius)

	if det > 0.0 {
		*thit = b - math.Sqrt(det)

		if *thit >= 0.0 && *thit < ray.MaxT {
			return true
		} 
	}

	return false
}
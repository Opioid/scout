package camera

import (
	"github.com/Opioid/scout/base/math"
)

type Orthographic struct {
	genericCamera
	extent math.Vector2
}

func NewOrthographic(dimensions math.Vector2, film Film) *Orthographic {
	o := new(Orthographic)
	o.film = film
	
	if 0.0 == dimensions.X {
		o.extent = math.Vector2{dimensions.Y * (float32(film.Dimensions.X) / float32(film.Dimensions.Y)), dimensions.Y}
	} else if 0.0 == dimensions.Y {
		o.extent = math.Vector2{dimensions.X, dimensions.X * (float32(film.Dimensions.Y) / float32(film.Dimensions.X))}
	} else {
		o.extent = dimensions
	}

	return o
}

func (o *Orthographic) Position() math.Vector3 {
	return o.Entity.Transformation.Position
}

func (o *Orthographic) Film() Film {
	return o.film
}

func (o *Orthographic) GenerateRay(sample *Sample, ray *math.Ray) {
	x := sample.Coordinates.X / float32(o.film.Dimensions.X)
	y := sample.Coordinates.Y / float32(o.film.Dimensions.Y)

	offset := math.Vector3{x * o.extent.X - 0.5 * o.extent.X, 0.5 * o.extent.Y - y * o.extent.Y, 0.0}

	ray.Origin    = o.Entity.Transformation.Position.Add(offset)
	ray.Direction = math.Vector3{0.0, 0.0, 1.0}
	ray.MaxT      = 1000.0
}
package camera

import (
	"github.com/Opioid/scout/base/math"
)

type Orthographic struct {
	genericCamera
	dimensions math.Vector2i
}

func NewOrthographic(dimensions math.Vector2i, film Film) *Orthographic {
	o := new(Orthographic)
	o.dimensions = dimensions
	o.film = film
	return o
}

func (o *Orthographic) Position() math.Vector3 {
	return o.Entity.Position
}

func (o *Orthographic) Film() Film {
	return o.film
}

func (o *Orthographic) GenerateRay(sample *Sample, ray *math.Ray) {
	width  := float32(4.0)
	height := width * (float32(o.film.Dimensions.Y) / float32(o.film.Dimensions.X))

	x := float32(sample.Coordinates.X) / float32(o.film.Dimensions.X)
	y := float32(sample.Coordinates.Y) / float32(o.film.Dimensions.Y)

	offset := math.Vector3{x * width - 0.5 * width, 0.5 * height - y * height, 0.0}

	ray.Origin    = o.Entity.Position.Add(offset)
	ray.Direction = math.Vector3{0.0, 0.0, 1.0}
	ray.MaxT      = 1000.0
}
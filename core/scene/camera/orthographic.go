package camera

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Orthographic struct {
	genericCamera
}

func NewOrthographic(dimensions math.Vector2, film film.Film) *Orthographic {
	o := new(Orthographic)
	o.film = film
	o.dimensions = calculateDimensions(dimensions, film)
	o.Entity.Transformation.ObjectToWorld.SetIdentity()
	return o
}

func (o *Orthographic) UpdateView() {}

func (o *Orthographic) Transformation() *entity.ComposedTransformation {
	return &o.Entity.Transformation
}

func (o *Orthographic) Position() math.Vector3 {
	return o.Entity.Transformation.Position
}

func (o *Orthographic) Film() film.Film {
	return o.film
}

func (o *Orthographic) GenerateRay(sample *sampler.Sample, ray *math.OptimizedRay) {
	x := sample.Coordinates.X / float32(o.film.Dimensions().X)
	y := sample.Coordinates.Y / float32(o.film.Dimensions().Y)

	offset := math.MakeVector3(x * o.dimensions.X - 0.5 * o.dimensions.X, 0.5 * o.dimensions.Y - y * o.dimensions.Y, 0.0)
	offset = o.Entity.Transformation.Rotation.TransformVector3(offset)
	ray.Origin = o.Entity.Transformation.Position.Add(offset)

	ray.SetDirection(o.Entity.Transformation.Rotation.Row(2))

	ray.MaxT  = 1000.0
	ray.Depth = 0
}
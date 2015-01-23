package camera

import (
	"github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Orthographic struct {
	projectiveCamera
}

func NewOrthographic(dimensions math.Vector2, film film.Film, shutterSpeed float32) *Orthographic {
	o := new(Orthographic)
	o.dimensions = calculateDimensions(dimensions, film)
	o.film = film
//	o.Entity.Transformation.ObjectToWorld.SetIdentity()
	o.shutterSpeed = shutterSpeed
	return o
}

func (o *Orthographic) UpdateView() {}

func (o *Orthographic) GenerateRay(sample *sampler.CameraSample, shutterOpen, shutterClose float32, transformation *math.ComposedTransformation, ray *math.OptimizedRay) {
	x := sample.Coordinates.X / float32(o.film.Dimensions().X)
	y := sample.Coordinates.Y / float32(o.film.Dimensions().Y)

	offset := math.MakeVector3(x * o.dimensions.X - 0.5 * o.dimensions.X, 0.5 * o.dimensions.Y - y * o.dimensions.Y, 0.0)

	o.entity.TransformationAt(sample.Time, transformation)

	offset = transformation.Rotation.TransformVector3(offset)
	ray.Origin = transformation.Position.Add(offset)

	ray.SetDirection(transformation.Rotation.Direction())

	ray.MaxT  = 1000.0
	ray.Depth = 0
}
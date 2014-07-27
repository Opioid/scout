package camera

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Sample struct {
	Coordinates math.Vector2
}

func NewSample(x, y float32) *Sample {
	return &Sample{math.Vector2{x, y}}
}

type Camera interface {
	UpdateView()
	Transformation() *entity.ComposedTransformation
	Position() math.Vector3
	Film() Film
	GenerateRay(sample *Sample, ray *math.Ray)
}

type genericCamera struct {
	entity.Entity
	dimensions math.Vector2
	film Film
}

func calculateDimensions(dimensions math.Vector2, film Film) math.Vector2 {
	var r math.Vector2

	if 0.0 == dimensions.X {
		r = math.Vector2{dimensions.Y * (float32(film.Dimensions.X) / float32(film.Dimensions.Y)), dimensions.Y}
	} else if 0.0 == dimensions.Y {
		r = math.Vector2{dimensions.X, dimensions.X * (float32(film.Dimensions.Y) / float32(film.Dimensions.X))}
	} else {
		r = dimensions
	}

	return r
}
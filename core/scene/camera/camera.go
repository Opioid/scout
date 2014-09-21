package camera

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Camera interface {
	UpdateView()
	Transformation() *entity.ComposedTransformation
	Position() math.Vector3
	Film() film.Film
	GenerateRay(sample *sampler.Sample, ray *math.OptimizedRay)
}

type genericCamera struct {
	entity.Entity
	dimensions math.Vector2
	film film.Film
}

func calculateDimensions(dimensions math.Vector2, film film.Film) math.Vector2 {
	var r math.Vector2

	if 0.0 == dimensions.X {
		r = math.MakeVector2(dimensions.Y * (float32(film.Dimensions().X) / float32(film.Dimensions().Y)), dimensions.Y)
	} else if 0.0 == dimensions.Y {
		r = math.MakeVector2(dimensions.X, dimensions.X * (float32(film.Dimensions().Y) / float32(film.Dimensions().X)))
	} else {
		r = dimensions
	}

	return r
}
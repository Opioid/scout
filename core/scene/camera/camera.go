package camera

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Camera interface {
	Entity() *entity.Entity
	UpdateView()
	Film() film.Film
	GenerateRay(sample *sampler.CameraSample, shutterOpen, shutterClose float32, ray *math.OptimizedRay)
}

type projectiveCamera struct {
	entity entity.Entity
	dimensions math.Vector2
	film film.Film
}

func (p *projectiveCamera) Entity() *entity.Entity {
	return &p.entity
}

func (p *projectiveCamera) Film() film.Film {
	return p.film
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
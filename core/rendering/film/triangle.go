package film

import (
	"github.com/Opioid/scout/core/rendering/film/filter"
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_"fmt"
)

type Triangle struct {
	film
	filter filter.Triangle
}

func NewTriangle(dimensions math.Vector2i, exposure float32, tonemapper tonemapping.Tonemapper) *Triangle {
	f := new(Triangle)
	f.resize(dimensions)
	f.exposure = exposure
	f.tonemapper = tonemapper

	radius := math32.Sqrt(2)
	f.filter.SetWidth(math.MakeVector2(radius, radius))
	return f
}

func (f *Triangle) AddSample(sample *sampler.Sample, color math.Vector3) {
	x, y := int32(sample.Coordinates.X), int32(sample.Coordinates.Y)

	w := f.filter.Evaluate(sample.RelativeOffset)

	f.addPixel(x, y, color, w)
}
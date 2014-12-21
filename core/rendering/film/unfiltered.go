package film

import (
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	_"fmt"
)

type Unfiltered struct {
	film
}

func NewUnfiltered(dimensions math.Vector2i, exposure float32, tonemapper tonemapping.Tonemapper) *Unfiltered {
	f := new(Unfiltered)
	f.resize(dimensions)
	f.exposure = exposure
	f.tonemapper = tonemapper
	return f
}

func (f *Unfiltered) AddSample(sample *sampler.CameraSample, color math.Vector3, start, end math.Vector2i) {
	x, y := int32(sample.Coordinates.X), int32(sample.Coordinates.Y)

	f.addPixel(x, y, color, 1)
}

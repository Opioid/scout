package film

import (
	"github.com/Opioid/scout/core/rendering/film/filter"
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	_"github.com/Opioid/math32"
	_ "fmt"
)

type Filtered struct {
	film
	filter filter.Triangle
}

func NewFiltered(dimensions math.Vector2i, exposure float32, tonemapper tonemapping.Tonemapper) *Filtered {
	f := new(Filtered)
	f.resize(dimensions)
	f.exposure = exposure
	f.tonemapper = tonemapper

	radius := float32(1)//*/math32.Sqrt(0.5)
	f.filter.SetWidth(math.MakeVector2(radius, radius))
	return f
}

func (f *Filtered) AddSample(sample *sampler.Sample, color math.Vector3) {
	x, y := int32(sample.Coordinates.X), int32(sample.Coordinates.Y)

	o := sample.RelativeOffset
	o.X -= 1
	o.Y -= 1
	w := f.filter.Evaluate(o)
	f.addPixel(x + 1, y + 1, color, w)

	o = sample.RelativeOffset
	o.Y -= 1
	w = f.filter.Evaluate(o)
	f.addPixel(x, y + 1, color, w)

	o = sample.RelativeOffset
	o.X += 1
	o.Y -= 1
	w = f.filter.Evaluate(o)
	f.addPixel(x - 1, y + 1, color, w)

	o = sample.RelativeOffset
	o.X -= 1
	w = f.filter.Evaluate(o)
	f.addPixel(x + 1, y, color, w)

	// center
	w = f.filter.Evaluate(sample.RelativeOffset)
	f.addPixel(x, y, color, w)

	o = sample.RelativeOffset
	o.X += 1
	w = f.filter.Evaluate(o)
	f.addPixel(x - 1, y, color, w)

	o = sample.RelativeOffset
	o.X -= 1
	o.Y += 1
	w = f.filter.Evaluate(o)
	f.addPixel(x + 1, y - 1, color, w)

	o = sample.RelativeOffset
	o.Y += 1
	w = f.filter.Evaluate(o)
	f.addPixel(x, y - 1, color, w)

	o = sample.RelativeOffset
	o.X += 1
	o.Y += 1
	w = f.filter.Evaluate(o)
	f.addPixel(x - 1, y - 1, color, w)

}


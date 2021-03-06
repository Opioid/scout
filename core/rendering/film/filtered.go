package film

import (
	"github.com/Opioid/scout/core/rendering/film/filter"
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Filtered struct {
	film
	filter filter.Filter
}

func NewFiltered(dimensions math.Vector2i, exposure float32, tonemapper tonemapping.Tonemapper) *Filtered {
	f := new(Filtered)
	f.resize(dimensions)
	f.exposure = exposure
	f.tonemapper = tonemapper

	radius := float32(0.8)//*/math32.Sqrt(0.5)

//	f.filter = filter.NewTriangle(math.MakeVector2(radius, radius))
	f.filter = filter.NewGaussian(math.MakeVector2(radius, radius), 0.3)
//	f.filter = filter.NewMitchellNetravali(math.MakeVector2(radius, radius), 1.0 / 3.0, 1.0 / 3.0)

	return f
}

func (f *Filtered) AddSample(sample *sampler.CameraSample, color math.Vector3, start, end math.Vector2i) {
	x, y := int32(sample.Coordinates.X), int32(sample.Coordinates.Y)

	leftEdge, rightEdge, topEdge, bottomEdge := false, false, false, false
	if x == start.X && x != 0 {
		leftEdge = true
	}
	if x == end.X - 1 && x < f.dimensions.X - 1 {
		rightEdge = true
	}
	if y == start.Y && y != 0 {
		topEdge = true
	}
	if y == end.Y - 1  && y < f.dimensions.Y - 1 {
		bottomEdge = true
	}

	o := sample.RelativeOffset
	o.X += 1.0
	o.Y += 1.0
	w := f.filter.Evaluate(o)

	if leftEdge || topEdge {
		f.atomicAddPixel(x - 1, y - 1, color, w)
	} else {
		f.addPixel(x - 1, y - 1, color, w)
	}

	o = sample.RelativeOffset
	o.Y += 1.0
	w = f.filter.Evaluate(o)

	if topEdge {
		f.atomicAddPixel(x, y - 1, color, w)
	} else {
		f.addPixel(x, y - 1, color, w)
	}

	o = sample.RelativeOffset
	o.X -= 1.0
	o.Y += 1.0
	w = f.filter.Evaluate(o)

	if rightEdge || topEdge {
		f.atomicAddPixel(x + 1, y - 1, color, w)
	} else { 
		f.addPixel(x + 1, y - 1, color, w)
	}	

	o = sample.RelativeOffset
	o.X += 1.0
	w = f.filter.Evaluate(o)

	if leftEdge {
		f.atomicAddPixel(x - 1, y, color, w)
	} else {
		f.addPixel(x - 1, y, color, w)
	}
	
	// center
	w = f.filter.Evaluate(sample.RelativeOffset)

	if leftEdge || rightEdge || topEdge || bottomEdge {
		f.atomicAddPixel(x, y, color, w)
	} else {
		f.addPixel(x, y, color, w)
	}
	
	o = sample.RelativeOffset
	o.X -= 1.0
	w = f.filter.Evaluate(o)

	if rightEdge {
		f.atomicAddPixel(x + 1, y, color, w)
	} else {
		f.addPixel(x + 1, y, color, w)
	}

	o = sample.RelativeOffset
	o.X += 1.0
	o.Y -= 1.0
	w = f.filter.Evaluate(o)

	if leftEdge || bottomEdge {
		f.atomicAddPixel(x - 1, y + 1, color, w)
	} else {
		f.addPixel(x - 1, y + 1, color, w)
	}
	
	o = sample.RelativeOffset
	o.Y -= 1.0
	w = f.filter.Evaluate(o)

	if bottomEdge {
		f.atomicAddPixel(x, y + 1, color, w)
	} else {
		f.addPixel(x, y + 1, color, w)
	}

	o = sample.RelativeOffset
	o.X -= 1.0
	o.Y -= 1.0
	w = f.filter.Evaluate(o)

	if rightEdge || bottomEdge {
		f.atomicAddPixel(x + 1, y + 1, color, w)
	} else {
		f.addPixel(x + 1, y + 1, color, w)
	}
}


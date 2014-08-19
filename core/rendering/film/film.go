package film

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"image"
)

type Film interface {
	Dimensions() math.Vector2i 

	AddSample(sample *sampler.Sample, color math.Vector3)

	RGBA() *image.RGBA
}

type pixel struct {
	color math.Vector3
	weightSum float32
}

type film struct {
	dimensions math.Vector2i

	pixels []pixel
}

func (f *film) Dimensions() math.Vector2i {
	return f.dimensions
}

func (f *film) resize(dimensions math.Vector2i) {
	f.dimensions = dimensions
	f.pixels = make([]pixel, dimensions.X * dimensions.Y)
}

func (f *film) at(x, y int) pixel {
	return f.pixels[f.dimensions.X * y + x]
}

func (f *film) addPixel(x, y int, color math.Vector3) {
	p := &f.pixels[f.dimensions.X * y + x]
	p.color.AddAssign(color)
	p.weightSum += 1.0
}
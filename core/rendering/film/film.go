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

type base struct {
	dimensions math.Vector2i

	pixels []pixel
}

func (b *base) resize(dimensions math.Vector2i) {
	b.dimensions = dimensions
	b.pixels = make([]pixel, dimensions.X * dimensions.Y)
}

func (b *base) at(x, y int) pixel {
	return b.pixels[b.dimensions.X * y + x]
}

func (b *base) addPixel(x, y int, color math.Vector3) {
	p := &b.pixels[b.dimensions.X * y + x]
	p.color.AddAssign(color)
	p.weightSum += 1.0
}
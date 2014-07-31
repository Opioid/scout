package film

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"image"
	"image/color"
)

type Unfiltered struct {
	base
}

func NewUnfiltered(dimensions math.Vector2i) *Unfiltered {
	f := new(Unfiltered)
	f.resize(dimensions)
	return f
}

func (f *Unfiltered) Dimensions() math.Vector2i {
	return f.dimensions
}

func (f *Unfiltered) AddSample(sample *sampler.Sample, color math.Vector3) {
	x, y := int(sample.Coordinates.X), int(sample.Coordinates.Y)

	f.addPixel(x, y, color)
}

func (f *Unfiltered) RGBA() *image.RGBA {
	image := image.NewRGBA(image.Rect(0, 0, f.dimensions.X, f.dimensions.Y))

	for y := 0; y < f.dimensions.Y; y++ {
		for x := 0; x < f.dimensions.X; x++ {
			pixel := f.at(x, y)
			iw := 1.0 / pixel.weightSum
			r := uint8(255.0 * pixel.color.X * iw)
			g := uint8(255.0 * pixel.color.Y * iw)
			b := uint8(255.0 * pixel.color.Z * iw)

/*
			if pixel.weightSum > 1.0 {
				r = 255
				g = 127
				b = 127
			}
*/
			image.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return image
}
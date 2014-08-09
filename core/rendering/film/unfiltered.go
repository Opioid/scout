package film

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/math"
	"image"
	gocolor "image/color"
	"runtime"
	"sync"
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
	target := image.NewRGBA(image.Rect(0, 0, f.dimensions.X, f.dimensions.Y))

	numTaks := runtime.GOMAXPROCS(0)

	a := f.dimensions.Y / numTaks

	start := math.Vector2i{0, 0}
	end   := math.Vector2i{f.dimensions.X, a}

	wg := sync.WaitGroup{}

	for i := 0; i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i, t *image.RGBA) {
			f.process(s, e, t)
			wg.Done()
		}(start, end, target)

		start.Y += a

		if i == numTaks - 1 {
			end.Y = f.dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	return target
}

func (f *Unfiltered) process(start, end math.Vector2i, target *image.RGBA) {
	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			pixel := f.at(x, y)
			iw := 1.0 / pixel.weightSum
			r := uint8(255.0 * color.LinearToSrgb(pixel.color.X * iw))
			g := uint8(255.0 * color.LinearToSrgb(pixel.color.Y * iw))
			b := uint8(255.0 * color.LinearToSrgb(pixel.color.Z * iw))
/*
			if pixel.weightSum > 1.0 {
				r = 255
				g = 127
				b = 127
			}
*/
			target.Set(x, y, gocolor.RGBA{r, g, b, 255})
		}
	}
}
package film

import (
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/math"
	"image"
	gocolor "image/color"
	"runtime"
	"sync"
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

func (f *Unfiltered) AddSample(sample *sampler.Sample, color math.Vector3) {
	x, y := int32(sample.Coordinates.X), int32(sample.Coordinates.Y)

	f.addPixel(x, y, color, 1)
}

func (f *Unfiltered) RGBA() *image.RGBA {
	target := image.NewRGBA(image.Rect(0, 0, int(f.dimensions.X), int(f.dimensions.Y)))

	numTaks := int32(runtime.GOMAXPROCS(0))

	a := f.dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(f.dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := int32(0); i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i) {
			f.processTonemapped(s, e, target)
			wg.Done()
		}(start, end)

		start.Y += a

		if i == numTaks - 2 {
			end.Y = f.dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	return target
}

func (f *Unfiltered) Float32x3() []float32 {
	target := make([]float32, f.dimensions.X * f.dimensions.Y * 3)

	for y := int32(0); y < f.dimensions.Y; y++ {
		for x := int32(0); x < f.dimensions.X; x++ {
			pixel := f.at(x, y)
			c := pixel.color.Div(pixel.weightSum)

			o := (f.dimensions.X * y + x) * 3

			target[o + 0] = c.X
			target[o + 1] = c.Y
			target[o + 2] = c.Z
		}
	}

	return target
}

func (f *Unfiltered) processTonemapped(start, end math.Vector2i, target *image.RGBA) {
	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			pixel := f.at(x, y)
			c := pixel.color.Div(pixel.weightSum)

			exposed := expose(c, f.exposure)

			tonemapped := f.tonemapper.Tonemap(exposed)
			r := uint8(255 * color.LinearToSrgb(tonemapped.X))
			g := uint8(255 * color.LinearToSrgb(tonemapped.Y))
			b := uint8(255 * color.LinearToSrgb(tonemapped.Z))

			target.Set(int(x), int(y), gocolor.RGBA{r, g, b, 255})
		}
	}
}
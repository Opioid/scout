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
	x, y := int(sample.Coordinates.X), int(sample.Coordinates.Y)

	f.addPixel(x, y, color)
}

func (f *Unfiltered) RGBA() *image.RGBA {
	target := image.NewRGBA(image.Rect(0, 0, f.dimensions.X, f.dimensions.Y))

	numTaks := runtime.GOMAXPROCS(0)

	a := f.dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(f.dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := 0; i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i) {
			f.process(s, e, target)
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

func (f *Unfiltered) process(start, end math.Vector2i, target *image.RGBA) {
	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			pixel := f.at(x, y)
			c := pixel.color.Div(pixel.weightSum)
			exposed := expose(c, f.exposure)
			tonemapped := f.tonemapper.Tonemap(exposed)
			r := uint8(255.0 * color.LinearToSrgb(tonemapped.X))
			g := uint8(255.0 * color.LinearToSrgb(tonemapped.Y))
			b := uint8(255.0 * color.LinearToSrgb(tonemapped.Z))

			target.Set(x, y, gocolor.RGBA{r, g, b, 255})
		}
	}
}
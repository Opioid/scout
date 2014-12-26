package film

import (
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32/atomic"
	"image"
	gocolor "image/color"
	"runtime"
	"sync"
)

type Film interface {
	Dimensions() math.Vector2i 

	AddSample(sample *sampler.CameraSample, color math.Vector3, start, end math.Vector2i)

	RGBA() *image.RGBA

	Float32x3() []float32
}

type pixel struct {
	color math.Vector3
	weightSum float32
}

type film struct {
	dimensions math.Vector2i
	exposure float32

	pixels []pixel

	tonemapper tonemapping.Tonemapper
}

func (f *film) Dimensions() math.Vector2i {
	return f.dimensions
}


func (f *film) RGBA() *image.RGBA {
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

func (f *film) Float32x3() []float32 {
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

func (f *film) processTonemapped(start, end math.Vector2i, target *image.RGBA) {
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

func (f *film) resize(dimensions math.Vector2i) {
	f.dimensions = dimensions
	f.pixels = make([]pixel, dimensions.X * dimensions.Y)
}

func (f *film) at(x, y int32) pixel {
	return f.pixels[f.dimensions.X * y + x]
}

func (f *film) addPixel(x, y int32, color math.Vector3, weight float32) {
	if x < 0 || x >= f.dimensions.X || y < 0 || y >= f.dimensions.Y {
		return
	}

	p := &f.pixels[f.dimensions.X * y + x]
	p.color.AddAssign(color.Scale(weight))
	p.weightSum += weight
}

func (f *film) atomicAddPixel(x, y int32, color math.Vector3, weight float32) {
	if x < 0 || x >= f.dimensions.X || y < 0 || y >= f.dimensions.Y {
		return
	}

	p := &f.pixels[f.dimensions.X * y + x]
	atomic.AddFloat32(&p.color.X, color.X * weight)
	atomic.AddFloat32(&p.color.Y, color.Y * weight)
	atomic.AddFloat32(&p.color.Z, color.Z * weight)
	atomic.AddFloat32(&p.weightSum, weight)
}

func expose(color math.Vector3, exposure float32) math.Vector3 {
	return color.Scale(math.Exp2(exposure))
}
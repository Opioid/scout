package texture

import (
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/math"
	goimage "image"
	gocolor "image/color"
	"runtime"
	"sync"
)

type Buffer struct {
	dimensions math.Vector2i
	data []math.Vector4
}

func (b *Buffer) Dimensions() math.Vector2i {
	return b.dimensions
}

func (b *Buffer) Data() []math.Vector4 {
	return b.data
}

func (b *Buffer) Resize(dimensions math.Vector2i) {
	b.dimensions = dimensions
	b.data = make([]math.Vector4, dimensions.X * dimensions.Y)
}

func (b *Buffer) At(x, y int32) math.Vector4 {
	return b.data[b.dimensions.X * y + x]
}

func (b *Buffer) Set(x, y int32, color math.Vector4) {
	b.data[b.dimensions.X * y + x] = color
}

func (b *Buffer) SetRgb(x, y int32, color math.Vector3) {
	v := &b.data[b.dimensions.X * y + x]

	v.X = color.X
	v.Y = color.Y
	v.Z = color.Z
}

func (b *Buffer) SetChannel(x, y, c int32, value float32) {
	b.data[b.dimensions.X * y + x].Set(c, value)
}

func (b *Buffer) RGBA() *goimage.RGBA {
	target := goimage.NewRGBA(goimage.Rect(0, 0, int(b.dimensions.X), int(b.dimensions.Y)))

	numTaks := int32(runtime.GOMAXPROCS(0))

	a := b.dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(b.dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := int32(0); i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i) {
			b.process(s, e, target)
			wg.Done()
		}(start, end)

		start.Y += a

		if i == numTaks - 2 {
			end.Y = b.dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	return target
}

func (buf *Buffer) process(start, end math.Vector2i, target *goimage.RGBA) {
	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			c := buf.At(x, y)
			r := uint8(255.0 * color.LinearToSrgb(c.X))
			g := uint8(255.0 * color.LinearToSrgb(c.Y))
			b := uint8(255.0 * color.LinearToSrgb(c.Z))
			a := uint8(255.0 * c.W)

			target.Set(int(x), int(y), gocolor.RGBA{r, g, b, a})
		}
	}
}
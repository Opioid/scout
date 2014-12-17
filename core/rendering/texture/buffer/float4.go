package buffer

import (
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/math"
	"image"
	gocolor "image/color"
	"runtime"
	"sync"
)

type Float4 struct {
	buffer
	data []math.Vector4
}

func NewFloat4(dimensions math.Vector2i) *Float4 {
	b := Float4{}
	b.Resize(dimensions)
	return &b
}

func (b *Float4) Resize(dimensions math.Vector2i) {
	b.dimensions = dimensions
	b.data = make([]math.Vector4, dimensions.X * dimensions.Y)
}

func (b *Float4) At(x, y int32) math.Vector4 {
	return b.data[b.dimensions.X * y + x]
}

func (b *Float4) Set(x, y int32, color math.Vector4) {
	b.data[b.dimensions.X * y + x] = color
}

func (b *Float4) SetRgb(x, y int32, color math.Vector3) {
	v := &b.data[b.dimensions.X * y + x]

	v.X = color.X
	v.Y = color.Y
	v.Z = color.Z
}

func (b *Float4) SetChannel(x, y, c int32, value float32) {
	b.data[b.dimensions.X * y + x].Set(c, value)
}

func (b *Float4) RGBA() *image.RGBA {
	target := image.NewRGBA(image.Rect(0, 0, int(b.dimensions.X), int(b.dimensions.Y)))

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

func (buf *Float4) process(start, end math.Vector2i, target *image.RGBA) {
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
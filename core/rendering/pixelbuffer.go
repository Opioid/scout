package rendering

import (
	"github.com/Opioid/scout/base/math"
	"image"
	"image/color"
)

type PixelBuffer struct {
	dimensions math.Vector2i

	pixels []math.Vector4
}

func NewPixelBuffer(dimensions math.Vector2i) *PixelBuffer {
	buffer := new(PixelBuffer)
	buffer.dimensions = dimensions
	buffer.pixels = make([]math.Vector4, dimensions.X * dimensions.Y)
	return buffer 
}

func (b *PixelBuffer) Dimensions() math.Vector2i {
	return b.dimensions
}

func (b *PixelBuffer) At(x, y int) math.Vector4 {
	return b.pixels[b.dimensions.X * y + x]
}

func (b *PixelBuffer) Set(x, y int, color math.Vector4) {
	b.pixels[b.dimensions.X * y + x] = color
}

func (b *PixelBuffer) RGBA() *image.RGBA {
	image := image.NewRGBA(image.Rect(0, 0, b.dimensions.X, b.dimensions.Y))

	for y := 0; y < b.dimensions.Y; y++ {
		for x := 0; x < b.dimensions.X; x++ {
			pixel := b.At(x, y)
			r := uint8(255.0 * pixel.X)
			g := uint8(255.0 * pixel.Y)
			b := uint8(255.0 * pixel.Z)
			a := uint8(255.0 * pixel.W)

			image.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return image
}
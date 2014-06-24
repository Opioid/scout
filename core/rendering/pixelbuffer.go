package rendering

import (
	"github.com/Opioid/scout/base/math"
	"image"
	"image/color"
)

type PixelBuffer struct {
	dimensions math.Vector2

	pixels []math.Vector4
}

func NewPixelBuffer(dimensions math.Vector2) *PixelBuffer {
	buffer := new(PixelBuffer)
	buffer.dimensions = dimensions
	buffer.pixels = make([]math.Vector4, dimensions.X * dimensions.Y)
	return buffer 
}

func (r *PixelBuffer) Dimensions() math.Vector2 {
	return r.dimensions
}

func (r *PixelBuffer) At(x, y int) math.Vector4 {
	return r.pixels[r.dimensions.X * y + x]
}

func (r *PixelBuffer) Set(x, y int, color math.Vector4) {
	r.pixels[r.dimensions.X * y + x] = color
}

func (r *PixelBuffer) RGBA() *image.RGBA {
	image := image.NewRGBA(image.Rect(0, 0, r.dimensions.X, r.dimensions.Y))

	for y := 0; y < r.dimensions.Y; y++ {
		for x := 0; x < r.dimensions.X; x++ {
			pixel := r.At(x, y)
			r := uint8(255.0 * pixel.X)
			g := uint8(255.0 * pixel.Y)
			b := uint8(255.0 * pixel.Z)
			a := uint8(255.0 * pixel.W)

			image.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return image
}
package texture

import (
	"github.com/Opioid/scout/base/math"
)

type imageBuffer struct {
	dimensions math.Vector2i
	data []math.Vector4
}

func (b *imageBuffer) resize(dimensions math.Vector2i) {
	b.dimensions = dimensions
	b.data = make([]math.Vector4, dimensions.X * dimensions.Y)
}

func (b *imageBuffer) at(x, y int) math.Vector4 {
	return b.data[b.dimensions.X * y + x]
}

func (b *imageBuffer) set(x, y int, color math.Vector4) {
	b.data[b.dimensions.X * y + x] = color
}

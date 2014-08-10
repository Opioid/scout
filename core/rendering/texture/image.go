package texture

import (
	"github.com/Opioid/scout/base/math"
)

type image struct {
	dimensions math.Vector2i
	data []math.Vector4
}

func (i *image) resize(dimensions math.Vector2i) {
	i.dimensions = dimensions
	i.data = make([]math.Vector4, dimensions.X * dimensions.Y)
}

func (i *image) at(x, y int) math.Vector4 {
	return i.data[i.dimensions.X * y + x]
}

func (i *image) set(x, y int, color math.Vector4) {
	i.data[i.dimensions.X * y + x] = color
}

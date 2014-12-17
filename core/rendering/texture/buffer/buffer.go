package buffer

import (
	"github.com/Opioid/scout/base/math"
	"image"
)

type Buffer interface {
	Dimensions() math.Vector2i
	Resize(dimensions math.Vector2i)
	At(x, y int32) math.Vector4
	Set(x, y int32, color math.Vector4)
	SetRgb(x, y int32, color math.Vector3)
	SetChannel(x, y, c int32, value float32)
	RGBA() *image.RGBA
}

type buffer struct {
	dimensions math.Vector2i
}

func (b *buffer) Dimensions() math.Vector2i {
	return b.dimensions
}

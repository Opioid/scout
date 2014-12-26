package buffer

import (
	"github.com/Opioid/scout/base/math"
	_ "image"
)

const (
	Float3 = iota
	Float4 = iota
)

type Buffer interface {
	Type() uint32
	Dimensions() math.Vector2i
	Resize(dimensions math.Vector2i)
	At(x, y int32) math.Vector4
	Set(x, y int32, color math.Vector4)
	SetRgb(x, y int32, color math.Vector3)
	SetChannel(x, y, c int32, value float32)
//	RGBA() *image.RGBA
}

func New(t uint32, dimensions math.Vector2i) Buffer {
	var buffer Buffer

	switch t {
	case Float3:
		buffer = newFloat3(dimensions)
	case Float4:
		buffer = newFloat4(dimensions)
	default: 
		buffer = newFloat4(dimensions)
	}

	return buffer
}

type buffer struct {
	dimensions math.Vector2i
}

func (b *buffer) Dimensions() math.Vector2i {
	return b.dimensions
}
